package common_test

import (
	"testing"

	"github.com/toudi/kwity/internal/common"
)

func TestCalculateAmount(t *testing.T) {
	type testCase struct {
		quantity  common.PriceNormalized
		unitPrice common.UnitPrice
		vatRate   common.VatRate
		expected  common.Amount
	}

	for _, test := range []testCase{
		{
			quantity:  common.PriceNormalized{Price: 1, Multiplier: 0},
			unitPrice: common.UnitPrice{Price: common.PriceNormalized{Price: 100, Multiplier: 0}},
			vatRate:   common.VatRate{Numerator: 23},
			expected:  common.Amount{Net: 10000, Gross: 12300, Vat: 2300, Multiplier: 2},
		},
		{
			quantity:  common.PriceNormalized{Price: 10},
			unitPrice: common.UnitPrice{Price: common.PriceNormalized{Price: 14927, Multiplier: 2}},
			vatRate:   common.VatRate{Numerator: 23},
			expected:  common.Amount{Net: 149270, Gross: 183602, Vat: 34332, Multiplier: 2},
		},
	} {
		amount := common.CalculateAmount(test.quantity, test.unitPrice, test.vatRate)
		amount.RoundUp(2)
		if *amount != test.expected {
			t.Errorf("unexpected amount: %v vs %v", amount, test.expected)
		}
	}
}
