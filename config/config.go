package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"thwInventoryMerge/utils"
)

type Config struct {
	WorkingDir    string `json:"working_dir"`
	InventoryCSVFileName string `json:"inventory_csv_file_name"`
	InventoryCSVConfig   struct {
		EquipmentIDColumnName         string `json:"equipment_id_column_name"`
		EquipmentAvailableColumnName  string `json:"equipment_available_column_name"`
	} `json:"inventory_csv_config"`
	
	logger utils.Logger
}

func (c *Config) GetCSVFilesWithRecordedEquipment() ([]string, error) {
	var csvFiles []string

	files, err := os.ReadDir(c.WorkingDir)
	if err != nil {
		return nil, err
	}

	firstEquipment := true

	for _, file := range files {
		if !file.IsDir() &&
			filepath.Ext(file.Name()) == ".csv" &&
			filepath.Base(file.Name()) != c.InventoryCSVFileName {

				if firstEquipment {
					c.logger.Info("files with recorded equipment:")
					c.logger.Info("")
					firstEquipment = false
				}
			
			c.logger.InfoIndented(fmt.Sprintf("using '%s'", file.Name()))
			csvFiles = append(csvFiles, filepath.Join(c.WorkingDir, file.Name()))
		}
	}

	c.logger.Info("")

	return csvFiles, nil
}

func (c *Config) GetAbsoluteInventoryCSVFileName() string {
	return filepath.Join(c.WorkingDir, c.InventoryCSVFileName)
}

func LoadConfig(filePath string, logger utils.Logger) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config = Config{
		logger: logger,
	}

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
	if c.InventoryCSVFileName == "" {
		return errors.New("property inventory_csv_file_name is required")
	}
	if c.InventoryCSVConfig.EquipmentIDColumnName == "" {
		return errors.New("property inventory_csv_config.equipment_id_column_name is required")
	}
	if c.InventoryCSVConfig.EquipmentAvailableColumnName == "" {
		return errors.New("property inventory_csv_config.equipment_available_column_name is required")
	}
	return nil
}
