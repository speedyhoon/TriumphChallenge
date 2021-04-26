package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	// Column headings
	hPosition     = "Pos"
	hRacingNumber = "#"
	hDriver       = "Driver"
	hQualify      = "Qualify"
	hFastest      = "Fastest"
	hSlowest      = "Slowest"
	hAverage      = "Slow Ave"
	hPercentage   = "Percentage"
	hRuns         = "Runs"
	hLaps         = "Laps"
	hSeconds      = "Secs"
	hMissing      = "Missing:"
	hCompetitors  = "Competitors:"
)

func render(drivers []Driver, eventName string, missingCars []string, longestNameLen uint) {
	l := uint(len(drivers))
	excel, spreadsheetRow := excelHeading(eventName)
	html := htmlHeading(eventName, l)
	txt := txtHeading(eventName, l, longestNameLen)

	var position int
	var isEqual bool
	for i := range drivers {
		// Calculate if the next or previous competitor had an identical score
		if i >= 1 && drivers[i].Percentage == drivers[i-1].Percentage && drivers[i].Runs == drivers[i-1].Runs && drivers[i].Laps == drivers[i-1].Laps {
			isEqual = true
		} else {
			position++

			isEqual = i+1 < len(drivers) && drivers[i].Percentage == drivers[i+1].Percentage && drivers[i].Runs == drivers[i+1].Runs && drivers[i].Laps == drivers[i+1].Laps
		}

		ord := ordinal(position, isEqual)

		htmlRow(html, &drivers[i], ord)
		textRow(txt, &drivers[i], ord, longestNameLen)
		excelRow(excel, &drivers[i], ord, &spreadsheetRow)
	}

	htmlFooter(html, missingCars)
	textFooter(txt, missingCars)
	excelFooter(excel, &spreadsheetRow, missingCars)

	// Print text output to screen
	fmt.Println(txt.String())

	fileName := time.Now().Format("results-2006-01-02 15;04;06")

	checkErr(ioutil.WriteFile(fileName+".txt", txt.Bytes(), filePermission))
	checkErr(ioutil.WriteFile(fileName+".html", html.Bytes(), filePermission))
	checkErr(excel.SaveAs(fileName + ".xlsx"))
}

// Ordinal gives you the input number in a rank/ordinal format.
// Ordinal(3, true) -> "=3rd"
func ordinal(x int, isEqual bool) string {
	suffix := "th"

	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}

	if isEqual {
		return "=" + strconv.Itoa(x) + suffix
	}

	return strconv.Itoa(x) + suffix
}
