package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	championship  = "All Triumph Challenge"
	decimalPlaces = 4 // How many decimal places are used by Natsoft and to display in event results.
	natSoftURL    = "http://racing.natsoft.com.au/results/"
	help          = `Instructions to use:
	Open Natsoft racing results for the event ` + natSoftURL + `

	Select all of the individual lap times by pressing Ctrl + A
	Then copy by pressing Ctrl + C
	Run the TriumphChallenge program
	Type in all the competitors racing numbers that are in the event, each separated by a space.
	Press Enter
	Results will be generated in the same folder with today's date in spreadsheet, HTML and text format.`

	// Regular expressions
	rDriverName   = `([a-zA-Z_\-'/]+ )+`
	rRacingNumber = `\d{1,3}`
)

var (
	rNonLaps  = fmt.Sprintf(`\*:\*{2}\.\*{%d}|-:-{2}.-{%[1]d}`, decimalPlaces) // *:**.**** or -:--.----.
	rLapTimes = fmt.Sprintf(`(\d:\d{2}\.\d{%d}|%s)`, decimalPlaces, rNonLaps)  // lap time OR *:**.****.

	// matches a list of lap times by a driver.
	reHasDrivers = regexp.MustCompile(fmt.Sprintf(`\n *%s( %s)+ +((\s*\d{1,2}0 )*(%s[ p])*)*`, rRacingNumber, rDriverName, rLapTimes))
	reLapTime    = regexp.MustCompile(rLapTimes)
	reNonLaps    = regexp.MustCompile(rNonLaps)
	reRacingNum  = regexp.MustCompile(fmt.Sprintf("^%s ", rRacingNumber))
	reDriverName = regexp.MustCompile(rDriverName)

	lineDelimiter = []byte("\n")
)

// Driver represents a competitor entered in the event.
type Driver struct {
	RaceNumber string
	Name       string
	Fastest    time.Duration // The fastest time excluding Qualifying session
	Slowest    time.Duration // The slowest time excluding Qualifying session
	Qualify    time.Duration // aka Practice time
	Percentage float64       //
	Runs       uint          // aka Session. Zero based index, but the first run is ignored for Qualifying.
	Laps       uint          // Quantity of laps completed excluding Qualifying session.
	Position   uint          // Only assigned once Driver's slice has been sorted.
}

// sortResults returns a list of entered drivers, the event name, a list of drivers who are missing results and the longest driver name
// given the Natsoft results and a list of competitors entered in the event.
func sortResults(results []byte, enteredCars [][]byte) (drivers []Driver, eventName string, missing []string, longestNameLen uint) {
	eventName = eventTitle(results)

	matches := reHasDrivers.FindAll(results, -1)

	// Iterate through all competitors lap times.
	for i := range matches {
		// If this driver is a competitor.
		if driver, ok := newDriver(matches[i], enteredCars); ok {
			drivers = append(drivers, driver)

			// Work out driver names table column length used in text file output.
			if l := uint(len(driver.Name)); l > longestNameLen {
				longestNameLen = l
			}
		}
	}

	sortDrivers(drivers)

	// Find if there are any missing competitors.
	if len(drivers) != len(enteredCars) {
		for i := range enteredCars {
			if !hasRacingNum(drivers, enteredCars[i]) {
				missing = append(missing, string(enteredCars[i]))
			}
		}
	}

	return
}

func eventTitle(results []byte) string {
	lines := bytes.Split(results, lineDelimiter)
	for i := range lines {
		lines[i] = bytes.TrimSpace(lines[i])
		if len(lines[i]) >= 1 {
			return fmt.Sprintf("%s - %s", championship, string(lines[i]))
		}
	}

	return ""
}

