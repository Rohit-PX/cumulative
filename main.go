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

func GetEventsFromFile(fileName string) (*Events, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s, error: %v", fileName, err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s, error: . ")
	}

	eventsDB := Init()

	for _, event := range records {
		e := &Event{
			Type:       EventType(event[0]),
			EmployeeID: event[1],
			Name:       event[2],
			AwardID:    event[3],
			EventDate:  ConvertToDate(event[4]),
		}
		// Add to code validation part
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

// Implement delimiter as a variable it could be "," or "|"
const (
	VEST   EventType = "VEST"
	CANCEL EventType = "CANCEL"
)

func Init() *Events {
	return &Events{}
}

func PrintEvent(e *Event) {
	fmt.Printf("%s %s %s %s %s %f\n", e.Type, e.EmployeeID, e.Name, e.AwardID, (e.EventDate).Format("2006-01-02"), e.Quantity)
}

func PrintEvents(events *Events) {
	for _, e := range *events {
		PrintEvent(e)
	}
}

func GetEventsBefore(t time.Time, events *Events) {
	for _, e := range *events {
		if e.EventDate.Before(t) || e.EventDate == t {
			PrintEvent(e)
		}
	}
}

func ConvertToDate(str string) time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func main() {
	fileName := os.Args[1]
	targetDate := ConvertToDate(os.Args[2])
	allEvents, err := GetEventsFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	PrintEvents(allEvents)
	fmt.Println("------------------------------------")
	GetEventsBefore(targetDate, allEvents)
}
