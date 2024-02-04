package template

import (
	"errors"
)

var errNoRatesFound = errors.New("no rate could be found for selected date")
var errInvalidRates = errors.New("invalid rates")
