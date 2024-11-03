package app

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVContent [][]string

type CSVFile interface {
	Read(filePath string) (CSVContent, error)

	Write(filePath string, content CSVContent) error
}

type csvFile struct {
}

func NewCSVFile() CSVFile {
	return &csvFile{}
}

func (c *csvFile) Read(filePath string) (CSVContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	content, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file '%s': %w", filePath, err)
	}
  
	return content, nil
}

func (c *csvFile) Write(filePath string, content CSVContent) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Write the UTF-8 BOM for Excel on Windows compatibility
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return fmt.Errorf("failed to write UTF-8 BOM to csv file: %v", err)
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
