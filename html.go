package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func htmlHeading(eventName string, driversQty uint) *bytes.Buffer {
	return bytes.NewBufferString(
		fmt.Sprintf(`<!DOCTYPE html><html lang=en><title>%s</title><link rel=icon href="data:image/x-icon;base64,%s"><style>body{font-family:sans-serif}h1{color:#07f;text-align:center}table{width:100%%}th{text-align:left}</style><h1>%[1]s</h1><b>%[3]s %d</b><table><thead><tr><th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%s<th>%[13]s<th>%[14]s<tbody>`,
			eventName,
			favicon,
			hCompetitors,
			driversQty,
			hPosition,
			hRacingNumber,
			hDriver,
			hQualify,
			hSeconds,
			hFastest,
			hSlowest,
			hAverage,
			hPercentage,
			hRuns,
			hLaps,
		),
	)
}

func htmlRow(html io.Writer, d *Driver, ordinal string) {
	_, err := fmt.Fprintf(html, "<tr><td>%s<td>%s<td>%s<td>%v<td>%.4f<td>%v<td>%.4f<td>%v<td>%.4f<td>%.5f<td>%.8f<td>%d<td>%d",
		ordinal,
		d.RaceNumber,
		d.Name,
		d.Qualify, d.Qualify.Seconds(),
		d.Fastest, d.Fastest.Seconds(),
		d.Slowest, d.Slowest.Seconds(),
		(d.Slowest.Seconds()+d.Qualify.Seconds())/2,
		d.Percentage,
		d.Runs,
		d.Laps,
	)
	checkErr(err)
}

func htmlFooter(html *bytes.Buffer, missingCars []string) {
	_, err := fmt.Fprint(html, "</table>")
	checkErr(err)

	if len(missingCars) >= 1 {
		_, err = fmt.Fprintf(html, "<h3>%s</h3><ul><li>%s</ul>", hMissing, strings.Join(missingCars, "<li>"))
		checkErr(err)
	}
}
