package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toudi/kwity/internal/common"
)

func TestPriceNormalizedFormat(t *testing.T) {
	type unitTest struct {
		thousandSeparator string
		price             common.PriceNormalized
		expected          string
	}

	for _, test := range []unitTest{
		{
			expected: "12.34",
			price:    common.PriceNormalized{Price: 1234, Multiplier: 2},
		},
		{
			expected: "16.80",
			price:    common.PriceNormalized{Price: 168, Multiplier: 1},
		},
		{
			expected: "16.08",
			price:    common.PriceNormalized{Price: 1608, Multiplier: 2},
		},
		{
			expected: "123.405",
			price:    common.PriceNormalized{Price: 123405, Multiplier: 3},
		},
	} {
		test := test
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.price.Format(test.thousandSeparator))
		})
	}
}
