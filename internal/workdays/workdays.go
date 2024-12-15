package workdays

import (
	"errors"
	"strings"
	"time"
)

type ptoRange struct {
	start time.Time
	end   time.Time
}

func (r ptoRange) Includes(date time.Time) bool {
	return r.start.Equal(date) || r.end.Equal(date) || (r.start.Before(date) && date.Before(r.end))
}

type WorkingDaysCalculator struct {
	bankHolidaysSet   map[string]bool
	paidTimeOffRanges []ptoRange
}

type WorkingDaysResult struct {
	WorkingDays    int
	BankHolidays   int
	PtoDays        int
	LastWorkingDay time.Time
}

func NewWorkingDaysCalculator(bankHolidaysSet map[string]bool, ptoRanges []string) (calculator *WorkingDaysCalculator, err error) {
	calculator = &WorkingDaysCalculator{
		bankHolidaysSet: bankHolidaysSet,
	}

	var rangeStart, rangeEnd time.Time

	// let's iterate over pto ranges (if any) and try to convert them into time objects.
	for _, ptoRangeString := range ptoRanges {
		// a pto range (defined as string) could either be:
		// a single date (e.g. "2018-04-05") or
		// a date range (e.g. "2018-04-05 - 2018-04-11") or
		// a range without an end, which will indicate range until the end of specified month
		//  e.g. "2018-04-05 -" will indicate a break from 2018-04-05 -- 2024-04-31 (or whatever the last day happens to be)
		// we can detect which type of range it is by counting the number of dashes:
		numDashes := strings.Count(ptoRangeString, "-")
		if numDashes <= 5 {
			rangeStart, err = time.ParseInLocation("2006-01-02", ptoRangeString, time.Local)
			if err != nil {
				return nil, err
			}
			if numDashes == 5 {
				rangeEnd, err = time.ParseInLocation("2006-01-02", ptoRangeString[len(ptoRangeString)-10:], time.Local)
				if err != nil {
					return nil, err
				}
			} else if numDashes == 3 {
				rangeEnd = rangeStart.AddDate(0, 1, -rangeStart.Day())
			} else if numDashes == 2 {
				rangeEnd = rangeStart
			} else {
				return nil, errors.New("Unrecognized PTO range format")
			}
			calculator.paidTimeOffRanges = append(calculator.paidTimeOffRanges, ptoRange{start: rangeStart, end: rangeEnd})
			continue
		}
		// more than 5 dashes means we cannot recognize what it is.
		return nil, errors.New("unrecognized PTO range format")
	}
	return calculator, err
}

func (w *WorkingDaysCalculator) CalculateWorkingDays(startTime time.Time, endTime time.Time, lastDayCounts bool) WorkingDaysResult {
	var result WorkingDaysResult

	var isWeekend bool

	if lastDayCounts {
		endTime = endTime.AddDate(0, 0, 1)
	}

	for startTime.Before(endTime) {
		// let's check if this is a working day:
		isWeekend = startTime.Weekday() == time.Sunday || startTime.Weekday() == time.Saturday

		if !isWeekend {
			// it could be a working day but let's check if it's not a banking holiday first
			formattedDate := startTime.Format("2006-01-02")
			if w.bankHolidaysSet[formattedDate] {
				result.BankHolidays += 1
			} else {
				result.WorkingDays += 1
				result.LastWorkingDay = startTime

				// now we can also check if it's a PTO day:
				for _, ptoRange := range w.paidTimeOffRanges {
					if ptoRange.Includes(startTime) {
						result.PtoDays += 1
						break
					}
				}
			}
		}

		startTime = startTime.AddDate(0, 0, 1)
	}

	return result
}

func (w *WorkingDaysCalculator) CalculateFromStartOfYear(endTime time.Time, lastDayCounts bool) WorkingDaysResult {
	now := time.Now()
	// calculate first day of year
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	// make sure that we truncate the date to the start of day. Otherwise, we'd be comparing against
	// the next day due to call to .Before()
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local)
	return w.CalculateWorkingDays(startOfYear, endTime, lastDayCounts)
}

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
