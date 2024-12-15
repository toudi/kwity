package fs

import (
	"errors"
	"fmt"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
	"github.com/toudi/yti"
)

type InvoicesDB struct {
	root         string
	db           *FSDb
	uidGenerator func() string
}

var (
	ErrOpeningInvoices          = errors.New("unable to open invoices db")
	ErrInstantiatingIDGenerator = errors.New("unable to instantiate nanoid generator")
	ErrDraftDoesNotExist        = errors.New("draft with the specified ID does not exist")
	ErrDraftAlreadyCommitted    = errors.New("draft with the specified ID was already committed")
	ErrUnknownEntity            = errors.New("unknown entity")
)

const InvoiceIndexId = "id"

func (f *FSDb) Invoices() db.InvoicesInterface {
	generator, err := nanoid.Standard(4)
	if err != nil {
		panic("unable to instantiate generator")
	}
	return &InvoicesDB{root: f.config.Root, db: f, uidGenerator: generator}
}

func (idb *InvoicesDB) loadVatRates(invoice *invoice.Invoice) {
	for _, item := range invoice.Items {
		item.VatRate, _ = idb.db.VATRates().GetByID(item.VatRateId)
	}
}

type InvoicesTable struct {
	*yti.Table[*invoice.Invoice]
}

func (idb *InvoicesDB) getTable(date time.Time) *InvoicesTable {
	dbPath := fmt.Sprintf("invoices/%d.yaml", date.Year())

	return getTable(idb.db, dbPath, func(filename string) (*InvoicesTable, error) {
		instance, err := yti.OpenFile[*invoice.Invoice](
			filename,
			&yti.TableOptions[*invoice.Invoice]{
				Indices: map[string]yti.Indexer[*invoice.Invoice]{
					InvoiceIndexId: func(item *invoice.Invoice) interface{} {
						return item.Id
					},
				},
			},
		)

		return &InvoicesTable{
			Table: instance,
		}, err
	})
}

func (idb *InvoicesDB) Save(i *invoice.Invoice) error {
	// determine path based on the sales date
	table := idb.getTable(i.SaleDate)

	if i.Id == "" {
		id, err := table.EnsureIndexDoesNotContain(
			InvoiceIndexId,
			func() interface{} { return idb.uidGenerator() },
			100,
		)
		if err != nil {
			return err
		}
		i.Id = id.(string)
	}

	for _, item := range i.Items {
		item.VatRateId = item.VatRate.Id
	}

	return table.UpdateOrCreateByIndexValue(InvoiceIndexId, i.Id, i)
}

func (idb *InvoicesDB) FilterByRecipient(
	recipientID string,
	issued time.Time,
	id string,
) ([]*invoice.Invoice, error) {
	table := idb.getTable(issued)

	var dest []*invoice.Invoice

	table.ForEach(func(item *invoice.Invoice) bool {
		var match = item.RecipientId == recipientID

		if id != "" && item.Id != id {
			match = false
		}

		if match {
			idb.loadVatRates(item)
			dest = append(dest, item)
		}
		return false
	})

	return dest, nil
}

func (idb *InvoicesDB) GetNextSequenceNumber(
	saleDate time.Time,
	draft bool,
) (common.SequenceNumber, error) {
	var sn = common.SequenceNumber{}

	table := idb.getTable(saleDate)
	table.ForEach(func(invoice *invoice.Invoice) bool {
		if invoice.Number != "" && invoice.Draft == draft {
			sn.Year += 1
			if invoice.SaleDate.Month() == saleDate.Month() {
				sn.Month += 1
			}
		}

		return false
	})

	sn.Year += 1
	sn.Month += 1

	return sn, nil
}

func (idb *InvoicesDB) GenerateInvoiceFromDraft(
	id string,
	issueDate time.Time,
) (*invoice.Invoice, error) {
	// first, narrow down the search to current year
	if issueDate.IsZero() {
		issueDate = time.Now().Local()
	}

	table := idb.getTable(issueDate)

	var draft *invoice.Invoice
	var committed *invoice.Invoice

	table.ForEach(func(i *invoice.Invoice) bool {
		if i.Id == id {
			draft = i
		}
		if i.DraftId == id {
			committed = i
		}
		return draft != nil && committed != nil
	})

	// does the draft even exist ?
	if draft == nil || (draft != nil && !draft.Draft) {
		return nil, ErrDraftDoesNotExist
	}

	// seems that it does. Was it already committed?
	if committed != nil {
		return nil, ErrDraftAlreadyCommitted
	}

	idb.loadVatRates(draft)

	// seems this is our lucky day.
	committed = &invoice.Invoice{
		DraftId:     id,
		IssueDate:   issueDate,
		SaleDate:    draft.SaleDate,
		Draft:       false,
		RecipientId: draft.RecipientId,
		BuyerId:     draft.BuyerId,
	}

	if committed.RecipientId != "" {
		if recipient, err := idb.db.Entities().GetByNIP(committed.RecipientId); err != nil {
			return nil, ErrUnknownEntity
		} else {
			committed.Recipient = &recipient
		}
	}
	if committed.BuyerId != "" {
		if buyer, err := idb.db.Entities().GetByNIP(committed.BuyerId); err != nil {
			return nil, ErrUnknownEntity
		} else {
			committed.Buyer = &buyer
		}
	}

	for _, item := range draft.Items {
		committed.AddItem(item)
	}

	return committed, nil
}

func (idb *InvoicesDB) GetByID(id string) (*invoice.Invoice, error) {
	now := time.Now().Local()

	table := idb.getTable(now)

	invoice, err := table.GetByIndex(InvoiceIndexId, id)
	if err != nil {
		return nil, err
	}

	// found the invoice. let's just read it back and recalculate total amounts
	invoice.CalculateTotalAmount()

	return invoice, nil
}
