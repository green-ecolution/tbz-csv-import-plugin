package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

type CSVConverter struct {
	csvFile       multipart.File
	csvFileHeader *multipart.FileHeader
}

func NewCSVConverter(fileHeader *multipart.FileHeader, file multipart.File) *CSVConverter {
	return &CSVConverter{
		csvFile:       file,
		csvFileHeader: fileHeader,
	}
}

func (c *CSVConverter) Convert(ctx context.Context) ([]CsvTree, error) {
	if _, err := c.csvFile.Seek(0, 0); err != nil {
		return nil, err
	}

	trees, err := c.mapCSVToTrees(ctx)
	if err != nil {
		return nil, err
	}

	return trees, nil
}

func (c *CSVConverter) IsValid() bool {
	if strings.ToLower(filepath.Ext(c.csvFileHeader.Filename)) != ".csv" || c.csvFileHeader.Header.Get("Content-Type") != "text/csv" {
		return false
	}

	csvReader := csv.NewReader(c.csvFile)
	headers, err := csvReader.Read()
	if err != nil {
		return false
	}

	if !hasExpectedHeaders(headers) {
		return false
	}

	if _, err := csvReader.ReadAll(); err != nil {
		return false
	}

	return true
}

func hasExpectedHeaders(headers []string) bool {
	expectedHeaders := cfg.CsvHeaders
	if len(headers) != len(expectedHeaders) {
		return false
	}

	for i, header := range headers {
		if header != expectedHeaders[i] {
			return false
		}
	}

	return true
}

func (c *CSVConverter) mapCSVToTrees(_ context.Context) ([]CsvTree, error) {
	r := csv.NewReader(c.csvFile)
	r.LazyQuotes = true
	header, err := r.Read()
	if err != nil {
		slog.Error("failed to read csv", "error", err)
		return nil, errors.Join(err, errors.New("failed to read CSV"))
	}

	headerIndexMap := c.createHeaderIndexMap(header)
	var trees []CsvTree
	for i := range NumberSequence(1) {
		row, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		tree, err := c.parseRowToTree(i, row, headerIndexMap)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func (c *CSVConverter) createHeaderIndexMap(header []string) map[string]int {
	headerIndexMap := make(map[string]int, len(header))
	for i, h := range header {
		headerIndexMap[h] = i
	}
	return headerIndexMap
}

func (c *CSVConverter) parseRowToTree(rowIdx int, row []string, headerIndexMap map[string]int) (CsvTree, error) {
	// Helper function for validating and retrieving a field from the row
	getField := func(header string) (string, error) {
		idx, exists := headerIndexMap[header]
		if !exists || idx >= len(row) {
			return "", errors.New(fmt.Sprintf("header '%s' not found or index out of bounds at row: %d", header, rowIdx))
		}
		value := row[idx]
		if value == "" {
			return "", errors.New(fmt.Sprintf("invalid '%s' value at row: %d", header, rowIdx))
		}
		return value, nil
	}

	parseFloat := func(value string, fieldName string) (float64, error) {
		parsedValue, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
		if err != nil {
			return 0, errors.Join(err, fmt.Errorf("invalid '%s' value at row: %d", fieldName, rowIdx))
		}
		return parsedValue, nil
	}

	parseInt := func(value string, fieldName string) (int, error) {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.Join(err, fmt.Errorf("invalid '%s' value at row: %d", fieldName, rowIdx))
		}
		return parsedValue, nil
	}

	expectedHeaders := cfg.CsvHeaders
	area, err := getField(expectedHeaders[0])
	if err != nil {
		return CsvTree{}, err
	}

	street, err := getField(expectedHeaders[1])
	if err != nil {
		return CsvTree{}, err
	}

	treeNumber, err := getField(expectedHeaders[2])
	if err != nil {
		return CsvTree{}, err
	}

	species, err := getField(expectedHeaders[3])
	if err != nil {
		species = "" // Default to empty string
	}

	latitudeStr, err := getField(expectedHeaders[4])
	if err != nil {
		return CsvTree{}, err
	}
	latitude, err := parseFloat(latitudeStr, "Hochwert")
	if err != nil {
		return CsvTree{}, err
	}

	longitudeStr, err := getField(expectedHeaders[5])
	if err != nil {
		return CsvTree{}, err
	}
	longitude, err := parseFloat(longitudeStr, "Rechtswert")
	if err != nil {
		return CsvTree{}, err
	}

	plantingYearStr, err := getField(expectedHeaders[6])
	if err != nil {
		return CsvTree{}, err
	}
	plantingYear, err := parseInt(plantingYearStr, "Pflanzjahr")
	if err != nil {
		return CsvTree{}, err
	}

	tree := CsvTree{
		Area:         area,
		Street:       street,
		TreeNumber:   treeNumber,
		Species:      species,
		Hochwert:     latitude,
		Rechtswert:   longitude,
		PlantingYear: plantingYear,
	}

	return tree, nil
}
