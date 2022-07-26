package vesting

type Event struct {
	Type       EventType
	EmployeeID string
	Name       string
	AwardID    string
	Date       string
	Quantity   float64
}

type Events []Event

type EventType string

const (
	VEST   EventType = "VEST"
	CANCEL EventType = "CANCEL"
)
