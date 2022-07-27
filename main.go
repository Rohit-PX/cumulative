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

	allEvents, err := vesting.GetVestingFromFile(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	filteredEvents := vesting.GetVestingBefore(targetDate, allEvents, precision)
	cumulativeVesting := vesting.GetVestingSchedule(filteredEvents, precision)
	vesting.PrintEvents(cumulativeVesting, precision, true)
}
