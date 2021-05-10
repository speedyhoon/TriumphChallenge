package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func txtHeading(eventName string, driversQty, longestNameLen uint) *bytes.Buffer {
	//	-	Pad with spaces on the right rather than the left (left-justify the field).
	//	*	Width or precision value taken from the integer preceding the one to format.
	return bytes.NewBufferString(
		fmt.Sprintf("   %s%sCompetitors: %d%s%-5s  %4s %-*s  %-10s    %-8s    %-10s    %-8s    %-10s    %-8s    %-9s    %-11s    %4s    %4s%[2]s",
			eventName,
			newLine,
			driversQty,
			newLine,
			hPosition,
			hRacingNumber,
			longestNameLen, hDriver,
			hQualify, hSeconds,
			hFastest, hSeconds,
			hSlowest, hSeconds,
			hAverage,
			hPercentage,
			hRuns,
			hLaps,
		))
}

func textRow(txt io.Writer, d *Driver, ordinal string, longestNameLen uint) {
	/*	-	Pad with spaces on the right rather than the left (left-justify the field).
		*	Width or precision value taken from the integer preceding the one to format.
		%9f    width 9, default precision
		%9.4f  width 9, precision 4 */
	_, err := fmt.Fprintf(txt, "%-5s  %4s %-*s  %-10v    %-8.4f    %-10v    %-8.4f    %-10v    %-8.4f    %9.5f    %11.8f    %4d    %4d%s",
		ordinal,
		d.RaceNumber,
		longestNameLen, d.Name,
		d.Qualify, d.Qualify.Seconds(),
		d.Fastest, d.Fastest.Seconds(),
		d.Slowest, d.Slowest.Seconds(),
		d.SlowAv,
		d.Percentage,
		d.Runs,
		d.Laps,
		newLine,
	)
	checkErr(err)
}

func textFooter(txt io.Writer, missingCars []string) {
	if len(missingCars) >= 1 {
		_, err := fmt.Fprintf(txt, "%s%s%[1]s%[3]s", newLine, hMissing, strings.Join(missingCars, newLine))
		checkErr(err)
	}
}
