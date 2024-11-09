package app

import (
	"fmt"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils"
)

type InitInventoryCSVStep interface {
	Init() error
}

type initInventoryCSVStep struct {
	config config.Config
	logger utils.Logger
}

func NewInitInventoryCSVStep(config config.Config, logger utils.Logger) InitInventoryCSVStep {
	return &initInventoryCSVStep{
		config: config,
		logger: logger,
	}
}

func (s *initInventoryCSVStep) Init() error {

  csvFile := NewCSVFile()

  content, err := csvFile.Read(s.config.GetAbsoluteInventoryCSVFileName())
	if err != nil {
		return fmt.Errorf("failed to read CSV file '%s': %v", s.config.GetAbsoluteInventoryCSVFileName(), err)
	}

  s.addActualEquipmentColumn(content)

	inventoryData, err := NewInventoryData(content, s.config, s.logger)
	if err != nil {
		return fmt.Errorf("failed to init inventory data: %v", err)
	}

	inventoryData.GeneratePsydoEquipmentIDs()

	inventoryData.GetContent()

  csvFile.Write(s.config.GetAbsoluteInventoryCSVFileName(), inventoryData.GetContent())

  return nil
}

func (s *initInventoryCSVStep) addActualEquipmentColumn(content CSVContent) {
	for i := range content {
		if i == 0 {
			content[i] = append(content[i], s.config.Columns.EquipmentCountActual)
		} else {
			content[i] = append(content[i], "")
		}
	}
}