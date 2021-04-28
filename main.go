package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

const (
	competitorsFile = "competitors.txt"
	filePermission  = 0600
)

func main() {
	flag.Usage = func() {
		fmt.Println(help)
		flag.PrintDefaults()
	}
	flag.Parse()

	fmt.Println(championship)

	src := getEventResults()

	comps := getCompetitorsFile()
	if len(comps) == 0 {
		// Keep checking standard input for a list of competitors numbers to be entered
		for len(comps) == 0 {
			comps = prepareComps(input())
		}

		checkErr(ioutil.WriteFile(competitorsFile, bytes.Join(comps, []byte(" ")), filePermission))
	}

	render(sortResults(src, comps))
}

func getEventResults() (src []byte) {
	// Try to find a text file with today's date.
	filename := time.Now().Format("event-2006-01-02.txt")

	printOnce := true

	for {
		// Check clipboard for event results
		s, _ := clipboard.ReadAll()
		if reHasDrivers.MatchString(s) {
			src = []byte(s)
			break
		}

		// Check if clipboard contained a URL
		src = retrieveBody(s)
		if reHasDrivers.Match(src) {
			break
		}

		// Check file
		src, _ = ioutil.ReadFile(filename)

		if reHasDrivers.Match(src) {
			fmt.Println("Using the results from", filename)
			return
		}

		if printOnce {
			fmt.Printf("\n\nNo results found in %s or the clipboard. Please copy event results from %s\nDo you want to open Natsoft in your default browser? [ y / n ]\n", filename, natSoftURL)
			if yes(input()) {
				openBrowser()
			}
			printOnce = false
		}

		// Check standard input
		src = input()
		if reHasDrivers.Match(src) {
			break
		}

		// Check if standard input contained a URL
		if body := retrieveBody(string(src)); body != nil {
			src = body
		}
		if reHasDrivers.Match(src) {
			break
		}

		// Provide some commands to exit if user gets ultra stuck
		exit(src)
	}

	checkErr(ioutil.WriteFile(filename, src, filePermission))

	return
}

func yes(input []byte) bool {
	input = bytes.TrimSpace(input)

	return bytes.EqualFold(input, []byte("y")) || bytes.EqualFold(input, []byte("yes"))
}

func input() []byte {
	// ReadString will block until the delimiter is entered.
	input, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return input
	}

	// Remove the '\n' delimiter from the string.
	return bytes.TrimSpace(input)
}

func getCompetitorsFile() [][]byte {
	src, err := ioutil.ReadFile(competitorsFile)
	if err != nil || len(src) == 0 {
		fmt.Println("Please enter racing numbers separated by a space.")
		return nil
	}

	fmt.Println("Using the list of competitors in", competitorsFile)

	return prepareComps(src)
}

func prepareComps(src []byte) (competitors [][]byte) {
	src = bytes.TrimSpace(src)
	lines := bytes.Split(src, lineDelimiter)
	for i := range lines {
		lines[i] = bytes.TrimSpace(lines[i])
		// Ignore any commented out lines prefixed with #
		if !bytes.HasPrefix(lines[i], []byte("#")) {
			words := bytes.Split(lines[i], []byte(" "))
			for j := range words {
				words[j] = bytes.TrimSpace(words[j])
				// Check if the competitor is already in the list
				if len(words[j]) >= 1 && !has(competitors, words[j]) {
					competitors = append(competitors, words[j])
				}
			}
		}
	}

	return
}

func exit(src []byte) {
	switch strings.ToLower(string(src)) {
	case "x", "exit", "q", "quit", "s", "stop", "h", "halt":
		os.Exit(5)
	}
}
