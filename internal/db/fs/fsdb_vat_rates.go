package fs

import (
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/yti"
)

const vatRateIndexId string = "id"
const vatRatesDbFilename = "vat-rates.yaml"

type VatRatesDB struct {
	*yti.Table[common.VatRate]
}

func (f *FSDb) VATRates() db.VATRatesInterface {
	return getTable(f, vatRatesDbFilename, func(filename string) (*VatRatesDB, error) {
		instance, err := yti.OpenFile[common.VatRate](
			filename,
			&yti.TableOptions[common.VatRate]{
				Indices: map[string]yti.Indexer[common.VatRate]{
					vatRateIndexId: func(item common.VatRate) interface{} {
						return item.Id
					},
				},
			},
		)
		return &VatRatesDB{
			Table: instance,
		}, err
	})
}

func (v *VatRatesDB) Sync(vatRates []common.VatRate) error {
	var err error

	for _, vatRate := range vatRates {
		if _, err = v.CreateNXByIndexValue(vatRateIndexId, vatRate.Id, vatRate); err != nil {
			return err
		}
	}

	return nil
}

func (v *VatRatesDB) GetByID(id string) (common.VatRate, error) {
	return v.GetByIndex(vatRateIndexId, id)
}
