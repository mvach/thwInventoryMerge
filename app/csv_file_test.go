package app_test

import (
	"bufio"
	"encoding/csv"
	"os"
	"path/filepath"
	"thwInventoryMerge/app"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSVFile", func() {

	var (
		filePath string
	)

	BeforeEach(func() {
	})

	AfterEach(func() {
		if filePath != "" {
			os.Remove(filePath)
		}
	})

	var _ = Describe("Read", func() {
		Context("when the CSV file is valid", func() {
			BeforeEach(func() {
				filePath = filepath.Join(os.TempDir(), "valid.csv")
				file, err := os.Create(filePath)
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				writer := csv.NewWriter(file)
				writer.Comma = ';'
				writer.WriteAll([][]string{
					{"name", "age", "city"},
					{"Alice", "30", "New York"},
					{"Bob", "25", "San Francisco"},
				})
				writer.Flush()
			})

			It("should read the CSV file successfully", func() {
				content, err := app.NewCSVFile().Read(filePath)
				Expect(err).NotTo(HaveOccurred())

				Expect(content).To(HaveLen(3)) // header + 2 rows
				Expect(content[0]).To(Equal([]string{"name", "age", "city"}))
				Expect(content[1]).To(Equal([]string{"Alice", "30", "New York"}))
				Expect(content[2]).To(Equal([]string{"Bob", "25", "San Francisco"}))
			})
		})

		Context("when the CSV file starts like the original from thw", func() {
			BeforeEach(func() {
				csvContent := `"Ebene";"OE";"Art";"FB";"Menge";"Menge Ist";"Verfügbar";"Ausstattung | Hersteller | Typ";"Sachnummer";"Inventar Nr";"Gerätenr.";"Status"
;"OV Speyer";"";"";"";"";"";"1. Technischer Zug/Fachgruppe Notversorgung und Notinstandsetzung";"";"";"";"V"
1;"";"";"";"";"";"";"Geringwertiges Material";"";"";"";"V"
2;"";"Gwm";"";"";"1";"1";"  Eiskratzer, handelsüblich";"2540T21171";"---------------";"--------------------";"V"`

				filePath = filepath.Join(os.TempDir(), "valid.csv")
				file, err := os.Create(filePath)
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				writer := bufio.NewWriter(file)
				_, err = writer.WriteString(csvContent)
				Expect(err).NotTo(HaveOccurred())

				err = writer.Flush()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should read the CSV file successfully", func() {
				_, err := app.NewCSVFile().Read(filePath)
				Expect(err).NotTo(HaveOccurred())

				// Expect(content).To(HaveLen(3)) // header + 2 rows
				// Expect(content[0]).To(Equal([]string{"name", "age", "city"}))
				// Expect(content[1]).To(Equal([]string{"Alice", "30", "New York"}))
				// Expect(content[2]).To(Equal([]string{"Bob", "25", "San Francisco"}))
			})
		})



		Context("when the file does not exist", func() {
			It("should return an error", func() {
				_, err := app.NewCSVFile().Read("nonexistent.csv")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to open CSV file"))
			})
		})

		Context("when the CSV file is malformed", func() {
			BeforeEach(func() {
				// Create a temporary CSV file with malformed content
				filePath = filepath.Join(os.TempDir(), "malformed.csv")
				file, err := os.Create(filePath)
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				file.WriteString("name;age;city\nAlice;30\nBob;25;San Francisco") // Missing one field in Alice's row
			})

			It("should return an error when reading", func() {
				_, err := app.NewCSVFile().Read(filePath)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to read CSV file"))
			})
		})
	})

	var _ = Describe("Write", func() {
		var (
			content  [][]string
		)

		BeforeEach(func() {
			content = [][]string{
				{"name", "age", "city"},
				{"Alice", "30", "New York"},
				{"Bob", "25", "San Francisco"},
			}
		})

		Context("when the CSV file path is valid", func() {
			BeforeEach(func() {
				filePath = filepath.Join(os.TempDir(), "output.csv")
			})

			It("should write the content to the CSV file successfully", func() {
				err := app.NewCSVFile().Write(filePath, content)
				Expect(err).NotTo(HaveOccurred())

				file, err := os.Open(filePath)
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				scanner := bufio.NewScanner(file)

				// \ufeff is the UTF-8 BOM for Excel on Windows compatibility
				expectedLines := []string{
					"\ufeffname;age;city",
					"Alice;30;New York",
					"Bob;25;San Francisco",
				}

				i := 0
				for scanner.Scan() {
					Expect(scanner.Text()).To(Equal(expectedLines[i]))
					i++
				}
				Expect(scanner.Err()).NotTo(HaveOccurred())
				Expect(i).To(Equal(len(expectedLines)))
			})
		})

		Context("when the file path is invalid", func() {
			It("should return an error", func() {
				invalidPath := "/invalid/output.csv" // Likely to be invalid on most systems
				err := app.NewCSVFile().Write(invalidPath, content)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to create CSV file"))
			})
		})
	})
})
