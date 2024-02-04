package invoice

import (
	"sort"
	"strings"

	"github.com/toudi/kwity/internal/common"
)

type totalPerVATRate struct {
	VatRate common.VatRate
	common.AmountOfPrices
}

func (i *Invoice) CalculateTotalAmount() {
	i.TotalAmount = nil

	for _, item := range i.Items {
		i.UpdateTotalAmount(item)
	}
}

func (i *Invoice) UpdateTotalAmount(item *Item) {
	if i.TotalAmount == nil {
		i.TotalAmount = &common.Amount{}
		i._aggregatesPerVAT = make(map[common.VatRate]*common.Amount)
	}

	vatRate := item.VatRate
	itemAmount := item.Amount()

	i.TotalAmount.Add(itemAmount)

	// just so that it outputs on the struct nicely
	if _, exists := i._aggregatesPerVAT[vatRate]; !exists {
		i._aggregatesPerVAT[vatRate] = &common.Amount{}
	}
	i._aggregatesPerVAT[vatRate].Add(itemAmount)
}

func (i *Invoice) GetTotalsPerVATRateArray() []totalPerVATRate {
	if i._aggregatesPerVAT == nil {
		i.CalculateTotalAmount()
	}

	var totals = make([]totalPerVATRate, 0, len(i._aggregatesPerVAT))

	for vatRate, total := range i._aggregatesPerVAT {
		totals = append(totals, totalPerVATRate{
			VatRate: vatRate,
			AmountOfPrices: common.AmountOfPrices{
				Net:   common.PriceNormalized{Price: total.Net, Multiplier: total.Multiplier},
				Gross: common.PriceNormalized{Price: total.Gross, Multiplier: total.Multiplier},
				VAT:   common.PriceNormalized{Price: total.Vat, Multiplier: total.Multiplier},
			},
		})
	}

	sort.Slice(totals, func(i, j int) bool {
		left := totals[i]
		right := totals[j]

		if left.VatRate.Rate == right.VatRate.Rate {
			return strings.Compare(
				left.VatRate.GetDescription(),
				right.VatRate.GetDescription(),
			) == -1
		}

		return left.VatRate.Rate < right.VatRate.Rate
	})

	return totals
}
