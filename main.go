package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"thwInventoryMerge/app"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils"
	"time"
)

func main() {
	
	logger := utils.NewLogger()

	var configPath string
	
	flag.StringVar(&configPath, "c", "config.json", "the config file path")
	flag.Parse()

	config, err := config.LoadConfig(configPath, logger)
	if err != nil {
		log.Fatalf("Failed to load config from path %s: %v", configPath, err)
	}

	if config.WorkingDir == "" {
		config.WorkingDir = getExecutablePath()
	}

	if configPath == "" {
		filepath.Join(config.WorkingDir, "config.json")
	}

	csvFiles, err := config.GetCSVFilesWithRecordedEquipment()
	if err != nil {
		log.Fatalf("Failed to get CSV files: %v", err)
	}

	csvFile := app.NewCSVFile()

	var recordedInventoryData []app.CSVContent
	for _, file := range csvFiles {
		content, err := csvFile.Read(file)
		if err != nil {
			log.Fatalf("Failed to read CSV file '%s': %v", file, err)
		}
		recordedInventoryData = append(recordedInventoryData, content)
	}

	recordedInventory := app.NewRecordedInventory(recordedInventoryData)

	content, err := csvFile.Read(config.GetAbsoluteInventoryCSVFileName())
	if err != nil {
		log.Fatalf("Failed to read CSV file '%s': %v", config.GetAbsoluteInventoryCSVFileName(), err)
	}

	inventoryData, err := app.NewInventoryData(content, logger)
	if err != nil {
		log.Fatalf("Failed to init inventory data: %v", err)
	}

	inventoryMap, err := recordedInventory.AsMap()
	if err != nil {
		log.Fatalf("Failed to convert recorded inventory to map: %v", err)
	}

	logger.Info("recorded equipment:")
	logger.Info("")
	logger.InfoIndented("equipment     : amount")
	logger.InfoIndented("----------------------")
	for key, value := range inventoryMap {
		logger.InfoIndented(fmt.Sprintf("%-13s : %5d", key, value))
	}
	logger.Info("")

	inventoryData.UpdateInventory(inventoryMap)

	resultDir := filepath.Join(config.WorkingDir, "result")

	err = os.MkdirAll(resultDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create result directory: %v", err)
	}

	err = csvFile.Write(
		filepath.Join(resultDir, fmt.Sprintf("result_%s.csv", time.Now().Format("2006-01-02_15-04-05"))),
		inventoryData.GetContent(),
	)
	if err != nil {
		log.Fatalf("Failed to write result csv: %v", err)
	}

	// Keep the terminal open
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func getExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	return filepath.Dir(exePath)   
}