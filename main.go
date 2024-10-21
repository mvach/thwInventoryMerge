package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"thwInventoryMerge/app"
	"thwInventoryMerge/config"
)

func main() {
    var configPath string
    
    flag.StringVar(&configPath, "c", "config.json", "the config file path")
    flag.Parse()

    config, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load config from path %s: %v", configPath, err)
    }

    if config.WorkingDir == "" {
        config.WorkingDir = getExecutablePath()
    }

    if configPath == "" {
        filepath.Join(config.WorkingDir, "config.json")
    }

    err = app.NewUpdateExcel(*config).Update()
    if err != nil {
        log.Fatalf("Failed to update Excel: %v", err)
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