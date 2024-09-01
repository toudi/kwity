package template

import (
	"errors"
	"strings"

	"github.com/phuslu/log"
	"github.com/toudi/kwity/internal/common"
)

const quantityNumWorkdays = "working-days"
const quantityNumWorkingHours = "working-hours"
const suffixWithBankHolidays = "+bank-holidays"

var ErrUnknownVATRate = errors.New("unknown VAT rate")

type TemplateItem struct {
	Name           string  `yaml:"name"`
	Quantity       string  `yaml:"quantity"`
	HoursPerDay    int     `yaml:"hours-per-day"`
	UnitPriceNet   *string `yaml:"unit-price-net"`
	UnitPriceGross *string `yaml:"unit-price-gross"`
	UnitPrice      common.UnitPrice
	DaysFree       int    `yaml:"days-free"`
	VatRateId      string `yaml:"vat"`
	RateId         string `yaml:"rate"`
	Accumulate     bool   `yaml:"accumulate"`
}

func (ti *TemplateItem) GetUnitPrice() (common.UnitPrice, error) {
	var err error
	var unitPrice common.UnitPrice

	isGross := ti.UnitPriceNet == nil

	if isGross {
		if unitPrice.Price, err = common.ParsePrice(*ti.UnitPriceGross); err != nil {
			return unitPrice, err
		}
		unitPrice.IsGross = true
	} else {
		if unitPrice.Price, err = common.ParsePrice(*ti.UnitPriceNet); err != nil {
			return unitPrice, err
		}
	}

	return unitPrice, nil
}

func (ti *TemplateItem) GetQuantity(
	workingDays int,
	bankHolidays int,
) (common.PriceNormalized, error) {
	// check if this is a preautomated way of delivering quantity:
	log.Trace().
		Int("workingDays", workingDays).
		Int("bankHolidays", bankHolidays).
		Any("quantity", ti.Quantity).Msg("TemplateItem::GetQuantity")

	if strings.HasPrefix(ti.Quantity, quantityNumWorkdays) ||
		strings.HasPrefix(ti.Quantity, quantityNumWorkingHours) {
		log.Trace().Msg("preautomation detected")

		quantity := workingDays
		if !strings.HasSuffix(ti.Quantity, suffixWithBankHolidays) {
			log.Trace().Msgf("subtracting %d bank holidays", bankHolidays)
			quantity -= bankHolidays
		}
		if strings.HasPrefix(ti.Quantity, quantityNumWorkingHours) {
			log.Trace().Msgf("calculating working hours")
			hoursPerDay := 8
			if ti.HoursPerDay > 0 {
				hoursPerDay = ti.HoursPerDay
			}
			quantity *= hoursPerDay
			log.Trace().Int("quantity", quantity).Msg("")
		}
		return common.PriceNormalized{
			Price:      quantity,
			Multiplier: 0,
		}, nil
	}
	return common.ParsePrice(ti.Quantity)
}
