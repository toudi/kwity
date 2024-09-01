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
	GenerateInvoiceFromDraft(id string, issueDate time.Time) (*invoice.Invoice, error)
}

type VATRatesInterface interface {
	GetByID(id string) (common.VatRate, error)
	Sync(vatRates []common.VatRate) error
}

type EntitiesInterface interface {
	UpdateOrCreate(common.Entity) error
	GetByNIP(id string) (common.Entity, error)
}

type BankHolidaysInterface interface {
	GetBankHolidays() []string
}

type Database interface {
	Invoices() InvoicesInterface
	VATRates() VATRatesInterface
	Entities() EntitiesInterface
	BankHolidays() BankHolidaysInterface
	Close()
}
