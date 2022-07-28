package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cumulative/pkg/utils"
	"github.com/cumulative/pkg/vesting"
)

func main() {
	// TODO add validation of input params

	if len(os.Args) < 4 {
		log.Fatal(fmt.Errorf("Not enough arguments passed to run the program.\nExample usage: ./cumulative <path-to-file> 2020-01-01 2"))
	}
	fileName, targetDate, precision, err := utils.ValidateParams(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	// Get all vesting events from file
	allEvents, err := vesting.GetVestingFromFile(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// Filter events based on the target date
	filteredEvents := vesting.GetVestingBefore(targetDate, allEvents, precision)

	// Calculate cumulative vesting for the filtered events
	cumulativeVesting := vesting.GetVestingSchedule(filteredEvents, precision)

	// Print cumulative vesting
	vesting.PrintEvents(cumulativeVesting, precision, true)
}
