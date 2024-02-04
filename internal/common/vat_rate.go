package common

import (
	"fmt"
	"strconv"
)

type VatRate struct {
	Id          string  `mapstructure:"id"`
	Rate        float64 `mapstructure:"rate"`
	Description string  `mapstructure:"description,omitempty" yaml:"description,omitempty"`
}

func (v VatRate) GetDescription() string {
	if v.Description != "" {
		return v.Description
	}
	return fmt.Sprintf("%s %%", strconv.FormatFloat(float64(v.Rate*100), 'f', -1, 64))
}

func (v VatRate) ID() string {
	return v.Id
}

func (v VatRate) NetToGross(amount int) int {
	return int(float64(amount) * (1.0 + v.Rate))
}

func (v VatRate) GrossToNet(amount int) int {
	return int(float64(amount) / (1 + v.Rate))
}
