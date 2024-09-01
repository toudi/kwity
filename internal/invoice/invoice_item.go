package invoice

import (
	"github.com/toudi/kwity/internal/common"
)

type Item struct {
	Name      string                 `yaml:"name"`
	VatRate   common.VatRate         `yaml:"-"`
	VatRateId string                 `yaml:"vat-rate"`
	Quantity  common.PriceNormalized `yaml:"quantity"`
	UnitPrice common.UnitPrice       `yaml:"unit-price"`
	Metadata  map[string]interface{} `yaml:"metadata,omitempty"`
}

func (ii *Item) Amount() *common.Amount {
	amount := common.CalculateAmount(ii.Quantity, ii.UnitPrice, ii.VatRate)
	amount.RoundUp(2)
	return amount
}
