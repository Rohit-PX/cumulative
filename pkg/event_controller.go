package vesting

import (
	"encoding/csv"
	"fmt"
	"os"
)

func Init() *Events {
	return &Events{}
}

func GetEventsFromFile(fileName string) (*Events, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s, error: %v", fileName, err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s, error: . ")
	}

	defer file.Close()
}

// Implement delimiter as a variable it could be "," or "|"
