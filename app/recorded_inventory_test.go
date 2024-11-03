package app_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"thwInventoryMerge/app"
)

var _ = Describe("RecordedInventory", func() {
	var (
		csvFile1 *os.File
		csvFile2 *os.File
		err      error
	)

	BeforeEach(func() {
		csvFile1, err = os.CreateTemp("", "csv1.csv")
		Expect(err).ToNot(HaveOccurred())

		csvFile2, err = os.CreateTemp("", "csv2.csv")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(csvFile1.Name())
		os.Remove(csvFile2.Name())
	})

	var _ = Describe("GetRecordedInventory", func() {
		It("returns the inventory recorded in multiple csv", func() {
			recordedInventory := app.NewRecordedInventory(
				[]app.CSVContent{[][]string{
					{"0001-S001304"},
					{"0509-002494"},
					{"0591-S002360"},
					{"0509-002494"},
				}, [][]string{
					{"0591-002781"},
					{"0591-S002319"},
					{"0591-002781"},
					{"0591-S002319"},
					{"0591-002781"},
				}})
			inventoryMap, err := recordedInventory.AsMap()
			Expect(err).ToNot(HaveOccurred())
			Expect(inventoryMap).To(HaveLen(5))
			Expect(inventoryMap).To(HaveKeyWithValue("0001-s001304", 1))
			Expect(inventoryMap).To(HaveKeyWithValue("0509-002494", 2))
			Expect(inventoryMap).To(HaveKeyWithValue("0591-s002360", 1))
			Expect(inventoryMap).To(HaveKeyWithValue("0591-002781", 3))
			Expect(inventoryMap).To(HaveKeyWithValue("0591-s002319", 2))
		})
	})
})
