package app

import (
	"fmt"
	"os"
	"path/filepath"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils"
	"time"
)

type ProcessInvetoryStep interface {
	Process() error
}

type inventoryProcessor struct {
	config config.Config
	logger utils.Logger
}

func NewProcessInvetoryStep(config config.Config, logger utils.Logger) ProcessInvetoryStep {
	return &inventoryProcessor{
		config: config,
		logger: logger,
	}
}

func (p *inventoryProcessor) Process() error {
	csvFiles, err := p.config.GetCSVFilesWithRecordedEquipment()
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("Failed to get CSV files: %v", err))
	}

	csvFile := NewCSVFile()

	var recordedInventoryData []CSVContent
	for _, file := range csvFiles {
		content, err := csvFile.Read(file)
		if err != nil {
			p.logger.Fatal(fmt.Sprintf("Failed to read CSV file '%s': %v", file, err))
		}
		recordedInventoryData = append(recordedInventoryData, content)
	}

	recordedInventory := NewRecordedInventory(recordedInventoryData)

	content, err := csvFile.Read(p.config.GetAbsoluteInventoryCSVFileName())
	if err != nil {
		return fmt.Errorf("failed to read CSV file '%s': %v", p.config.GetAbsoluteInventoryCSVFileName(), err)
	}

	inventoryData, err := NewInventoryData(content, p.config, p.logger)
	if err != nil {
		return fmt.Errorf("failed to init inventory data: %v", err)
	}

	inventoryMap, err := recordedInventory.AsMap()
	if err != nil {
		return fmt.Errorf("failed to convert recorded inventory to map: %v", err)
	}

	p.logger.Info("recorded equipment:")
	p.logger.Info("")
	p.logger.InfoIndented("equipment     : amount")
	p.logger.InfoIndented("----------------------")
	for key, value := range inventoryMap {
		p.logger.InfoIndented(fmt.Sprintf("%-13s : %5d", key, value))
	}
	p.logger.Info("")

	err = inventoryData.UpdateInventory(inventoryMap)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	resultDir := filepath.Join(p.config.WorkingDir, "result")

	err = os.MkdirAll(resultDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create result directory: %v", err)
	}

	err = csvFile.Write(
		filepath.Join(resultDir, fmt.Sprintf("result_%s.csv", time.Now().Format("2006-01-02_15-04-05"))),
		inventoryData.GetContent(),
	)
	if err != nil {
		return fmt.Errorf("failed to write result csv: %v", err)
	}

	return nil
}
