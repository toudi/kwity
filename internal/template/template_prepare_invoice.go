package template

import (
	"errors"
	"time"

	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
)

var ErrContractorIdMissing = errors.New("contractor ID missing")
var ErrUnableToGetVATRate = errors.New("unable to get VAT rate")
var ErrUnableToGetUnitPrice = errors.New("unable to get unit price")
var ErrUnableToParseQuantity = errors.New("unable to parse quantity")

const LastDayOfMonth = "last-day-of-month"

func (t *Template) PrepareInvoice(_db db.Database) (*invoice.Invoice, error) {
	if t.Recipient.NIP == "" {
		return nil, ErrContractorIdMissing
	}
	_invoice := &invoice.Invoice{
		Recipient:   t.Recipient,
		Buyer:       t.Buyer,
		RecipientId: t.Recipient.NIP,
		IssueDate:   time.Now(),
		SaleDate:    time.Now(),
		Items:       make([]*invoice.Item, 0, len(t.Items)),
	}

	if t.Buyer != nil {
		_invoice.BuyerId = t.Buyer.NIP
	}

	if t.SaleDate == LastDayOfMonth {
		// take the issue date and set the sale date to last day of it's month.
		_invoice.SetSaleDateToEndOfMonth()
	}

	for _, item := range t.Items {
		unitPrice, err := item.GetUnitPrice()
		if err != nil {
			return nil, errors.Join(ErrUnableToGetUnitPrice, err)
		}
		vatRate, err := _db.VATRates().GetByID(item.VatRateId)
		if err != nil {
			return nil, errors.Join(ErrUnableToGetVATRate, err)
		}
		quantity, err := item.GetQuantity()
		if err != nil {
			return nil, errors.Join(ErrUnableToParseQuantity, err)
		}
		invoiceItem := &invoice.Item{
			Name:      item.Name,
			UnitPrice: unitPrice,
			VatRate:   vatRate,
			Quantity:  quantity,
		}

		_invoice.AddItem(invoiceItem)
	}

	return _invoice, nil
}
