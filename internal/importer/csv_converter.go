package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/green-ecolution/tbz-csv-import-plugin/internal/entities"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/utils"
	"github.com/pkg/errors"
)

type CSVConverter struct {
	expectedHeaders []string
	fromEPSG        int
	toEPSG          int
	csvFile         *os.File
}

func NewCSVConverter(file *os.File) *CSVConverter {
	expectedCSVHeaders := strings.Split(strings.Trim(os.Getenv("CSV_HEADERS"), " "), ",")
	if len(expectedCSVHeaders) == 0 {
		log.Fatalf("Error getting CSV headers from environment variable. Please check the CSV_HEADERS variable.\n")
	}

	fromEPSGStr := os.Getenv("CSV_USED_EPSG")
	fromEPSG, err := strconv.Atoi(fromEPSGStr)
	if err != nil {
		log.Fatalf("Error converting EPSG from string to int: %v\n", err)
	}

	toEPSGStr := os.Getenv("CSV_TO_EPSG")
	toEPSG := 4326 // default to WGS84
	if toEPSGStr != "" {
		toEPSG, err = strconv.Atoi(toEPSGStr)
		if err != nil {
			log.Fatalf("Error converting EPSG from string to int: %v\n", err)
		}
	}

	return &CSVConverter{
		expectedHeaders: expectedCSVHeaders,
		fromEPSG:        fromEPSG,
		toEPSG:          toEPSG,
		csvFile:         file,
	}
}

func (c *CSVConverter) Convert(ctx context.Context) ([]*entities.TreeImport, error) {
	start := time.Now()
	if err := c.validateCsv(); err != nil {
		return nil, err
	}

	if _, err := c.csvFile.Seek(0, 0); err != nil {
		return nil, err
	}

	trees, err := c.mapCSVToTrees(ctx)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	slog.Info("Imported trees from CSV", "elapsed", elapsed)

	return trees, nil
}

func (c *CSVConverter) validateCsv() error {
	if !c.isCsvFile() {
		return errors.New("file is not a CSV file")
	}

	csvReader := csv.NewReader(c.csvFile)
	headers, err := csvReader.Read()
	if err != nil {
		return err
	}

	if !c.hasExpectedHeaders(headers) {
		return errors.New("csv file does not contain the expected headers")
	}

	_, err = csvReader.ReadAll()
	return err
}

func (c *CSVConverter) hasExpectedHeaders(headers []string) bool {
	if len(headers) != len(c.expectedHeaders) {
		return false
	}

	for i, header := range headers {
		if header != c.expectedHeaders[i] {
			return false
		}
	}

	return true
}

func (c *CSVConverter) isCsvFile() bool {
	fileExt := strings.ToLower(filepath.Ext(c.csvFile.Name()))
	slog.Debug("File extension", "ext", fileExt)
	return fileExt == ".csv"
}

func (c *CSVConverter) mapCSVToTrees(_ context.Context) ([]*entities.TreeImport, error) {
	r := csv.NewReader(c.csvFile)
	r.LazyQuotes = true
	header, err := r.Read()
	if err != nil {
		slog.Error("Failed to read CSV", "error", err)
		return nil, errors.Wrap(err, "failed to read CSV")
	}

	headerIndexMap := c.createHeaderIndexMap(header)
	transformer, err := NewGeoTransformer(c.fromEPSG, c.toEPSG)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transformer")
	}

	var trees []*entities.TreeImport
	for i := range utils.NumberSequence(1) {
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

	geoPoints := utils.Map(trees, func(tree *entities.TreeImport) GeoPoint {
		return GeoPoint{X: tree.Latitude, Y: tree.Longitude}
	})

	transformedPoints, err := transformer.TransformBatch(geoPoints)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to transform batch of points from EPSG %d to EPSG %d. err: %s", c.fromEPSG, c.toEPSG, err))
	}

	for i, tree := range trees {
		tree.Latitude = transformedPoints[i].X
		tree.Longitude = transformedPoints[i].Y
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

func (c *CSVConverter) parseRowToTree(rowIdx int, row []string, headerIndexMap map[string]int) (*entities.TreeImport, error) {
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
			return 0, errors.Wrap(err, fmt.Sprintf("invalid '%s' value at row: %d", fieldName, rowIdx))
		}
		return parsedValue, nil
	}

	parseInt := func(value string, fieldName string) (int, error) {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("invalid '%s' value at row: %d", fieldName, rowIdx))
		}
		return parsedValue, nil
	}

	area, err := getField(c.expectedHeaders[0])
	if err != nil {
		return nil, err
	}

	street, err := getField(c.expectedHeaders[1])
	if err != nil {
		return nil, err
	}

	treeNumber, err := getField(c.expectedHeaders[2])
	if err != nil {
		return nil, err
	}

	species, err := getField(c.expectedHeaders[3])
	if err != nil {
		species = "" // Default to empty string
	}

	latitudeStr, err := getField(c.expectedHeaders[4])
	if err != nil {
		return nil, err
	}
	latitude, err := parseFloat(latitudeStr, "Hochwert")
	if err != nil {
		return nil, err
	}

	longitudeStr, err := getField(c.expectedHeaders[5])
	if err != nil {
		return nil, err
	}
	longitude, err := parseFloat(longitudeStr, "Rechtswert")
	if err != nil {
		return nil, err
	}

	plantingYearStr, err := getField(c.expectedHeaders[6])
	if err != nil {
		return nil, err
	}
	plantingYear, err := parseInt(plantingYearStr, "Pflanzjahr")
	if err != nil {
		return nil, err
	}

	tree := &entities.TreeImport{
		Area:         area,
		Street:       street,
		Number:       treeNumber,
		Species:      species,
		Latitude:     latitude,
		Longitude:    longitude,
		PlantingYear: int32(plantingYear),
	}

	return tree, nil
}
