package common

import (
	"fmt"
	"strconv"
)

type VatRate struct {
	Id          string `mapstructure:"id"`
	Numerator   uint   `mapstructure:"numerator"`
	Denominator uint   `mapstructure:"denominator"`
	Description string `mapstructure:"description,omitempty" yaml:"description,omitempty"`
}

func (v VatRate) GetDenominator() uint {
	if v.Denominator == 0 {
		return 100
	}
	return v.Denominator
}

func (v VatRate) GetRate() float64 {
	return float64(v.Numerator) / float64(v.GetDenominator())
}

func (v VatRate) GetDescription() string {
	if v.Description != "" {
		return v.Description
	}
	return fmt.Sprintf("%s %%", strconv.FormatFloat(v.GetRate()*100.0, 'f', -1, 64))
}

func (v VatRate) ID() string {
	return v.Id
}

func (v VatRate) NetToGross(amount int) int {
	denominator := v.GetDenominator()
	gross := amount * (int(denominator + v.Numerator)) / int(denominator)
	remainder := amount * (int(denominator + v.Numerator)) % int(denominator)
	if remainder > int(denominator)/2 {
		gross += 1
	}
	return gross
}

func (v VatRate) GrossToNet(amount int) int {
	denominator := v.GetDenominator()
	net := amount * int(denominator) / (int(denominator + v.Numerator))
	remainder := amount * int(denominator) % (int(denominator + v.Numerator))
	if remainder > int(denominator)/2 {
		net += 1
	}
	return net
}
