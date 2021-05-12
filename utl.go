package main

import (
	"bytes"
	"fmt"
)

// hasRacingNum returns true if raceNumber is one of the drivers racing number.
func hasRacingNum(drivers []Driver, raceNumber []byte) bool {
	for i := range drivers {
		if bytes.EqualFold([]byte(drivers[i].RaceNumber), raceNumber) {
			return true
		}
	}

	return false
}

// has checks if item is within the list.
func has(list [][]byte, item []byte) bool {
	for i := range list {
		if bytes.EqualFold(list[i], item) {
			return true
		}
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func yes(input []byte) bool {
	input = bytes.TrimSpace(input)

	return bytes.EqualFold(input, []byte("y")) || bytes.EqualFold(input, []byte("yes"))
}
