package fs

import (
	"path"

	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
)

type VatRatesDB struct {
	itemsView *FSDBView[common.VatRate]
}

var instance *VatRatesDB

func (f *FSDb) VATRates() db.VATRatesInterface {
	dbview, err := FSDBView_init[common.VatRate](
		path.Join(f.config.Root, "vat-rates.yaml"),
		&FSDBViewParams[common.VatRate]{
			indexer: func(document common.VatRate) []DocIndex {
				return []DocIndex{
					{Name: "id", Value: document.Id},
				}
			},
		},
	)
	if err != nil {
		panic(err)
	}
	instance = &VatRatesDB{
		itemsView: dbview,
	}
	f.OnClose(instance.itemsView.save)
	return instance
}

func (v *VatRatesDB) Sync(vatRates []common.VatRate) error {
	var err error

	for _, vatRate := range vatRates {
		if _, err = v.itemsView.AddDocument(vatRate); err != nil {
			return err
		}
	}

	return nil
}

func (v *VatRatesDB) GetByID(id string) (common.VatRate, error) {
	return v.itemsView.GetByID(id)
}
