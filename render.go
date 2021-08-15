package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/speedyhoon/utl"
)

const (
	// Column headings.
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

	var position, i uint
	var isEqual bool
	for i = range drivers {
		// Calculate if the next or previous competitor had an identical score.
		if i >= 1 && drivers[i].Percentage == drivers[i-1].Percentage && drivers[i].Runs == drivers[i-1].Runs && drivers[i].Laps == drivers[i-1].Laps {
			isEqual = true
		} else {
			// Ignore positions occupied by drivers with equal results, like: =1st, =1st, =3rd, =3rd, =5th, =5th. Change to `position++` for the opposite effect like: =1st, =1st, =2nd, =2nd, =3rd, =3rd.
			position = i + 1

			isEqual = position < uint(len(drivers)) && drivers[i].Percentage == drivers[i+1].Percentage && drivers[i].Runs == drivers[i+1].Runs && drivers[i].Laps == drivers[i+1].Laps
		}

		ord := utl.Ordinal(position, isEqual)

		htmlRow(html, &drivers[i], ord)
		textRow(txt, &drivers[i], ord, longestNameLen)
		excelRow(excel, &drivers[i], ord, &spreadsheetRow)
	}

	htmlFooter(html, missingCars)
	textFooter(txt, missingCars)
	excelFooter(excel, &spreadsheetRow, missingCars)

	// Print text output to screen.
	fmt.Println(txt.String())

	fileName := time.Now().Format("results-2006-01-02 3;4;05")

	checkErr(ioutil.WriteFile(fileName+".txt", txt.Bytes(), filePermission))
	checkErr(ioutil.WriteFile(fileName+".html", html.Bytes(), filePermission))
	checkErr(excel.SaveAs(fileName + ".xlsx"))
}
