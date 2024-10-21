package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	WorkingDir    string `json:"working_dir"`
	ExcelFileName string `json:"excel_file_name"`
	ExcelConfig   struct {
		WorksheetName                 string `json:"worksheet_name"`
		EquipmentIDColumnName         string `json:"equipment_id_column_name"`
		EquipmentIDColumnIndex        *int
		EquipmentAvailableColumnName  string `json:"equipment_available_column_name"`
		EquipmentAvailableColumnIndex *int
	} `json:"excel_config"`
}

func (c *Config) GetCSVFiles() ([]string, error) {
	var csvFiles []string

	files, err := os.ReadDir(c.WorkingDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			csvFiles = append(csvFiles, filepath.Join(c.WorkingDir, file.Name()))
		}
	}

	return csvFiles, nil
}

func (c *Config) GetAbsoluteExcelFileName() string {
	return filepath.Join(c.WorkingDir, c.ExcelFileName)
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load invalid config file, %w", err)
	}

	err = config.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate the config file, %w", err)
	}

	return &config, nil
}

func (c Config) validate() error {
	if c.ExcelFileName == "" {
		return errors.New("property excel_file_name is required")
	}
	if c.ExcelConfig.WorksheetName == "" {
		return errors.New("property excel_config.worksheet_name is required")
	}
	if c.ExcelConfig.EquipmentIDColumnName == "" {
		return errors.New("property excel_config.equipment_id_column_name is required")
	}
	if c.ExcelConfig.EquipmentAvailableColumnName == "" {
		return errors.New("property excel_config.equipment_available_column_name is required")
	}
	return nil
}
