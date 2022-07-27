package vesting

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/cumulative/pkg/utils"
)

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

const (
	VEST   EventType = "VEST"
	CANCEL EventType = "CANCEL"
)

func InitVestingEvents() *Events {
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

	eventsDB = InitVestingEvents()

	for _, event := range records {
		e := &Event{
			Type:       EventType(event[0]),
			EmployeeID: event[1],
			Name:       event[2],
			AwardID:    event[3],
		}
		eventDate, err := utils.ConvertToDate(event[4])
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

	// closes the file after everything is done
	defer file.Close()
	return eventsDB, nil

}

// Implement the sort interface on Events
func (e Events) Len() int { return len(e) }

func (e Events) Less(i, j int) bool {
	if e[i].EmployeeID == e[j].EmployeeID {
		return e[i].AwardID < e[j].AwardID
	}
	return e[i].EmployeeID < e[j].EmployeeID
}

func (e Events) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// GetVestingBefore - from the given list of events returns events before (and including) the given date, with given precision
func GetVestingBefore(t time.Time, events *Events, precision int) Events {
	vestingBefore := InitVestingEvents()
	for _, e := range *events {
		if !(e.EventDate.Before(t) || e.EventDate == t) {
			e.Quantity = 0
		}
		*vestingBefore = append(*vestingBefore, e)
	}
	return *vestingBefore
}

func GetVestingSchedule(events Events, precision int) *Events {
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
			// Key not found, so creating a new entry
			vestingMap[key] = &Event{
				EmployeeID: e.EmployeeID,
				Name:       e.Name,
				AwardID:    e.AwardID,
				Quantity:   e.Quantity,
			}
		}
	}

	// Create an event list from the vesting map so that it can be sorted
	flattenedVestingSchedule := InitVestingEvents()
	for _, event := range vestingMap {
		*flattenedVestingSchedule = append(*flattenedVestingSchedule, event)
	}
	sort.Sort(flattenedVestingSchedule)
	return flattenedVestingSchedule

}

func PrintEvent(e *Event, precision int, printOnlyCumulativeVesting bool) {
	if printOnlyCumulativeVesting {
		fmt.Printf("%s %s %s %s\n", e.EmployeeID, e.Name, e.AwardID, strconv.FormatFloat(e.Quantity, 'f', precision, 64))
	} else {
		fmt.Printf("%s %s %s %s %s %s\n", e.Type, e.EmployeeID, e.Name, e.AwardID, (e.EventDate).Format("2006-01-02"), strconv.FormatFloat(e.Quantity, 'f', precision, 64))
	}
}

func PrintEvents(events *Events, precision int, printOnlyCumulativeVesting bool) {
	for _, e := range *events {
		PrintEvent(e, precision, printOnlyCumulativeVesting)
	}
}
