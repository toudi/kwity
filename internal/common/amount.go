package common

import (
	"math"
)

type Amount struct {
	Net        int
	Gross      int
	Vat        int
	Multiplier int
}

type AmountOfPrices struct {
	Net   PriceNormalized
	Gross PriceNormalized
	VAT   PriceNormalized
}

func CalculateAmount(quantity PriceNormalized, unitPrice UnitPrice, vatRate VatRate) *Amount {
	multiplier := quantity.Multiplier + unitPrice.Price.Multiplier

	amount := int(quantity.Price * unitPrice.Price.Price)

	var net, gross int

	if unitPrice.IsGross {
		gross = amount
		net = vatRate.GrossToNet(gross)
	} else {
		net = amount
		gross = vatRate.NetToGross(net)
	}

	return &Amount{
		Net:        net,
		Gross:      gross,
		Vat:        gross - net,
		Multiplier: multiplier,
	}
}

func (a *Amount) RoundUp(decimalPlaces int) {
	if decimalPlaces == a.Multiplier {
		return
	}

	if a.Multiplier > decimalPlaces {
		quantizer := math.Pow10(a.Multiplier - decimalPlaces)
		a.Net = int(math.Ceil(float64(a.Net) / quantizer))
		a.Gross = int(math.Ceil(float64(a.Gross) / quantizer))
		a.Vat = int(math.Ceil(float64(a.Vat) / quantizer))
	} else if a.Multiplier < decimalPlaces {
		quantizer := int(math.Pow10(decimalPlaces - a.Multiplier))
		a.Net *= quantizer
		a.Gross *= quantizer
		a.Vat *= quantizer
	}

	a.Multiplier = decimalPlaces
}

func (a *Amount) Add(b *Amount) {
	var quantizer int

	if a.Multiplier < b.Multiplier {
		quantizer = int(math.Pow10(b.Multiplier - a.Multiplier))
		a.Net = (a.Net * quantizer) + b.Net
		a.Gross = (a.Gross * quantizer) + b.Gross
		a.Vat = (a.Vat * quantizer) + b.Vat
	} else if a.Multiplier > b.Multiplier {
		quantizer = int(math.Pow10(a.Multiplier - b.Multiplier))
		a.Net += b.Net * quantizer
		a.Gross += b.Gross * quantizer
		a.Vat += b.Vat * quantizer
	} else {
		a.Net += b.Net
		a.Gross += b.Gross
		a.Vat += b.Vat
	}
	a.Multiplier = b.Multiplier
}
