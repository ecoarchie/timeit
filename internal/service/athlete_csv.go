package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"unicode/utf8"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type AthleteCSV struct {
	Event       string `csv:"event"`
	Wave        string `csv:"wave"`
	Bib         int    `csv:"bib"`
	Chip        int    `csv:"tag"`
	FirstName   string `csv:"name"`
	LastName    string `csv:"surname"`
	Gender      string `csv:"gender"`
	DateOfBirth string `csv:"date of birth"`
	Phone       string `csv:"phone"`
	Comments    string `csv:"comments"`
}

type (
	InvalidHeader = string
	ValidHeader   = string
)

type AthleteCSVParser struct {
	FileAddr     string
	Separator    string
	tmpFolder    string
	validHeaders []string
}

func NewAthleteCSVParser(file, separator string) *AthleteCSVParser {
	validHeaders := []string{
		"event",
		"wave",
		"bib",
		"tag",
		"name",
		"surname",
		"gender",
		"date of birth",
		"phone",
		"comments",
	}
	return &AthleteCSVParser{
		FileAddr:     file,
		Separator:    separator,
		tmpFolder:    "tmp",
		validHeaders: validHeaders,
	}
}

func StoreTmpFile(r *http.Request) (string, error) {
	file, fileHeader, err := r.FormFile("athletes")
	if err != nil {
		return "", fmt.Errorf("error reading file from request")
	}
	defer file.Close()

	contentType := fileHeader.Header["Content-Type"][0]
	if contentType != "text/csv" {
		fmt.Println("file format", contentType)
		return "", fmt.Errorf("error file format is not csv")
	}
	token := uuid.New().String()
	tempFilePath := fmt.Sprintf("tmp/%s.csv", token)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("error could not save the file")
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", fmt.Errorf("error copying the file to tmp")
	}
	return token, nil
}

func (pp AthleteCSVParser) CompareHeaders() (userHeaders []string, matchingValidHeaders []string, err error) {
	userHeaders, err = pp.getHeaders()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading headers from file")
	}
	validHeadersMap := map[string]struct{}{}
	for _, h := range pp.validHeaders {
		validHeadersMap[h] = struct{}{}
	}
	fmt.Println("valid headers map", validHeadersMap)
	fmt.Println("user headers", userHeaders)
	fmt.Println("user headers len", len(userHeaders))
	for _, uh := range userHeaders {
		if _, ok := validHeadersMap[uh]; ok {
			matchingValidHeaders = append(matchingValidHeaders, uh)
		} else {
			matchingValidHeaders = append(matchingValidHeaders, "")
		}
	}
	return userHeaders, matchingValidHeaders, nil
}

func (pp AthleteCSVParser) ReadCSV(validatedHeaders []string) ([]*AthleteCSV, error) {
	newFileName, err := pp.copyToNewCSVWithValidHeaders(validatedHeaders)
	if err != nil {
		return nil, fmt.Errorf("error copying data to new csv")
	}
	file, err := os.Open(newFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma, _ = utf8.DecodeRuneInString(pp.Separator)
		return r // Allows use dot as delimiter and use quotes in CSV
	})
	athletes := []*AthleteCSV{}
	if err := gocsv.UnmarshalFile(file, &athletes); err != nil { // Load clients from file
		return nil, err
	}
	return athletes, nil
}

// func (pp AthleteCSVParser) validateHeaders(headersRow []string) []string {
// 	validHeadersMap := map[string]struct{}{}
// 	for _, h := range pp.validHeaders {
// 		validHeadersMap[h] = struct{}{}
// 	}

// 	var invalidHeaders []string
// 	for _, h := range headersRow {
// 		if _, ok := validHeadersMap[h]; !ok {
// 			invalidHeaders = append(invalidHeaders, h)
// 		}
// 	}
// 	return invalidHeaders
// }

func (pp AthleteCSVParser) getHeaders() ([]string, error) {
	file, err := os.Open(pp.tmpFolder + "/" + pp.FileAddr)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma, _ = utf8.DecodeRuneInString(pp.Separator)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return headers, nil
}

func (pp AthleteCSVParser) copyToNewCSVWithValidHeaders(newHeaders []string) (string, error) {
	file, err := os.Open(pp.tmpFolder + "/" + pp.FileAddr)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma, _ = utf8.DecodeRuneInString(pp.Separator)
	reader.LazyQuotes = true

	newFileName := pp.FileAddr + "New"
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", fmt.Errorf("error creating new file for copy")
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	writer.Comma, _ = utf8.DecodeRuneInString(pp.Separator)
	err = writer.Write(newHeaders)
	if err != nil {
		return "", err
	}
	_, err = reader.Read() // get rid of header record
	if err != nil {
		return "", err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return "", err
		}
		writer.Write(record)
	}

	err = writer.Error()
	if err != nil {
		return "", err
	}
	return newFileName, nil
}

// the input here depends on ui, may change
// func (pp AthleteCSVParser) HeadersToReplace(valid, invalid []string) (map[InvalidHeader]ValidHeader, error) {
// 	if len(valid) != len(invalid) {
// 		return nil, fmt.Errorf("number of invalid fields not equal to number of valid fields")
// 	}
// 	res := make(map[InvalidHeader]ValidHeader, len(invalid))
// 	for i, h := range invalid {
// 		res[h] = valid[i]
// 	}
// 	return res, nil
// }

// func (pp AthleteCSVParser) ReplaceInvalidHeaders(oldHeaders []string, headersToReplace map[InvalidHeader]ValidHeader) []string {
// 	newHeaders := []string{}
// 	for _, h := range oldHeaders {
// 		if newH, ok := headersToReplace[strings.TrimSpace(h)]; ok {
// 			newHeaders = append(newHeaders, newH)
// 		} else {
// 			newHeaders = append(newHeaders, strings.TrimSpace(h))
// 		}
// 	}
// 	return newHeaders
// }
