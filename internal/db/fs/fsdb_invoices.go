package fs

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
)

type InvoicesDB struct {
	root string
	db   *FSDb
}

var (
	ErrOpeningInvoices          = errors.New("unable to open invoices db")
	ErrInstantiatingIDGenerator = errors.New("unable to instantiate nanoid generator")
)

func (f *FSDb) Invoices() db.InvoicesInterface {
	return &InvoicesDB{root: f.config.Root, db: f}
}

func invoiceIndexer(i *invoice.Invoice) []DocIndex {
	return []DocIndex{
		{Name: "id", Value: i.Id},
	}
}

func invoiceSaveHook(i *invoice.Invoice) {
	for _, item := range i.Items {
		item.VatRateId = item.VatRate.Id
	}
}

func invoiceLoadHookFactory(idb *InvoicesDB) LoadHookFunc[*invoice.Invoice] {
	return func(i *invoice.Invoice) {
		for _, item := range i.Items {
			item.VatRate, _ = idb.db.VATRates().GetByID(item.VatRateId)
		}
	}
}

func (idb *InvoicesDB) Save(i *invoice.Invoice) error {
	// determine path based on the sales date
	dbPath := path.Join(idb.root, fmt.Sprintf("%d/%d.yaml", i.SaleDate.Year(), i.SaleDate.Year()))
	view, err := FSDBView_init[*invoice.Invoice](dbPath, &FSDBViewParams[*invoice.Invoice]{
		indexer:  invoiceIndexer,
		saveHook: invoiceSaveHook,
		loadHook: invoiceLoadHookFactory(idb),
	})
	if err != nil {
		return errors.Join(ErrOpeningInvoices, err)
	}
	if i.Id == "" {
		generator, err := nanoid.Standard(4)
		if err != nil {
			return errors.Join(ErrInstantiatingIDGenerator, err)
		}
		var contains = true
		for contains {
			i.Id = generator()
			contains, err = view.IndexContainsValue("id", i.Id)
			if err != nil {
				return err
			}
		}
	}
	if err = view.UpsertDocument(DocIndex{Name: "id", Value: i.Id}, i); err != nil {
		return err
	}
	return view.save()
}

func (idb *InvoicesDB) FilterByRecipient(
	recipientID string,
	issued time.Time,
	id string,
) ([]*invoice.Invoice, error) {
	dbPath := path.Join(idb.root, fmt.Sprintf("%d/%d.yaml", issued.Year(), issued.Year()))
	view, err := FSDBView_init[*invoice.Invoice](dbPath, &FSDBViewParams[*invoice.Invoice]{
		indexer:  invoiceIndexer,
		loadHook: invoiceLoadHookFactory(idb),
	})
	if err != nil {
		return nil, ErrOpeningInvoices
	}
	var dest []*invoice.Invoice = make([]*invoice.Invoice, 0)

	view.ForEach(func(document *invoice.Invoice) {
		var match = document.RecipientId == recipientID

		if id != "" && document.Id != id {
			match = false
		}

		if match {
			view.params.loadHook(document)
			dest = append(dest, document)
		}

	})

	return dest, nil
}

func (idb *InvoicesDB) GetNextSequenceNumber(
	saleDate time.Time,
	draft bool,
) (common.SequenceNumber, error) {
	var sn = common.SequenceNumber{}

	dbPath := path.Join(idb.root, fmt.Sprintf("%d/%d.yaml", saleDate.Year(), saleDate.Year()))
	view, err := FSDBView_init[*invoice.Invoice](dbPath, nil)
	if err != nil {
		return sn, err
	}

	view.ForEach(func(invoice *invoice.Invoice) {
		if invoice.Number != "" && invoice.Draft == draft {
			sn.Year += 1
			if invoice.SaleDate.Month() == saleDate.Month() {
				sn.Month += 1
			}
		}
	})

	sn.Year += 1
	sn.Month += 1

	return sn, nil
}
