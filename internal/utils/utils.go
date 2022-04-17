package utils

import (
	"errors"
	"time"
)

// Date in format YYYY-MM-DD = 2006-01-02
func DatesBetween(startDate, endDate, layout string, allowFridays, allowWeekends bool) ([]string, error) {

	datesBetween := []string{}

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, err
	}

	for start.Before(end) {

		if !allowWeekends && IsWeekend(start) {
			start = start.Add(time.Hour * 24)
			continue
		}

		if !allowFridays && IsFriday(start) {
			start = start.Add(time.Hour * 24)
			continue
		}

		datesBetween = append(datesBetween, start.Format(layout))

		start = start.Add(time.Hour * 24)
	}

	if start.Equal(end) && (allowWeekends || !IsWeekend(start)) && (allowFridays || !IsFriday(start)) {
		datesBetween = append(datesBetween, start.Format(layout))
	}

	return datesBetween, nil
}

func DatesFromNow(dateLayout string, numDays int) ([]string, error) {

	dates := []string{}

	if numDays < 1 {
		return dates, errors.New("number of days must be greater than 0")
	}

	timestamp := time.Now()
	dates = append(dates, timestamp.Format(dateLayout))

	for i := 1; i < numDays; i++ {
		timestamp = timestamp.Add(time.Hour * 24)
		dates = append(dates, timestamp.Format(dateLayout))
	}

	return dates, nil
}

func IsWeekend(t time.Time) bool {
	t = t.UTC()
	switch t.Weekday() {
	case time.Saturday:
		return true
	case time.Sunday:
		return true
	}

	return false
}

func IsFriday(t time.Time) bool {
	t = t.UTC()
	switch t.Weekday() {
	case time.Friday:
		return true
	}

	return false
}

func ValidDates(startDate, endDate, dateLayout string) (bool, error) {
	timeStartDate, err := time.Parse(dateLayout, startDate)
	if err != nil {
		return false, err
	}
	timeEndDate, err := time.Parse(dateLayout, endDate)
	if err != nil {
		return false, err
	}

	if timeStartDate.After(timeEndDate) {
		return false, nil
	}

	return true, nil
}
