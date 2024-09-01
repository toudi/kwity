package common_test

import (
	"testing"

	"github.com/toudi/kwity/internal/common"
)

func TestVatRateNetToGross(t *testing.T) {
	type testCase struct {
		rate     common.VatRate
		net      int
		expected int
	}

	for _, test := range []testCase{
		{
			net:      100,
			rate:     common.VatRate{Numerator: 23},
			expected: 123,
		},
		{
			net:      14927,
			rate:     common.VatRate{Numerator: 23},
			expected: 18360,
		},
		{
			net:      11783,
			rate:     common.VatRate{Numerator: 5},
			expected: 12372,
		},
		{
			net:      100,
			rate:     common.VatRate{Numerator: 85, Denominator: 1000},
			expected: 108,
		},
		{
			net:      10700,
			rate:     common.VatRate{Numerator: 85, Denominator: 1000},
			expected: 11609,
		},
	} {
		gross := test.rate.NetToGross(test.net)
		if gross != test.expected {
			t.Errorf("unexpected gross amount: %v vs %v", gross, test.expected)
		}
	}
}

func TestVatRateGrossToNet(t *testing.T) {
	type testCase struct {
		rate     common.VatRate
		gross    int
		expected int
	}

	for _, test := range []testCase{
		{
			gross:    123,
			rate:     common.VatRate{Numerator: 23},
			expected: 100,
		},
		{
			gross:    18360,
			rate:     common.VatRate{Numerator: 23},
			expected: 14927,
		},
		{
			gross:    18360,
			rate:     common.VatRate{Numerator: 8},
			expected: 17000,
		},
	} {
		net := test.rate.GrossToNet(test.gross)
		if net != test.expected {
			t.Errorf("unexpected GrossToNet result: %d vs %d", net, test.expected)
		}
	}

}
