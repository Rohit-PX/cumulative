package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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
