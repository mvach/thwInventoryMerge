package app_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"runtime"
	"thwInventoryMerge/app"
	"thwInventoryMerge/utils/utilsfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var _ = Describe("EncodingProvider", func() {

	var (
		filePath string
		logger   *utilsfakes.FakeLogger
	)

	BeforeEach(func() {
		logger = &utilsfakes.FakeLogger{}
	})

	AfterEach(func() {
		if filePath != "" {
			os.Remove(filePath)
		}
	})

	var _ = Describe("GetFileEncoding", func() {

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

		It("should return utf-8 encoding", func() {
			// Get the directory of the current file
			_, currentFile, _, _ := runtime.Caller(0)

			filePath := filepath.Join(currentFile, "..", "..", "testdata", "app", "utf-8.csv")

			enc, err := app.NewEncodingProvider(logger).GetFileEncoding(filePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(enc).To(Equal(unicode.UTF8))
		})

		It("should return iso8859_1 encoding", func() {
			// Get the directory of the current file
			_, currentFile, _, _ := runtime.Caller(0)

			filePath := filepath.Join(currentFile, "..", "..", "testdata", "app", "iso-8859-1.csv")

			enc, err := app.NewEncodingProvider(logger).GetFileEncoding(filePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(enc).To(Equal(charmap.ISO8859_1))
		})

		It("should return an error if the encoding is unknown", func() {
			// Get the directory of the current file
			_, currentFile, _, _ := runtime.Caller(0)
			 
		 filePath := filepath.Join(currentFile, "..", "..", "testdata", "app", "cp437.csv")

		 _, err := app.NewEncodingProvider(logger).GetFileEncoding(filePath)

		 Expect(err.Error()).To(ContainSubstring("unsupported encoding"))
	 })

	})
})
