package app_test

import (
	"thwInventoryMerge/app"
	"thwInventoryMerge/config"
	"thwInventoryMerge/utils/utilsfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSVFile", func() {

	var (
		logger *utilsfakes.FakeLogger
	)

	BeforeEach(func() {
		logger = &utilsfakes.FakeLogger{}
	})

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

	var _ = Describe("GeneratePsydoEquipmentIDs", func() {

		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg = config.Config{
				Columns: config.ConfigColumns{
					EquipmentLayer:       "Ebene",
					EquipmentPartNumber:  "Sachnummer",
					EquipmentID:          "Inventar Nr",
					EquipmentCountActual: "Bestand IST",
				},
			}
		})

		var _ = Describe("when equipmentID does not start with a number", func() {
			It("skips ID generation on rows where layer is not a number ", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"", "Handlampe", "0591-S00001", ""},
					{"foo", "Fuchsschwanz", "1234", "---"}}

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[0][3]).To(Equal("Inventar Nr"))
				Expect(content[1][3]).To(Equal(""))
				Expect(content[2][3]).To(Equal("---"))

				Expect(logger.WarnCallCount()).To(Equal(2))
				Expect(logger.WarnArgsForCall(0)).To(Equal("failed to convert column 'Ebene' to number on line 2"))
				Expect(logger.WarnArgsForCall(1)).To(Equal("failed to convert column 'Ebene' to number on line 3"))
			})

			It("skips ID generation on rows where layer 1", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"1", "Handlampe", "0591-S00001", ""},
					{"1", "Fuchsschwanz", "1234", "---"}}

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal(""))
				Expect(content[2][3]).To(Equal("---"))
			})

			It("skips ID generation on rows where the upper line layer is not a number", func() {
			csvData := [][]string{
				{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
				{"not_a_number", "Werkzugkasten", "1234", ""},
				{"2", "Ratschenkasten", "1234", ""},
				{"3", "Hammer", "1234", ""}}

			inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
			Expect(err).ToNot(HaveOccurred())
			err = inventoryData.GeneratePsydoEquipmentIDs()
			Expect(err).ToNot(HaveOccurred())

			content := inventoryData.GetContent()

			Expect(content[1][3]).To(Equal(""))
			Expect(content[2][3]).To(Equal(""))
			Expect(content[3][3]).To(Equal(""))

			Expect(logger.WarnCallCount()).To(Equal(3))
			Expect(logger.WarnArgsForCall(0)).To(Equal("failed to convert column 'Ebene' to number on line 2"))
			Expect(logger.WarnArgsForCall(1)).To(Equal("skipping ID generation for line 3 (processed lines 3, 2). Column 'Ebene' of line 2 cannot be converted to number"))
			Expect(logger.WarnArgsForCall(2)).To(Equal("skipping ID generation for line 4 (processed lines 4, 3, 2). Column 'Ebene' of line 2 cannot be converted to number"))
			})

			It("skips ID generation on rows where no equipment id can be found on higher layers", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"1", "Werkzugkasten", "1234", ""},
					{"2", "Ratschenkasten", "1234", ""},
					{"3", "Hammer", "1234", ""}}

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal(""))
				Expect(content[2][3]).To(Equal(""))
				Expect(content[3][3]).To(Equal(""))

				Expect(logger.WarnCallCount()).To(Equal(3))
				Expect(logger.WarnArgsForCall(0)).To(Equal("skipping ID generation for line 2 (processed lines 2). Could not find a 'Inventar Nr' value up to 'Ebene' 1"))
				Expect(logger.WarnArgsForCall(1)).To(Equal("skipping ID generation for line 3 (processed lines 3, 2). Could not find a 'Inventar Nr' value up to 'Ebene' 1"))
				Expect(logger.WarnArgsForCall(2)).To(Equal("skipping ID generation for line 4 (processed lines 4, 3, 2). Could not find a 'Inventar Nr' value up to 'Ebene' 1"))
			})
		})

		var _ = Describe("when equipmentID start with a number", func() {
			It("generates an ID by using a higher equipment ID", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"1", "Werkzugkasten", "1234", "5555"},
					{"2", "Ratschenkasten", "2222", ""}}

					inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
					Expect(err).ToNot(HaveOccurred())
					err = inventoryData.GeneratePsydoEquipmentIDs()
					Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal("5555"))
				Expect(content[2][3]).To(Equal("5555__2222"))

				Expect(logger.InfoCallCount()).To(Equal(1))
				Expect(logger.InfoArgsForCall(0)).To(Equal("created ID for line 3 (processed lines 3, 2)"))
			})

			It("generates an ID by using a higher equipment ID jumping over irrelevant layers", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"1", "Werkzugkasten", "1111", "5678"},
					{"2", "Ratschenkasten", "2222", "3456"},
					{"3", "irrelevent Einsatz1", "3333", ""},
					{"4", "irrelevant Ratsche", "4444", "7654"},
					{"4", "irrelevant Verlängerung", "5555", ""},
					{"3", "Einsatz2", "3333", ""}}

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal("5678"))
				Expect(content[2][3]).To(Equal("3456"))
				Expect(content[3][3]).To(Equal("3456__3333"))
				Expect(content[4][3]).To(Equal("7654"))
				Expect(content[5][3]).To(Equal("3456__5555"))
				Expect(content[6][3]).To(Equal("3456__3333"))

				Expect(logger.InfoCallCount()).To(Equal(3))
				Expect(logger.InfoArgsForCall(0)).To(Equal("created ID for line 4 (processed lines 4, 3)"))
				Expect(logger.InfoArgsForCall(1)).To(Equal("created ID for line 6 (processed lines 6, 4, 3)"))
				Expect(logger.InfoArgsForCall(2)).To(Equal("created ID for line 7 (processed lines 7, 3)"))
			})

			It("generates an ID on inventory data with multiple Layer 1 Elements", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"},
					{"1", "Werkzugkasten1", "1111", "3456"},
					{"2", "Ratschenkasten1", "2222", ""},
					{"1", "Werkzugkasten2", "3333", "5678"},
					{"2", "Ratschenkasten2", "4444", ""}}

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal("3456"))
				Expect(content[2][3]).To(Equal("3456__2222"))
				Expect(content[3][3]).To(Equal("5678"))
				Expect(content[4][3]).To(Equal("5678__4444"))

				Expect(logger.InfoCallCount()).To(Equal(2))
				Expect(logger.InfoArgsForCall(0)).To(Equal("created ID for line 3 (processed lines 3, 2)"))
				Expect(logger.InfoArgsForCall(1)).To(Equal("created ID for line 5 (processed lines 5, 4)"))
			})

			It("generates an ID on a deep tree", func() {
				csvData := [][]string{
					{"Ebene", "Ausstattung", "Sachnummer", "Inventar Nr"}, // 1
					{"1", "Werkzugkasten", "1111", "5678"},                // 2
					{"2", "Ratschenkasten", "2222", "3456"},               // 3
					{"3", "Einsatz1", "3333", ""},                         // 4
					{"4", "Ratsche", "4444", "7654"},                      // 5
					{"4", "Verlängerung", "5555", ""},                     // 6
					{"3", "Einsatz2", "3333", ""},                         // 7
					{"4", "Nuss1", "6666", ""},                            // 8
					{"4", "Nuss2", "7777", ""},                            // 9
					{"4", "Bitset", "8888", ""},                           // 10
					{"5", "PX10", "9999", ""}}                             // 11

				inventoryData, err := app.NewInventoryData(csvData, cfg, logger)
				Expect(err).ToNot(HaveOccurred())
				err = inventoryData.GeneratePsydoEquipmentIDs()
				Expect(err).ToNot(HaveOccurred())

				content := inventoryData.GetContent()

				Expect(content[1][3]).To(Equal("5678"))
				Expect(content[2][3]).To(Equal("3456"))
				Expect(content[3][3]).To(Equal("3456__3333"))
				Expect(content[4][3]).To(Equal("7654"))
				Expect(content[5][3]).To(Equal("3456__5555"))
				Expect(content[6][3]).To(Equal("3456__3333"))
				Expect(content[7][3]).To(Equal("3456__6666"))
				Expect(content[8][3]).To(Equal("3456__7777"))
				Expect(content[9][3]).To(Equal("3456__8888"))
				Expect(content[10][3]).To(Equal("3456__9999"))

				Expect(logger.InfoCallCount()).To(Equal(7))
				Expect(logger.InfoArgsForCall(0)).To(Equal("created ID for line 4 (processed lines 4, 3)"))
				Expect(logger.InfoArgsForCall(1)).To(Equal("created ID for line 6 (processed lines 6, 4, 3)"))
				Expect(logger.InfoArgsForCall(2)).To(Equal("created ID for line 7 (processed lines 7, 3)"))
				Expect(logger.InfoArgsForCall(3)).To(Equal("created ID for line 8 (processed lines 8, 7, 3)"))
				Expect(logger.InfoArgsForCall(4)).To(Equal("created ID for line 9 (processed lines 9, 7, 3)"))
				Expect(logger.InfoArgsForCall(5)).To(Equal("created ID for line 10 (processed lines 10, 7, 3)"))
				Expect(logger.InfoArgsForCall(6)).To(Equal("created ID for line 11 (processed lines 11, 10, 7, 3)"))
			})
		})
	})
})
