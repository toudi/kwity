package invoice

import (
	"time"

	"github.com/toudi/kwity/workdays"
)

type WorkdaysWithRate struct {
	Workdays int
	Rate     Rate
}

func (c *Contractor) WorkdaysWithRates(startDate time.Time, endDate time.Time) ([]WorkdaysWithRate, error) {
	var result = make([]WorkdaysWithRate, 0)
	var selectedRate Rate
	var numWorkDays int
	var err error

	// ok so basically this function is supposed to handle cases where your rate would
	// change within a month.
	// for example: you have a month that has 31 days, a rate that starts at the beginning
	// of the year at value 1, then you get a change rate at the 15'th of this month. The
	// right thing to do is to split the items into two parts:
	// 1/ from the 1st day of the month till the 14'th
	// 2/ from the 15'th day of the month till end of month, with the new rate

	// let's iterate over the days within month
	for startDate.Before(endDate) {
		err = c.applicableRate(startDate, &selectedRate)
		// we cannot detect a valid rate for this range
		if err != nil {
			return nil, err
		}
		// ok, we have found the rate. now let's check if it has started before the first day
		// that the invoice accounts for
		// example: rate starts at 1st Jan and startDate = 1st July
		if selectedRate.StartTime.Before(startDate) {
			selectedRate.StartTime = startDate
		}
		// this case is to make sure we only account for the full month and no more.
		if selectedRate.EndTime.After(endDate) {
			selectedRate.EndTime = endDate
		}
		// calculate number of workdays
		numWorkDays = workdays.CalculateWorkingDays(selectedRate.StartTime, selectedRate.EndTime) + 1
		// and move on to the end of applicable rate. it could be that this would be next month
		// altogether therefore the forloop would stop.
		startDate = selectedRate.EndTime.AddDate(0, 0, 1)
		result = append(result, WorkdaysWithRate{Workdays: numWorkDays, Rate: selectedRate})
	}

	return result, nil
}
