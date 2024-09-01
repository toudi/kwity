package workdays

import (
	"time"
)

// https://stackoverflow.com/questions/31327124/how-to-calculate-number-of-business-days-in-golang
func CalculateWorkingDays(
	startTime time.Time,
	endTime time.Time,
	bankHolidaysSet map[string]bool,
) (numWorkingDays int, numBankHolidays int) {
	for startTime.Before(endTime.AddDate(0, 0, 1)) {
		if startTime.Weekday() > time.Sunday && startTime.Weekday() < time.Saturday {
			numWorkingDays += 1
		}
		// let's check if this is a public holiday / bank holiday
		formattedDate := startTime.Format("2006-01-02")
		if bankHolidaysSet[formattedDate] {
			numBankHolidays += 1
		}
		startTime = startTime.AddDate(0, 0, 1)
	}

	return numWorkingDays, numBankHolidays
}
