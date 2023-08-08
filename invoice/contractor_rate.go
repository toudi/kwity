package invoice

import (
	"errors"
	"fmt"
	"time"
)

var errNoRatesFound = errors.New("no rate could be found for selected date")
var errInvalidRates = errors.New("invalid rates")

type Rate struct {
	StartDate string `json:"start-date"`
	StartTime time.Time
	EndTime   time.Time
	Rate      float32 `json:"rate"`
	IsGross   bool    `json:"gross"`
	Vat       *Vat    `json:"vat"`
}

func (ct *Contractor) parseRates() error {
	var err error

	now := time.Now()
	var endOfYear = now.AddDate(0, 0, 365-now.YearDay())

	for i, rate := range ct.Rates {
		rate.StartTime, err = time.Parse("2006-01-02", rate.StartDate)
		if err != nil {
			return fmt.Errorf("could not parse rate start period: %v", err)
		}
		rate.StartTime = rate.StartTime.In(now.Location())
		if i > 0 {
			ct.Rates[i-1].EndTime = rate.StartTime.AddDate(0, 0, -1)
			// fmt.Printf("Set previous rate end to %v\n", ct.Rates[i-1].EndTime)
		}
		// fmt.Printf("parsed start of period as %v\n", rate.StartTime)
	}
	ct.Rates[len(ct.Rates)-1].EndTime = endOfYear
	// fmt.Printf("set last rate end to %v\n", endOfYear)
	// for i, rate := range ct.Rates {
	// 	// fmt.Printf("rate %d => %+v\n", i, rate)
	// }
	return nil
}

func (ct *Contractor) applicableRate(startDate time.Time, dest *Rate) error {
	var definition *Rate
	var i int
	var err error

	for i, definition = range ct.Rates {
		if definition.StartTime.After(startDate) {
			// if the next rate window starts *after* what we're looking for then it means that
			// the previous entry is applicable.
			i -= 1
			break
		}
	}

	if i < 0 {
		return errNoRatesFound
	}

	definition = ct.Rates[i]

	dest.StartTime = definition.StartTime
	dest.EndTime = definition.EndTime
	dest.Rate = definition.Rate
	dest.IsGross = definition.IsGross
	dest.Vat = definition.Vat

	return err
}
