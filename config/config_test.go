package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"thwInventoryMerge/config"
	"thwInventoryMerge/utils/utilsfakes"
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
		"inventory_csv_file_name": "foo_inventory_csv_file_name",
		"columns": {
			"equipment_id": "foo_equipment_id_column_name",
			"equipment_count_actual": "foo_equipment_available_column_name",
			"equipment_count_target": "foo_equipment_target_column_name"
		}
	}
	`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name(), nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.WorkingDir).To(Equal("foo_working_dir"))
			Expect(cfg.InventoryCSVFileName).To(Equal("foo_inventory_csv_file_name"))
			Expect(cfg.Columns.EquipmentID).To(Equal("foo_equipment_id_column_name"))
			Expect(cfg.Columns.EquipmentCountActual).To(Equal("foo_equipment_available_column_name"))
			Expect(cfg.Columns.EquipmentCountTarget).To(Equal("foo_equipment_target_column_name"))
		})

		var _ = Describe("config errors", func() {
			It("returns an error if mandatory inventory_csv_file_name is missing", func() {
				jsonContent := `
	{
		"working_dir": "foo_working_dir"
	}
	`
				_, err := tempFile.Write([]byte(jsonContent))
				Expect(err).ToNot(HaveOccurred())
				tempFile.Close()

				cfg, err := config.LoadConfig(tempFile.Name(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to validate the config file, property inventory_csv_file_name is required"))
				Expect(cfg).To(BeNil())
			})
		})

		It("returns an error if mandatory columns.equipment_id is missing", func() {
			jsonContent := `
	{
		"inventory_csv_file_name": "foo_inventory_csv_file_name",
		"columns": {
		}
	}
	`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name(), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to validate the config file, property columns.equipment_id is required"))
			Expect(cfg).To(BeNil())
		})

		It("returns an error if mandatory columns.equipment_count_actual is missing", func() {
			jsonContent := `
		{
			"inventory_csv_file_name": "foo_inventory_csv_file_name",
			"columns": {
				"equipment_id": "foo_equipment_id"
			}
		}
		`
			_, err := tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			cfg, err := config.LoadConfig(tempFile.Name(), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to validate the config file, property columns.equipment_count_actual is required"))
			Expect(cfg).To(BeNil())
		})
	})

	var _ = Describe("GetCSVFilesWithRecordedEquipment", func() {
		It("should return the CSV files", func() {

			tempDir, err := os.MkdirTemp("", "test-csv")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			fileNames := []string{"file1.csv", "file2.csv", "inventory_fgr_n.csv", "file3.csv"}
			for _, fileName := range fileNames {
				filePath := filepath.Join(tempDir, fileName)
				file, err := os.Create(filePath)
				Expect(err).ToNot(HaveOccurred())
				file.Close()
			}

			jsonContent := fmt.Sprintf(`
{
	"working_dir": "%s",
	"inventory_csv_file_name": "inventory_fgr_n.csv",
	"columns": {
		"equipment_id": "foo_equipment_id_column_name",
		"equipment_count_actual": "foo_equipment_available_column_name"
	}
}
`, strings.ReplaceAll(tempDir, "\\", "\\\\"))

			fmt.Println(jsonContent)
			_, err = tempFile.Write([]byte(jsonContent))
			Expect(err).ToNot(HaveOccurred())
			tempFile.Close()

			logger := &utilsfakes.FakeLogger{}

			cfg, err := config.LoadConfig(tempFile.Name(), logger)
			Expect(err).ToNot(HaveOccurred())

			files, err := cfg.GetCSVFilesWithRecordedEquipment()
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(3))
			Expect(files).To(ContainElement(filepath.Join(tempDir, "file1.csv")))
			Expect(files).To(ContainElement(filepath.Join(tempDir, "file2.csv")))
			Expect(files).To(ContainElement(filepath.Join(tempDir, "file3.csv")))
		})
	})

})
