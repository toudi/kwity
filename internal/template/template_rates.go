package template

import (
	"fmt"
	"time"

	"github.com/toudi/kwity/internal/common"
)

func (t *Template) parseRatesTable() error {
	if t.Rates == nil {
		return nil
	}

	var err error

	now := time.Now()
	var endOfYear = now.AddDate(0, 0, 365-now.YearDay())

	for i, rate := range t.Rates {
		rate.StartTime, err = time.Parse("2006-01-02", rate.StartDate)
		if err != nil {
			return fmt.Errorf("could not parse rate start period: %v", err)
		}
		rate.StartTime = rate.StartTime.In(now.Location())
		if i > 0 {
			t.Rates[i-1].EndTime = rate.StartTime.AddDate(0, 0, -1)
			// fmt.Printf("Set previous rate end to %v\n", ct.Rates[i-1].EndTime)
		}
		// fmt.Printf("parsed start of period as %v\n", rate.StartTime)
	}
	t.Rates[len(t.Rates)-1].EndTime = endOfYear
	// fmt.Printf("set last rate end to %v\n", endOfYear)
	// for i, rate := range ct.Rates {
	// 	// fmt.Printf("rate %d => %+v\n", i, rate)
	// }
	return nil
}

func (t *Template) applicableRateById(startDate time.Time, id string) (*common.Rate, error) {
	var definition *RateDefinition
	var i int
	// var err error

	for i, definition = range t.Rates {
		if definition.StartTime.After(startDate) {
			// if the next rate window starts *after* what we're looking for then it means that
			// the previous entry is applicable.
			i -= 1
			break
		}
	}

	if i < 0 {
		return nil, errNoRatesFound
	}

	definition = t.Rates[i]
	return &definition.Rate, nil

	// *dest = *definition

	// dest.StartTime = definition.StartTime
	// dest.EndTime = definition.EndTime
	// dest.Rate = definition.Rate
	// dest.IsGross = definition.IsGross
	// dest.Vat = definition.Vat

	// return err
}
