package template

import (
	"errors"

	"github.com/toudi/kwity/internal/common"
)

const quantityNumWorkdays = "\"month-workdays\""

var ErrUnknownVATRate = errors.New("unknown VAT rate")

type TemplateItem struct {
	Name           string  `yaml:"name"`
	Quantity       string  `yaml:"quantity"`
	UnitPriceNet   *string `yaml:"unit-price-net,omitempty"`
	UnitPriceGross *string `yaml:"unit-price-gross,omitempty"`
	UnitPrice      common.UnitPrice
	DaysFree       int    `yaml:"days-free"`
	VatRateId      string `yaml:"vat"`
	RateId         string `yaml:"rate,omitempty"`
}

func (ti TemplateItem) CalculateNumberOfWorkdays() bool {
	return ti.Quantity == quantityNumWorkdays
}

// 	return tmpQuantityValue == quantityNumWorkdays
// }

// func (ti *TemplateItem) parseQuantity() error {
// 	var err error

// 	if tmpValueStr, ok := ti.Quantity.(string); ok {
// 		ti.QuantityNumber, err = strconv.ParseFloat(tmpValueStr, 64)
// 		return err
// 	}

// 	if tmpValueInt, ok := ti.Quantity.(int); ok {
// 		ti.QuantityNumber = float64(tmpValueInt)
// 	}

// 	return nil
// }

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

func (ti *TemplateItem) GetQuantity() (common.PriceNormalized, error) {
	return common.ParsePrice(ti.Quantity)
}
