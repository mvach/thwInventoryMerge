package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"thwInventoryMerge/utils"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type CSVContent [][]string

type CSVFile interface {
	Read(filePath string, encoding encoding.Encoding) (CSVContent, error)

	Write(filePath string, content CSVContent) error
}

type csvFile struct {
	logger utils.Logger
}

func NewCSVFile(logger utils.Logger) CSVFile {
	return &csvFile{
		logger: logger,
	}
}

func (c *csvFile) Read(filePath string, encoding encoding.Encoding) (CSVContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(transform.NewReader(file, encoding.NewDecoder()))
	reader.Comma = ';'
	reader.LazyQuotes = true
	content, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file '%s': %w", filePath, err)
	}

	return content, nil
}

func (c *csvFile) Write(filePath string, content CSVContent) error {
	// Open the file with O_WRONLY, O_CREATE, and O_TRUNC flags to clear contents if it exists
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Write the UTF-8 BOM
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return fmt.Errorf("failed to write UTF-8 BOM to CSV file: %w", err)
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	err = writer.WriteAll(content)
	if err != nil {
		return fmt.Errorf("failed to write into CSV file: %w", err)
	}

	writer.Flush()
	return nil
}
