package app

import (
	"strings"
)

type RecordedInventoryMap map[string]int

type RecordedInventory interface {
	AsMap() (RecordedInventoryMap, error)
}

type recordedInventory struct {
	data []CSVContent
}

func NewRecordedInventory(data []CSVContent) RecordedInventory {
	return recordedInventory{
		data: data,
	}
}

func (r recordedInventory) AsMap() (RecordedInventoryMap, error) {
	inventoryNumbers := make(RecordedInventoryMap)

	for _, csvContent := range r.data {

		for _, record := range csvContent {
			if len(record) > 0 {
				inventoryNumbers[strings.ToLower(record[0])]++
			}
		}
	}

	return inventoryNumbers, nil
}