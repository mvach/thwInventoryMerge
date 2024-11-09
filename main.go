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
)

func main() {
	
	logger := utils.NewLogger()

	var configPath string
	var step string
	
	flag.StringVar(&configPath, "c", "config.json", "the config file path")
	flag.StringVar(&step, "s", "process", "the inventory step")
	flag.Parse()

	executablePath := getExecutablePath(logger)
	
	if configPath == "" {
		configPath = filepath.Join(executablePath, "config.json")
	}

	config, err := config.LoadConfig(configPath, logger)
	if err != nil {
		log.Fatalf("Failed to load config from path %s: %v", configPath, err)
	}

	if config.WorkingDir == "" {
		config.WorkingDir = executablePath
	}

	switch step {
	case "init":
			fmt.Println("Running initialization step")
			err := app.NewInitInventoryCSVStep(*config, logger).Init()

			if err != nil {
				logger.Fatal(fmt.Sprintf("Failed to init inventory csv: %v", err))
			}

	case "process":
			fmt.Println("Running inventory step")
			err := app.NewProcessInvetoryStep(*config, logger).Process()

			if err != nil {
				logger.Fatal(fmt.Sprintf("Failed to process inventory: %v", err))
			}
			
	default:
			logger.Fatal(fmt.Sprintf("Invalid step: %s", step))
	}

	// Keep the terminal open
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func getExecutablePath(logger utils.Logger) string {
	exePath, err := os.Executable()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to get executable path: %v", err))
	}
	return filepath.Dir(exePath)   
}