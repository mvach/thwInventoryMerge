package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"thwInventoryMerge/config"
)

var _ = Describe("Config", func() {
	var (
		tempFile *os.File
		err      error
	)

	BeforeEach(func() {
		tempFile, err = os.CreateTemp("", "config.yml")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(tempFile.Name())
	})

	var _ = Describe("LoadConfig", func() {

		It("should load the configuration", func() {
			jsonContent := `
{
    "working_dir": "foo_working_dir",
    "excel_file_name": "foo_excel_file_name",
    "excel_config": {
        "worksheet_name": "foo_worksheet_name",
        "equipment_id_column_name": "foo_equipment_id_column_name",
        "equipment_available_column_name": "foo_equipment_available_column_name",
        "equipment_available_value": "foo_equipment_available_value"
    }
}
`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.WorkingDir).To(Equal("foo_working_dir"))
			Expect(cfg.ExcelFileName).To(Equal("foo_excel_file_name"))
			Expect(cfg.ExcelConfig.WorksheetName).To(Equal("foo_worksheet_name"))
			Expect(cfg.ExcelConfig.EquipmentIDColumnName).To(Equal("foo_equipment_id_column_name"))
			Expect(cfg.ExcelConfig.EquipmentAvailableColumnName).To(Equal("foo_equipment_available_column_name"))
		})

		var _ = Describe("config errors", func() {
			It("returns an error if mandatory working_dir is missing", func() {
				jsonContent := `
{
}
`
				_, err := tempFile.Write([]byte(jsonContent))
				Expect(err).ToNot(HaveOccurred())
				tempFile.Close()

				cfg, err := config.LoadConfig(tempFile.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to validate the config file, property working_dir is required"))
				Expect(cfg).To(BeNil())
			})

			It("returns an error if mandatory excel_file_name is missing", func() {
				jsonContent := `
{
    "working_dir": "foo_working_dir"
}
`
				_, err := tempFile.Write([]byte(jsonContent))
				Expect(err).ToNot(HaveOccurred())
				tempFile.Close()

				cfg, err := config.LoadConfig(tempFile.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to validate the config file, property excel_file_name is required"))
				Expect(cfg).To(BeNil())
			})
		})

		It("returns an error if mandatory excel_config,worksheet_name", func() {
			jsonContent := `
{
    "working_dir": "foo_working_dir",
    "excel_file_name": "foo_excel_file_name"
}
`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to validate the config file, property excel_config.worksheet_name is required"))
			Expect(cfg).To(BeNil())
		})

		It("returns an error if mandatory excel_config,equipment_id_column_name", func() {
			jsonContent := `
{
    "working_dir": "foo_working_dir",
    "excel_file_name": "foo_excel_file_name",
    "excel_config": {
        "worksheet_name": "foo_worksheet_name"
    }
}
`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to validate the config file, property excel_config.equipment_id_column_name is required"))
			Expect(cfg).To(BeNil())
		})

	})

	It("returns an error if mandatory excel_config,equipment_available_column_name", func() {
		jsonContent := `
{
    "working_dir": "foo_working_dir",
    "excel_file_name": "foo_excel_file_name",
    "excel_config": {
        "worksheet_name": "foo_worksheet_name",
        "equipment_id_column_name": "foo_equipment_id_column_name"
    }
}
`
		_, err := tempFile.Write([]byte(jsonContent))
		Expect(err).ToNot(HaveOccurred())
		tempFile.Close()

		cfg, err := config.LoadConfig(tempFile.Name())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("failed to validate the config file, property excel_config.equipment_available_column_name is required"))
		Expect(cfg).To(BeNil())
	})
})
