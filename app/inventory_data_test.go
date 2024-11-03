package app_test

import (
	"thwInventoryMerge/app"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils/utilsfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSVFile", func() {
	var _ = Describe("GetContent", func() {
		It("returns the csv content", func() {
			csvData := [][]string{
				{"Verfügbar", "Ausstattung", "Inventar Nr", "Status"},
				{"4", "Handlampe", "0591-S00001", "V"},
				{"1", "Fuchsschwanz", "", "V"},
				{"1", "Rettungsweste", "0591-S00002", "U"}}

			data, err := app.NewInventoryData(csvData, config.Config{}, nil)
			Expect(err).ToNot(HaveOccurred())

			content := data.GetContent()

			Expect(content[0][0]).To(Equal("Verfügbar"))
			Expect(content[0][1]).To(Equal("Ausstattung"))
			Expect(content[0][2]).To(Equal("Inventar Nr"))
			Expect(content[0][3]).To(Equal("Status"))
			Expect(content[1][0]).To(Equal("4"))
			Expect(content[1][1]).To(Equal("Handlampe"))
			Expect(content[1][2]).To(Equal("0591-S00001"))
			Expect(content[1][3]).To(Equal("V"))
			Expect(content[2][0]).To(Equal("1"))
			Expect(content[2][1]).To(Equal("Fuchsschwanz"))
			Expect(content[2][2]).To(Equal(""))
			Expect(content[2][3]).To(Equal("V"))
			Expect(content[3][0]).To(Equal("1"))
			Expect(content[3][1]).To(Equal("Rettungsweste"))
			Expect(content[3][2]).To(Equal("0591-S00002"))
			Expect(content[3][3]).To(Equal("U"))
		})

		It("preserve leading and trailing spaces", func() {
			csvData := [][]string{
				{"Verfügbar", "Ausstattung", "Inventar Nr", "Status"},
				{"4", "  Handlampe", " 0591-S00001 ", "V  "}}

			data, err := app.NewInventoryData(csvData, config.Config{}, nil)
			Expect(err).ToNot(HaveOccurred())

			content := data.GetContent()

			Expect(content[0][0]).To(Equal("Verfügbar"))
			Expect(content[0][1]).To(Equal("Ausstattung"))
			Expect(content[0][2]).To(Equal("Inventar Nr"))
			Expect(content[0][3]).To(Equal("Status"))
			Expect(content[1][0]).To(Equal("4"))
			Expect(content[1][1]).To(Equal("  Handlampe"))
			Expect(content[1][2]).To(Equal(" 0591-S00001 "))
			Expect(content[1][3]).To(Equal("V  "))
		})
	})

	var _ = Describe("UpdateInventory", func() {
		It("updated the content", func() {
			csvData := [][]string{
				{"Verfügbar", "Ausstattung", "Inventar Nr", "Status"},
				{"4", "Handlampe", "0591-S00001", "V"},
				{"1", "Fuchsschwanz", "1234", "V"},
				{"1", "Rettungsweste", "0591-S00002", "U"}}

			data, err := app.NewInventoryData(csvData, config.Config{
				Columns: config.ConfigColumns{
					EquipmentID:          "Inventar Nr",
					EquipmentCountActual: "Verfügbar",
				},
			}, &utilsfakes.FakeLogger{})
			Expect(err).ToNot(HaveOccurred())

			data.UpdateInventory(app.RecordedInventoryMap{
				"0591-S00001": 100,
				"1234":        101,
				"0591-S00002": 0,
			})

			content := data.GetContent()

			Expect(content[0][0]).To(Equal("Verfügbar"))
			Expect(content[0][1]).To(Equal("Ausstattung"))
			Expect(content[0][2]).To(Equal("Inventar Nr"))
			Expect(content[0][3]).To(Equal("Status"))
			Expect(content[1][0]).To(Equal("100"))
			Expect(content[1][1]).To(Equal("Handlampe"))
			Expect(content[1][2]).To(Equal("0591-S00001"))
			Expect(content[1][3]).To(Equal("V"))
			Expect(content[2][0]).To(Equal("101"))
			Expect(content[2][1]).To(Equal("Fuchsschwanz"))
			Expect(content[2][2]).To(Equal("1234"))
			Expect(content[2][3]).To(Equal("V"))
			Expect(content[3][0]).To(Equal("0"))
			Expect(content[3][1]).To(Equal("Rettungsweste"))
			Expect(content[3][2]).To(Equal("0591-S00002"))
			Expect(content[3][3]).To(Equal("U"))
		})

		It("cuts the recorded values down to EquipmentCountTarget if provided", func() {
			csvData := [][]string{
				{"Verfügbar", "Menge", "Ausstattung", "Inventar Nr", "Status"},
				{"0", "50", "Handlampe", "0591-S00001", "V"},
			}

			data, err := app.NewInventoryData(csvData, config.Config{
				Columns: config.ConfigColumns{
					EquipmentID:          "Inventar Nr",
					EquipmentCountActual: "Verfügbar",
					EquipmentCountTarget: "Menge",
				},
			}, &utilsfakes.FakeLogger{})
			Expect(err).ToNot(HaveOccurred())

			data.UpdateInventory(app.RecordedInventoryMap{
				"0591-S00001": 100,
			})

			content := data.GetContent()

			Expect(content[0][0]).To(Equal("Verfügbar"))
			Expect(content[0][1]).To(Equal("Menge"))
			Expect(content[0][2]).To(Equal("Ausstattung"))
			Expect(content[0][3]).To(Equal("Inventar Nr"))
			Expect(content[0][4]).To(Equal("Status"))
			Expect(content[1][0]).To(Equal("50"))
			Expect(content[1][1]).To(Equal("50"))
			Expect(content[1][2]).To(Equal("Handlampe"))
			Expect(content[1][3]).To(Equal("0591-S00001"))
			Expect(content[1][4]).To(Equal("V"))
		})

		It("logs not existing equipment", func() {
			logger := &utilsfakes.FakeLogger{}

			csvData := [][]string{
				{"Verfügbar", "Ausstattung", "Inventar Nr", "Status"}}

			data, err := app.NewInventoryData(csvData, config.Config{}, logger)
			Expect(err).ToNot(HaveOccurred())

			data.UpdateInventory(app.RecordedInventoryMap{
				"not_existing": 1,
			})

			Expect(logger.WarnIndentedCallCount()).To(Equal(3))
			Expect(logger.WarnIndentedArgsForCall(2)).To(Equal("not_existing  :     1"))
		})
	})
})
