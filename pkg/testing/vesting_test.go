package unittest

import (
	"fmt"
	"testing"

	"github.com/cumulative/pkg/utils"
	"github.com/cumulative/pkg/vesting"
)

// Tests if vesting file is read correctly
func TestGetVestingFromFile(t *testing.T) {
	validFileName := "test.csv"
	allEvents, err := vesting.GetVestingFromFile(validFileName)
	if err != nil {
		t.Fatalf("Failed to read valid file: %s", validFileName)
	}
	if len(*allEvents) != 8 {
		t.Fatalf("Expected 8 records in the file but got %d", len(*allEvents))
	}
}

// Tests if invalid vesting file is identified correctly
func TestGetVestingFromInvalidFile(t *testing.T) {
	validFileName := "test_invalid.csv"
	_, err := vesting.GetVestingFromFile(validFileName)
	if err == nil {
		t.Fatalf("Should have failed to read because file with name %s does not exist or name has been specified incorrectly ", validFileName)
	}
}

// Tests if events are correctly filtered based on vesting date
func TestGetVestingBefore(t *testing.T) {
	validFileName := "test.csv"
	targetDate, _ := utils.ConvertToDate("2020-02-01")
	precision := 2
	allEvents, err := vesting.GetVestingFromFile(validFileName)
	if err != nil {
		t.Fatalf("Failed to read valid file: %s", validFileName)
	}
	filteredEvents := vesting.GetVestingBefore(targetDate, allEvents, precision)

	var expectedQuantity float64
	expectedQuantity = 1000.00
	if filteredEvents[0].Quantity != expectedQuantity {
		t.Fatalf("This event should have quantity %f since should NOT have been filtered", expectedQuantity)
	}
	if filteredEvents[1].Quantity != float64(0) {
		t.Fatalf("This event should have quantity 0 since should have been filtered, but instead has quantity: %f", filteredEvents[1].Quantity)
	}
}

// Tests if cumulative frequency is correctly calculated based on target date and precision
func TestGetCumulativeVesting(t *testing.T) {
	validFileName := "test_fraction.csv"
	targetDate, _ := utils.ConvertToDate("2021-01-01")
	precision := 1

	allEvents, _ := vesting.GetVestingFromFile(validFileName)
	filteredEvents := vesting.GetVestingBefore(targetDate, allEvents, precision)
	cumulativeVesting := vesting.GetVestingSchedule(filteredEvents, precision)

	expectedQuantity := "299.8"
	formattedQuantity := fmt.Sprintf("%.1f", (*cumulativeVesting)[0].Quantity)
	if formattedQuantity != expectedQuantity {
		t.Fatalf("This event should have quantity %s since should NOT have been filtered instead has %s", expectedQuantity, formattedQuantity)
	}

	precision = 2
	cumulativeVesting = vesting.GetVestingSchedule(filteredEvents, precision)

	expectedQuantity = "299.75"
	formattedQuantity = fmt.Sprintf("%.2f", (*cumulativeVesting)[0].Quantity)
	if formattedQuantity != expectedQuantity {
		t.Fatalf("This event should have quantity %s since should NOT have been filtered instead has %s", expectedQuantity, formattedQuantity)
	}

}
