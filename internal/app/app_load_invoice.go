package app

import (
	"github.com/toudi/kwity/internal/invoice"
)

func (a *App) LoadInvoice(invoiceID string) (*invoice.Invoice, error) {
	return a.db.Invoices().GetByID(invoiceID)
}
