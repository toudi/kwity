package common

import "math"

type UnitPrice struct {
	Price   PriceNormalized
	IsGross bool `yaml:"is-gross"`
}

func (u UnitPrice) AsFloat() float64 {
	return float64(u.Price.Price) / math.Pow10(u.Price.Multiplier)
}
