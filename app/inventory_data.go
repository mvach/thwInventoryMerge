package app

import (
	"fmt"
	"strconv"
	"strings"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils"
)

type csvHeader map[string]int
type csvHeaderReverse map[int]string
type csvContent []map[string]string

type InventoryData interface {
	GetContent() [][]string

	UpdateInventory(recordedInventory RecordedInventoryMap) error

	GeneratePsydoEquipmentIDs() error
}

type inventoryData struct {
	csvHeader        csvHeader
	csvHeaderReverse csvHeaderReverse
	content          csvContent
	config           config.Config
	logger           utils.Logger
}

func NewInventoryData(data [][]string, config config.Config, logger utils.Logger) (InventoryData, error) {

	csvHeader := make(csvHeader)

	var content csvContent
	for i := 0; i < len(data); i++ {
		record := data[i]

		// create header index on first row
		if i == 0 {
			for i, colName := range record {
				csvHeader[colName] = i
			}
		}

		row := make(map[string]string)
		for colName, colIndex := range csvHeader {
			if colIndex < len(record) {
				row[colName] = record[colIndex]
			}
		}
		content = append(content, row)
	}

	csvHeaderReverse := make(csvHeaderReverse)
	for key, value := range csvHeader {
		csvHeaderReverse[value] = key
	}

	return &inventoryData{
		csvHeader:        csvHeader,
		csvHeaderReverse: csvHeaderReverse,
		content:          content,
		config:           config,
		logger:           logger,
	}, nil
}

func (c *inventoryData) GetContent() [][]string {
	var result [][]string

	for _, row := range c.content {
		var resultRow []string

		for i := 0; i < len(c.csvHeaderReverse); i++ {
			resultRow = append(resultRow, row[c.csvHeaderReverse[i]])
		}

		result = append(result, resultRow)
	}

	return result
}

func (c *inventoryData) UpdateInventory(recordedInventory RecordedInventoryMap) error {

	firstEquipment := true

	for inventory, amount := range recordedInventory {
		inventoryFound := false
		actualValue := strconv.Itoa(amount)

		for _, row := range c.content {
			configColumns := c.config.Columns

			// ignore case comparison
			if strings.EqualFold(row[configColumns.EquipmentID], inventory) {
				inventoryFound = true

				if configColumns.EquipmentCountTarget != "" {
					targetValueInt, err := strconv.Atoi(row[configColumns.EquipmentCountTarget])
					if err != nil {
						return fmt.Errorf("error converting target value to int: %v", err)
					}
					if amount >= targetValueInt {
						amount = amount - targetValueInt
						actualValue = strconv.Itoa(targetValueInt)
					}
				}

				row[configColumns.EquipmentCountActual] = actualValue
				actualValue = strconv.Itoa(amount);
			}
		}

		if !inventoryFound {
			if firstEquipment {
				c.logger.Info("recorded equipment not available in the inventory:")
				c.logger.Info("")
				c.logger.WarnIndented("equipment     : amount")
				c.logger.WarnIndented("----------------------")
				firstEquipment = false
			}

			c.logger.WarnIndented(fmt.Sprintf("%-13s : %5d", inventory, amount))
		}
	}

	if !firstEquipment {
		c.logger.Info("")
	}

	return nil
}

func (c *inventoryData) GeneratePsydoEquipmentIDs() error {

	content := c.content
	columns := c.config.Columns
	// Forward iteration
	for i := range content {
		if !utils.StartsWithNumber(content[i][columns.EquipmentID]) {
			equipmentLayer, err := strconv.Atoi(content[i][columns.EquipmentLayer])
			if err != nil && i > 0 {
				c.logger.Warn(fmt.Sprintf("failed to convert column '%s' to number on line %d", columns.EquipmentLayer, i+1))
				continue
			}

			searchedEquipmentLayer := equipmentLayer - 1
			searchPath := fmt.Sprintf("%d", i+1)

			// iterate backwards to find the last equipment number in upper layers
			for j := i - 1; j >= 0; j-- {
				if searchedEquipmentLayer <= 0 {
					msg := fmt.Sprintf(
						"skipping ID generation for line %d (processed lines %s). Could not find a '%s' value up to '%s' 1",
						i+1,
						searchPath,
						columns.EquipmentID,
						columns.EquipmentLayer,
					)
					c.logger.Warn(msg)

					break
				}

				previousLineEquipmentLayer, err := strconv.Atoi(content[j][columns.EquipmentLayer])
				if err != nil {
					searchPath = searchPath + fmt.Sprintf(", %d", j+1)

					msg := fmt.Sprintf(
						"skipping ID generation for line %d (processed lines %s). Column '%s' of line %d cannot be converted to number",
						i+1,
						searchPath,
						columns.EquipmentLayer,
						j+1,
					)
					c.logger.Warn(msg)

					break
				}

				if previousLineEquipmentLayer == searchedEquipmentLayer {
					searchPath = searchPath + fmt.Sprintf(", %d", j+1)

					if utils.StartsWithNumber(content[j][columns.EquipmentID]) && !strings.Contains(content[j][columns.EquipmentID], "__") {
						content[i][columns.EquipmentID] = content[j][columns.EquipmentID] + "__" + content[i][columns.EquipmentPartNumber]

						msg := fmt.Sprintf("created ID for line %d (processed lines %s)", i+1, searchPath)
						c.logger.Info(msg)

						break
					} else {
						searchedEquipmentLayer = searchedEquipmentLayer - 1
					}
				}
			}
		}
	}

	return nil
}
