package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AthleteImporter interface{}

type AthleteImporterCSV struct {
	FileName     string
	Separator    string
	tmpFolder    string
	validHeaders []string
}

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

func NewAthleteImporterCSV(file, separator string) *AthleteImporterCSV {
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
	return &AthleteImporterCSV{
		FileName:     file,
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

func (ai AthleteImporterCSV) FilePath() string {
	return ai.tmpFolder + "/" + ai.FileName + ".csv"
}

func (ai AthleteImporterCSV) CompareHeaders() (userHeaders []string, matchingValidHeaders []string, err error) {
	userHeaders, err = ai.getHeaders()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading headers from file")
	}
	validHeadersMap := map[string]struct{}{}
	for _, h := range ai.validHeaders {
		validHeadersMap[h] = struct{}{}
	}
	for _, uh := range userHeaders {
		if _, ok := validHeadersMap[uh]; ok {
			matchingValidHeaders = append(matchingValidHeaders, uh)
		} else {
			matchingValidHeaders = append(matchingValidHeaders, "")
		}
	}
	return userHeaders, matchingValidHeaders, nil
}

func (ai AthleteImporterCSV) ReadCSV(validatedHeaders []string) ([]*AthleteCSV, error) {
	start := time.Now()
	_, err := ai.copyToNewCSVWithValidHeaders(validatedHeaders)
	if err != nil {
		return nil, errors.Wrap(err, "error copying data to new csv")
	}
	file, err := os.Open(ai.FilePath())
	if err != nil {
		return nil, err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma, _ = utf8.DecodeRuneInString(ai.Separator)
		return r // Allows use dot as delimiter and use quotes in CSV
	})
	athletes := []*AthleteCSV{}
	if err := gocsv.UnmarshalFile(file, &athletes); err != nil { // Load clients from file
		return nil, err
	}
	file.Close()
	err = os.Remove(ai.FilePath())
	if err != nil {
		return nil, fmt.Errorf("error deleting file")
	}
	fmt.Printf("ReadCSV took: %v\n", time.Since(start))
	return athletes, nil
}

func (ai AthleteImporterCSV) getHeaders() ([]string, error) {
	file, err := os.Open(ai.FilePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma, _ = utf8.DecodeRuneInString(ai.Separator)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return headers, nil
}

func (ai AthleteImporterCSV) copyToNewCSVWithValidHeaders(newHeaders []string) (string, error) {
	bb, err := os.ReadFile(ai.FilePath())
	if err != nil {
		return "", fmt.Errorf("error reading whole file")
	}
	reader := csv.NewReader(bytes.NewReader(bb))
	reader.Comma, _ = utf8.DecodeRuneInString(ai.Separator)
	reader.LazyQuotes = true

	data, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error reading whole csv file")
	}
	data[0] = newHeaders

	file, err := os.OpenFile(ai.FilePath(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Comma, _ = utf8.DecodeRuneInString(ai.Separator)
	err = writer.WriteAll(data)
	if err != nil {
		return "", fmt.Errorf("error writing to file")
	}
	return "", nil
}
