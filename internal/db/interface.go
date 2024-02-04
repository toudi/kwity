package db

import (
	"time"

	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/invoice"
)

type InvoicesInterface interface {
	Save(*invoice.Invoice) error
	FilterByRecipient(
		recipientID string,
		issued time.Time,
		id string,
	) ([]*invoice.Invoice, error)
	GetNextSequenceNumber(saleDate time.Time, draft bool) (common.SequenceNumber, error)
}

type VATRatesInterface interface {
	GetByID(id string) (common.VatRate, error)
	Sync(vatRates []common.VatRate) error
}

type EntitiesInterface interface {
	UpdateOrCreate(common.Entity) error
}

type Database interface {
	Invoices() InvoicesInterface
	VATRates() VATRatesInterface
	Entities() EntitiesInterface
	Close()
}
