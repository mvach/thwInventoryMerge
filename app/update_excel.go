package app

import (
    "encoding/csv"
    "fmt"
    "os"
    "strings"
    "thwInventoryMerge/config"

    "github.com/xuri/excelize/v2"
)

type UpdateExcel interface {
    Update() error
}

type updateExcel struct {
    config config.Config
}

func NewUpdateExcel(config config.Config) UpdateExcel {
    return &updateExcel{
        config: config,
    }
}

func (u *updateExcel) Update() error {
    recordedInventory, err := u.getRecordedInventory()
    if err != nil {
        return fmt.Errorf("failed to get recorded inventory: %w", err)
    }

    excelFile, err := excelize.OpenFile(u.config.GetAbsoluteExcelFileName())
    if err != nil {
        return fmt.Errorf("failed to open the Excel file: %w", err)
    }

    rows, err := excelFile.GetRows(u.config.ExcelConfig.WorksheetName)
    if err != nil {
        return fmt.Errorf("failed to get rows from Excel file: %w", err)
    }

    for rowIndex, row := range rows {
        if rowIndex == 0 {
            err := u.getHeaderIndices(row)
            if err != nil {
                return fmt.Errorf("failed to get headers from Excel file: %w", err)
            }
            continue
        }

        // Check if the equipment ID exists in the inventory numbers
        equipmentID := row[*u.config.ExcelConfig.EquipmentIDColumnIndex]
        if count, exists := recordedInventory[strings.ToLower(equipmentID)]; exists {

            // Set the value in the specified cell
            cell, err := excelize.CoordinatesToCellName(*u.config.ExcelConfig.EquipmentAvailableColumnIndex+1, rowIndex+1)
            if err != nil {
                return fmt.Errorf("failed to get cell from coordinates: %w", err)
            }
            excelFile.SetCellValue(
                u.config.ExcelConfig.WorksheetName,
                cell,
                count,
            )
        }
    }

    if err := excelFile.Save(); err != nil {
        return fmt.Errorf("failed to save the updated Excel file: %w", err)
    } else {
        fmt.Println("\nUpdated Excel file successfully.")
    }

    return nil
}

func (u *updateExcel) getHeaderIndices(row []string) error {
    for j, col := range row {
        if strings.EqualFold(col, u.config.ExcelConfig.EquipmentIDColumnName) {
            u.config.ExcelConfig.EquipmentIDColumnIndex = &j
        } else if strings.EqualFold(col, u.config.ExcelConfig.EquipmentAvailableColumnName) {
            u.config.ExcelConfig.EquipmentAvailableColumnIndex = &j
        }
    }

    if u.config.ExcelConfig.EquipmentIDColumnIndex == nil {
        return fmt.Errorf(
            "failed to find header %s in first row of worksheet %s",
            u.config.ExcelConfig.EquipmentIDColumnName,
            u.config.ExcelConfig.WorksheetName,
        )
    }
    if u.config.ExcelConfig.EquipmentAvailableColumnIndex == nil {
        return fmt.Errorf(
            "failed to find header %s in first row of worksheet %s",
            u.config.ExcelConfig.EquipmentAvailableColumnName,
            u.config.ExcelConfig.WorksheetName,
        )
    }

    return nil
}

func (u *updateExcel) getRecordedInventory() (map[string]int, error) {
    inventoryNumbers := make(map[string]int)

    csvFiles, err := u.config.GetCSVFiles()
    if err != nil {
        return nil, fmt.Errorf("failed to get CSV files: %w", err)
    }

    for _, file := range csvFiles {
        f, err := os.Open(file)
        if err != nil {
            return nil, fmt.Errorf("failed to open CSV file: %w", err)
        }
        defer f.Close()

        reader := csv.NewReader(f)
        records, err := reader.ReadAll()
        if err != nil {
            return nil, fmt.Errorf("failed to read CSV file: %w", err)
        }

        for _, record := range records {
            if len(record) > 0 {
                inventoryNumbers[strings.ToLower(record[0])]++
            }
        }
    }

    u.printInventoryNumbers(inventoryNumbers)

    return inventoryNumbers, nil
}

func (u *updateExcel)  printInventoryNumbers(inventoryNumbers map[string]int) {
    fmt.Printf("%-20s %s\n", "Inventory Number", "Quantity")
    fmt.Println(strings.Repeat("-", 30)) // Print a separator line

    for number, quantity := range inventoryNumbers {
        fmt.Printf("%-20s %d\n", number, quantity)
    }
}
