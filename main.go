package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// move to package vesting
type Event struct {
	Type       EventType
	EmployeeID string
	Name       string
	AwardID    string
	EventDate  time.Time
	Quantity   float64
}

type Events []*Event

type EventType string

var eventsDB *Events

// Implement delimiter as a variable it could be "," or "|"
const (
	VEST   EventType = "VEST"
	CANCEL EventType = "CANCEL"
)

func Init() *Events {
	return &Events{}
}

func GetVestingFromFile(fileName string) (*Events, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s, error: %v", fileName, err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read CSV file %s, error: %v", fileName, err)
	}

	eventsDB = Init()

	for _, event := range records {
		e := &Event{
			Type:       EventType(event[0]),
			EmployeeID: event[1],
			Name:       event[2],
			AwardID:    event[3],
		}
		eventDate, err := ConvertToDate(event[4])
		if err != nil {
			return nil, err
		}
		e.EventDate = eventDate
		quantity, err := strconv.ParseFloat(event[5], 2)
		if err != nil {
			return nil, err
		}
		e.Quantity = quantity
		*eventsDB = append(*eventsDB, e)
	}
	defer file.Close()
	return eventsDB, nil

}

func GetVestingBefore(t time.Time, events *Events, precision int) Events {
	vestingBefore := Init()
	for _, e := range *events {
		if !(e.EventDate.Before(t) || e.EventDate == t) {
			e.Quantity = 0
		}
		*vestingBefore = append(*vestingBefore, e)
	}
	return *vestingBefore
}

func GetVestingSchedule(events Events, precision int) {
	vestingMap := make(map[string]*Event)
	for _, e := range events {
		key := fmt.Sprintf("%s-%s", e.EmployeeID, e.AwardID)
		if current, ok := vestingMap[key]; ok {
			if e.Type == VEST {
				current.Quantity = current.Quantity + e.Quantity
			} else if e.Type == CANCEL {
				current.Quantity = current.Quantity - e.Quantity
			}
		} else {
			// Key not found
			vestingMap[key] = &Event{
				EmployeeID: e.EmployeeID,
				Name:       e.Name,
				AwardID:    e.AwardID,
				Quantity:   e.Quantity,
			}
		}
	}

	for _, event := range vestingMap {
		PrintEvent(event, precision)
	}
}

/*
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

*/

func ValidateParams(fileParam, dateParam, precisionParam string) (string, time.Time, int, error) {
	_, err := os.Open(fileParam)
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("Failed to open file %s, error: %v. Check if file exists or file name or path to file is incorrect", fileParam, err)
	}

	targetDate, err := ConvertToDate(dateParam)
	if err != nil {
		return "", time.Time{}, 0, err
	}

	precision, err := strconv.Atoi(precisionParam)
	if err != nil || precision < 0 {
		return "", time.Time{}, 0, fmt.Errorf("Incorrect precision param  %s. Error: %v", precisionParam, err)
	}
	return fileParam, targetDate, precision, nil
}

func ConvertToDate(str string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid date format. Error: %v. Please pass date in the following format - YYYY-MM-DD")
	}
	return t, nil
}

func PrintEvent(e *Event, precision int) {
	fmt.Printf("%s %s %s %s %s %s\n", e.Type, e.EmployeeID, e.Name, e.AwardID, (e.EventDate).Format("2006-01-02"), strconv.FormatFloat(e.Quantity, 'f', precision, 64))
}

func PrintEvents(events *Events, precision int) {
	for _, e := range *events {
		PrintEvent(e, precision)
	}
}

func main() {

	// TODO add validation of input params
	fileName, targetDate, precision, err := ValidateParams(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	// TODO add error classes
	// TODO move methods to controllers/packages
	allEvents, err := GetVestingFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
	//PrintEvents(eventsDB, precision)
	filteredEvents := GetVestingBefore(targetDate, allEvents, precision)
	GetVestingSchedule(filteredEvents, precision)
}
