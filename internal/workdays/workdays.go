package workdays

import (
	"math"
	"time"
)

// https://stackoverflow.com/questions/31327124/how-to-calculate-number-of-business-days-in-golang
func CalculateWorkingDays(startTime time.Time, endTime time.Time) int {
	// Reduce dates to previous Mondays
	startOffset := weekday(startTime)
	startTime = startTime.AddDate(0, 0, -startOffset)
	endOffset := weekday(endTime)
	endTime = endTime.AddDate(0, 0, -endOffset)

	// Calculate weeks and days
	dif := endTime.Sub(startTime)
	weeks := int(math.Round((dif.Hours() / 24) / 7))
	days := -min(startOffset, 5) + min(endOffset, 5)

	// Calculate total days
	return weeks*5 + days
}

func weekday(d time.Time) int {
	wd := d.Weekday()
	if wd == time.Sunday {
		return 6
	}
	return int(wd) - 1
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
