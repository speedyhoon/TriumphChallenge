package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	championship  = "All Triumph Challenge"
	decimalPlaces = 4 // How many decimal places are used by Natsoft and to display in event results

	// Regular expressions
	rDriverName   = `[a-z A-Z_\-']+`
	rRacingNumber = `\d{1,3}`
	natSoftURL    = "http://racing.natsoft.com.au/results/"
	help          = `Instructions to use:
	Open Natsoft racing results for the event ` + natSoftURL + `

	Select all of the individual lap times by pressing Ctrl + A
	Then copy by pressing Ctrl + C
	Run the TriumphChallenge program
	Type in all the competitors racing numbers that are in the event, each separated by a space.
	Press Enter
	Results will be generated in the same folder with today's date in spreadsheet, HTML and text format.`
)

var (
	lineDelimiter = []byte("\n")

	endOfSession = fmt.Sprintf(`\*:\*{2}\.\*{%d}|-:-{2}.-{%[1]d}`, decimalPlaces)    // *:**.**** or -:--.----
	rTimes       = fmt.Sprintf(`(\d:\d{2}\.\d{%d}|%s)`, decimalPlaces, endOfSession) // lap time OR *:**.****

	// matches a list of lap times by a driver
	hasDrivers = regexp.MustCompile(fmt.Sprintf(`\n *%s( %s)+ +((\s*\d{1,2}0 )*(%s[ p])*)*`, rRacingNumber, rDriverName, rTimes))
	lapTime    = regexp.MustCompile(rTimes)
	missingLap = regexp.MustCompile(endOfSession)
	racingNum  = regexp.MustCompile(fmt.Sprintf("^%s ", rRacingNumber))
	driverName = regexp.MustCompile(rDriverName)
)

// Driver represents a competitor entered in the event
type Driver struct {
	RaceNumber string
	Name       string
	Fastest    time.Duration // The fastest time excluding Qualifying session
	Slowest    time.Duration // The slowest time excluding Qualifying session
	Qualify    time.Duration // aka Practice time
	Percentage float64       //
	Runs       uint          // aka Session. Zero based index, but the first run is ignored for Qualifying.
	Laps       uint          // Quantity of laps completed excluding Qualifying session
	Position   uint          // Only assigned once Driver's slice has been sorted
}

// sortResults returns a list of entered drivers, the event name, a list of drivers who are missing results and the longest driver name
// given the Natsoft results and a list of competitors entered in the event.
func sortResults(results []byte, enteredCars [][]byte) (drivers []Driver, eventName string, missing []string, longestNameLen uint) {
	lines := bytes.Split(results, lineDelimiter)
	for i := range lines {
		lines[i] = bytes.TrimSpace(lines[i])
		if len(lines[i]) >= 1 {
			eventName = fmt.Sprintf("%s - %s", championship, string(lines[i]))
			break
		}
	}

	matches := hasDrivers.FindAll(results, -1)

	// Iterate through all competitors lap times
	for i := range matches {
		// If this driver is a competitor
		if driver, ok := newDriver(matches[i], enteredCars); ok {
			drivers = append(drivers, driver)

			// Work out driver names table column length used in text file output
			if l := uint(len(driver.Name)); l > longestNameLen {
				longestNameLen = l
			}
		}
	}

	sort.SliceStable(drivers, func(i, j int) bool {
		if drivers[i].Runs == drivers[j].Runs {
			return drivers[i].Percentage > drivers[j].Percentage
		}
		return drivers[i].Runs > drivers[j].Runs
	})

	// Find if there are any missing competitors
	if len(drivers) != len(enteredCars) {
		for i := range enteredCars {
			if !hasRacingNum(drivers, enteredCars[i]) {
				missing = append(missing, string(enteredCars[i]))
			}
		}
	}

	return
}

func newDriver(line []byte, competitors [][]byte) (driver Driver, ok bool) {
	line = bytes.TrimSpace(line)
	raceNum := bytes.TrimSpace(racingNum.Find(line))

	// Ignore any line/entry NOT in the list of paid competitors entered for the event
	if !has(competitors, raceNum) {
		return
	}

	line = bytes.TrimPrefix(line, raceNum)
	name := bytes.TrimSpace(driverName.Find(line))
	line = bytes.TrimSpace(bytes.TrimPrefix(line, name))

	driver.RaceNumber = string(raceNum)
	driver.Name = string(name)
	driver.Fastest = math.MaxInt64 // Default the Fastest Lap to the slowest possible time.

	var skipNextLap bool

	lapTimes := lapTime.FindAll(line, -1)

	// Loop through all lap times.
	for n := range lapTimes {
		// If the lap is missing a time
		if missingLap.Match(lapTimes[n]) {
			//... AND If there's another lap in the list, AND the next lap is NOT the end of the Run/Session.
			if n+1 < len(lapTimes) && !missingLap.Match(lapTimes[n+1]) {
				driver.Runs++
			}

			skipNextLap = true

			if driver.Runs == 1 {
				// Set the fastest lap time obtained during qualifying/practice session
				driver.Qualify = driver.Fastest
			}
			continue
		}

		// Skip the first lap of each run, allowing for a grid formation lap. This may change depending on which circuit the race is held at or if formation laps are organized.
		if skipNextLap {
			skipNextLap = false
			continue
		}

		// Convert time format 00:00.0000 to 00m00.0000s so it can be parsed.
		lapTime := strings.ReplaceAll(string(lapTimes[n]), ":", "m") + "s"
		d, err := time.ParseDuration(lapTime)
		if err != nil {
			log.Println(err)
			continue
		}

		if driver.Runs >= 1 {
			// Qualifying laps completed don't count towards the quantity of laps completed during the day.
			driver.Laps++

			// Only calculate the slowest lap when not in Practice/Qualifying
			if d.Seconds() > driver.Slowest.Seconds() {
				driver.Slowest = d
			}
		}

		// Calculate the fastest lap even when in Run=0 (practice) so the qualifying time is set.
		if d.Seconds() < driver.Fastest.Seconds() {
			driver.Fastest = d
		}
	}

	// If at least one session is completed,
	if driver.Runs >= 1 {
		// Calculate scoring percentage formula.
		driver.Percentage = driver.Fastest.Seconds() / ((driver.Slowest.Seconds() + driver.Qualify.Seconds()) / 2) * 100
	} else {
		// Otherwise finish setting the qualifying time
		driver.Qualify = driver.Fastest
		driver.Fastest = 0
	}

	return driver, true
}

func retrieveBody(path string) (src []byte) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}

	u, err := url.Parse(path)
	if err != nil || u.Host != "racing.natsoft.com.au" || u.Scheme != "http" && u.Scheme != "https" {
		return
	}

	fmt.Println("Attempting to retrieve results from Natsoft.")

	var resp *http.Response
	resp, err = http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	src, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return
}
