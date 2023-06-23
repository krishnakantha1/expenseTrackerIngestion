package util

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func GetYearMonth(date string) (int, time.Month, error) {
	year, err := strconv.Atoi(date[strings.Index(date, "-")+1:])

	if err != nil {
		return 0, 0, errors.New("Could not parse Year segment")
	}

	month, err := getMonth(date)

	if err != nil {
		return 0, 0, err
	}

	return year, month, nil
}

func getMonth(date string) (time.Month, error) {
	monthStr := date[:strings.Index(date, "-")]

	switch strings.ToLower(monthStr) {
	case "jan":
		return time.January, nil
	case "feb":
		return time.February, nil
	case "mar":
		return time.March, nil
	case "apr":
		return time.April, nil
	case "may":
		return time.May, nil
	case "jun":
		return time.June, nil
	case "jul":
		return time.July, nil
	case "aug":
		return time.August, nil
	case "sep":
		return time.September, nil
	case "oct":
		return time.October, nil
	case "nov":
		return time.November, nil
	case "dec":
		return time.December, nil
	default:
		return 0, errors.New("Could not parse month segment")
	}
}
