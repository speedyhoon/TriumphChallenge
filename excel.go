package main

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

const worksheet = "Sheet1"

func excelHeading(eventName string) (f *excelize.File, row int) {
	f = excelize.NewFile()
	row = 1
	excelStr(f, &row, "A", eventName)
	checkErr(f.MergeCell(worksheet, "A1", "M1"))
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center"},"font":{"bold":true,"size":16,"color":"#FF8800"}}`)
	checkErr(err)
	checkErr(f.SetCellStyle(worksheet, "A1", "M1", style))

	row++
	excelStr(f, &row, "A", hPosition)

	row++
	// Create worksheet column headings in the second row.
	excelStr(f, &row, "A", hPosition)
	excelStr(f, &row, "B", hRacingNumber)
	excelStr(f, &row, "C", hDriver)
	excelStr(f, &row, "D", hQualify)
	excelStr(f, &row, "E", hSeconds)
	excelStr(f, &row, "F", hFastest)
	excelStr(f, &row, "G", hSeconds)
	excelStr(f, &row, "H", hSlowest)
	excelStr(f, &row, "I", hSeconds)
	excelStr(f, &row, "J", hAverage)
	excelStr(f, &row, "K", hPercentage)
	excelStr(f, &row, "L", hRuns)
	excelStr(f, &row, "M", hLaps)

	return f, row
}

// excelRow populates spreadsheet cells.
func excelRow(f *excelize.File, d *Driver, ordinal string, row *int) {
	*row++

	excelStr(f, row, "A", ordinal)
	excelStr(f, row, "B", d.RaceNumber)
	excelStr(f, row, "C", d.Name)
	excelStr(f, row, "D", d.Qualify.String())
	excelFloat(f, row, "E", d.Qualify.Seconds())
	excelStr(f, row, "F", d.Fastest.String())
	excelFloat(f, row, "G", d.Fastest.Seconds())
	excelStr(f, row, "H", d.Slowest.String())
	excelFloat(f, row, "I", d.Slowest.Seconds())

	// Slow Average equals d.Qualify.Seconds() + d.Slowest.Seconds() / 2.
	excelFormula(f, row, "J", fmt.Sprintf("(E%d+I%[1]d)/2", *row))

	// Percentage equals d.Fastest.Seconds() / ((d.Slowest.Seconds() + d.Qualify.Seconds()) / 2) * 100.
	excelFormula(f, row, "K", fmt.Sprintf("G%d/J%[1]d * 100", *row))

	excelInt(f, row, "L", d.Runs)
	excelInt(f, row, "M", d.Laps)
}

func excelFooter(xlsx *excelize.File, spreadsheetRow *int, missingCars []string) {
	if len(missingCars) == 0 {
		return
	}

	*spreadsheetRow += 2
	checkErr(xlsx.SetCellStr(worksheet, axis(spreadsheetRow, "A"), hMissing))
	for i := range missingCars {
		*spreadsheetRow++
		checkErr(xlsx.SetCellStr(worksheet, axis(spreadsheetRow, "A"), missingCars[i]))
	}
}

func excelStr(f *excelize.File, spreadsheetRow *int, column, value string) {
	checkErr(f.SetCellStr(worksheet, axis(spreadsheetRow, column), value))
}

func excelFloat(f *excelize.File, spreadsheetRow *int, column string, value float64) {
	const bitSize = 64 // Float64 precision.
	checkErr(f.SetCellDefault(worksheet, axis(spreadsheetRow, column), strconv.FormatFloat(value, 'f', decimalPlaces, bitSize)))
}

func excelFormula(f *excelize.File, spreadsheetRow *int, column, value string) {
	checkErr(f.SetCellFormula(worksheet, axis(spreadsheetRow, column), value))
}

func excelInt(f *excelize.File, spreadsheetRow *int, column string, value uint) {
	checkErr(f.SetCellInt(worksheet, axis(spreadsheetRow, column), int(value)))
}

func axis(row *int, column string) string {
	return fmt.Sprintf("%s%d", column, *row)
}
