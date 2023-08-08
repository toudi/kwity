package invoice

import "math"

type Amount struct {
	Net   float32
	Gross float32
	Vat   float32
}

func addAmount(amount Amount, another Amount) Amount {
	return Amount{
		Net:   amount.Net + another.Net,
		Gross: amount.Gross + another.Gross,
		Vat:   amount.Vat + another.Vat,
	}
}

// these two functions I got from here:
// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func roundFloat(value float32, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(round(float64(value)*output)) / float32(output)
}

func calculateAmount(quantity float32, unitPrice UnitPrice, vat *Vat) Amount {
	var result Amount
	if vat == nil {
		// there is no vat therefore let's treat it as zero
		result.Net = quantity * unitPrice.Price
		result.Gross = result.Net
	} else {
		// let's check if the unit price is already in gross form
		if unitPrice.Gross {
			// yes it is - we need to calculate net from gross.
			result.Gross = quantity * unitPrice.Price
			result.Net = roundFloat(result.Gross/(1+vat.Rate), 2)
		} else {
			// the price is given as net therefore we need to calculate gross from it.
			result.Net = quantity * unitPrice.Price
			result.Gross = roundFloat(result.Net*(1+vat.Rate), 2)
		}
	}

	result.Vat = roundFloat(result.Gross-result.Net, 2)

	return result
}
