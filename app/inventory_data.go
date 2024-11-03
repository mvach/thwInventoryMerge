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

		for _, row := range c.content {
			configColumns := c.config.Columns

			// ignore case comparison
			if strings.EqualFold(row[configColumns.EquipmentID], inventory) {
				inventoryFound = true

				actualValue := strconv.Itoa(amount)

				if configColumns.EquipmentCountTarget != "" {
					targetValueInt, err := strconv.Atoi(row[configColumns.EquipmentCountTarget])
					if err != nil {
						return fmt.Errorf("error converting target value to int: %v", err)
					}
					if amount > targetValueInt {
						actualValue = strconv.Itoa(targetValueInt)
					}
				}

				row[configColumns.EquipmentCountActual] = actualValue
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
