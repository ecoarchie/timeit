package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/gocarina/gocsv"
)

var validHeaders = []string{"event", "wave", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}

type AthleteCSVParser struct {
	fileAddr  string
	separator string
}

func NewAthleteCSVParser(file, separator string) *AthleteCSVParser {
	return &AthleteCSVParser{
		fileAddr:  file,
		separator: separator,
	}
}

func (pp AthleteCSVParser) ReadCSV(f, separator string) ([]*entity.AthleteCSV, error) {
	file, err := os.Open(pp.fileAddr)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma, _ = utf8.DecodeRuneInString(pp.separator)
		return r // Allows use dot as delimiter and use quotes in CSV
	})
	athletes := []*entity.AthleteCSV{}
	if err := gocsv.UnmarshalFile(file, &athletes); err != nil { // Load clients from file
		return nil, err
	}
	return athletes, nil
}

func (pp AthleteCSVParser) ValidateHeaders(headersRow []string) []string {
	// var validHeaders []string = []string{"event", "wave", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}
	validHeadersMap := map[string]struct{}{}
	for _, h := range validHeaders {
		validHeadersMap[h] = struct{}{}
	}

	var invalidHeaders []string
	for _, h := range headersRow {
		if _, ok := validHeadersMap[h]; !ok {
			invalidHeaders = append(invalidHeaders, h)
		}
	}
	return invalidHeaders
}

func (pp AthleteCSVParser) GetHeaders(f, separator string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma, _ = utf8.DecodeRuneInString(separator)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return headers, nil
}

func (pp AthleteCSVParser) CopyToNewCSV(old, new, separator string, newHeaders []string) error {
	file, err := os.Open(old)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma, _ = utf8.DecodeRuneInString(separator)
	reader.LazyQuotes = true

	newFile, err := os.Create(new)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	writer.Comma, _ = utf8.DecodeRuneInString(separator)
	err = writer.Write(newHeaders)
	if err != nil {
		return err
	}
	_, err = reader.Read() // get rid of header record
	if err != nil {
		return err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return err
		}
		writer.Write(record)
	}

	err = writer.Error()
	if err != nil {
		return err
	}
	return nil
}

// the input here depends on ui, may change
func (pp AthleteCSVParser) HeadersToReplace(valid, invalid []string) (map[entity.InvalidHeader]entity.ValidHeader, error) {
	if len(valid) != len(invalid) {
		return nil, fmt.Errorf("number of invalid fields not equal to number of valid fields")
	}
	res := make(map[entity.InvalidHeader]entity.ValidHeader, len(invalid))
	for i, h := range invalid {
		res[h] = valid[i]
	}
	return res, nil
}

func (pp AthleteCSVParser) ReplaceInvalidHeaders(f, newF, separator string, oldHeaders []string, headersToReplace map[entity.InvalidHeader]entity.ValidHeader) []string {
	newHeaders := []string{}
	for _, h := range oldHeaders {
		if newH, ok := headersToReplace[strings.TrimSpace(h)]; ok {
			newHeaders = append(newHeaders, newH)
		} else {
			newHeaders = append(newHeaders, strings.TrimSpace(h))
		}
	}
	return newHeaders
}