func sortDrivers(drivers []Driver) {
	sort.SliceStable(drivers, func(i, j int) bool {
		// If either driver didn't set a qualifying lap time.
		if drivers[i].Qualify == 0 || drivers[j].Qualify == 0 {
			// Sort qualifying time in descending order (placing the slowest lap time first and zeros last).
			return drivers[i].Qualify > drivers[j].Qualify
		}

		// If both drivers didn't set a qualifying time.
		if drivers[i].Qualify == 0 && drivers[j].Qualify == 0 {
			// Sort by the quantity of laps completed in descending order (most laps first).
			return drivers[i].Laps > drivers[j].Laps
		}

		// If both drivers don't have a percentage.
		if drivers[i].Percentage == 0 && drivers[j].Percentage == 0 {
			// Sort by the qualifying time in ascending order (fastest qualifying time first).
			return drivers[i].Qualify < drivers[j].Qualify
		}

		// If both drivers have completed the same number of laps.
		if drivers[i].Runs == drivers[j].Runs {
			// Sort Percentage in descending order (highest percentage first).
			return drivers[i].Percentage > drivers[j].Percentage
		}

		// Sort by the quantity of runs/session completed in descending order (most runs/sessions first).
		return drivers[i].Runs > drivers[j].Runs
	})
}

func newDriver(line []byte, competitors [][]byte) (driver Driver, ok bool) {
	line = bytes.TrimSpace(line)
	raceNum := bytes.TrimSpace(reRacingNum.Find(line))

	// Ignore any line/entry NOT in the list of paid competitors entered for the event.
	if !has(competitors, raceNum) {
		return
	}

	driver = Driver{
		RaceNumber: string(raceNum),
		Name:       string(bytes.TrimSpace(reDriverName.Find(line))),
		Fastest:    math.MaxInt64, // Default the Fastest Lap and Qualifying Lap to the slowest possible time.
		Qualify:    math.MaxInt64,
	}

	var skipNextLap bool

	lapTimes := reLapTime.FindAll(line, -1)

	// Loop through all lap times.
	for n := range lapTimes {
		// If the lap is missing a time.
		if reNonLaps.Match(lapTimes[n]) {
			//... AND If there's another lap in the list, AND the next lap is NOT the end of the Run/Session.
			if n+1 < len(lapTimes) && !reNonLaps.Match(lapTimes[n+1]) {
				driver.Runs++
			}

			skipNextLap = true
			continue
		}

		// Skip the first lap of each run, allowing for a grid formation lap. This may change depending on which circuit the race is held at or if formation laps are organized.
		if skipNextLap {
			skipNextLap = false
			continue
		}

		lapTime, err := time.ParseDuration(
			// Convert time format 00:00.0000 to 00m00.0000s so it can be parsed.
			strings.ReplaceAll(string(lapTimes[n]), ":", "m") + "s",
		)
		if err != nil {
			log.Println(err)
			continue
		}

		if driver.Runs >= 1 {
			// Qualifying laps completed don't count towards the quantity of laps completed during the day.
			driver.Laps++

			// Only calculate the slowest lap when not in Practice/Qualifying.
			if lapTime.Seconds() > driver.Slowest.Seconds() {
				driver.Slowest = lapTime
			}

			// Calculate the fastest lap.
			if lapTime.Seconds() < driver.Fastest.Seconds() {
				driver.Fastest = lapTime
			}
		} else if lapTime.Seconds() < driver.Qualify.Seconds() {
			// Calculate the fastest qualifying lap only during the qualifying session/run.
			driver.Qualify = lapTime
		}
	}

	// If at least one session is completed,
	if driver.Runs >= 1 && driver.Qualify != math.MaxInt64 && driver.Fastest != math.MaxInt64 {
		// Calculate scoring percentage formula.
		driver.Percentage = driver.Fastest.Seconds() / ((driver.Slowest.Seconds() + driver.Qualify.Seconds()) / 2) * 100
	}

	// If Qualifying or Fastest lap times haven't been calculated, clear their values to prevent displaying erroneous results.
	if driver.Qualify == math.MaxInt64 {
		driver.Qualify = 0
	}
	if driver.Fastest == math.MaxInt64 {
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
	if err != nil || !strings.EqualFold(u.Host, "racing.natsoft.com.au") || !strings.EqualFold(u.Scheme, "http") {
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
