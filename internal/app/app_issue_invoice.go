package app

import (
	"errors"
	"os"
	"path"
	"time"

	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
	"github.com/toudi/kwity/internal/renderer"
	tmpl "github.com/toudi/kwity/internal/template"
)

type IssueParams struct {
	IssueDate    time.Time
	TemplateFile string
	Amend        bool
	Commit       bool
	Confirm      bool
	Dump         renderer.DumpParams
	InvoiceId    string
	MakeDirs     bool
}

var (
	ErrInvalidTemplate    = errors.New("invalid template")
	ErrCannotIssueInvoice = errors.New("cannot issue invoice")
	ErrCannotPrintInvoice = errors.New(
		"cannot print invoice; it was safely saved to the file though",
	)
	ErrCannotAmendCommitedInvoice = errors.New(
		"cannot amend already commited invoice. If you know what you're doing, use --confirm flag",
	)
)

type ErrSepecifyInvoiceId struct {
	Invoices []*invoice.Invoice
}

func (e *ErrSepecifyInvoiceId) Error() string {
	return "cannot uniquely identify invoice"
}

func (app *App) IssueInvoice(params IssueParams) error {
	var template *tmpl.Template
	var err error

	if template, err = tmpl.LoadFromFile(params.TemplateFile); err != nil {
		return errors.Join(ErrInvalidTemplate, err)
	}

	var inv *invoice.Invoice
	inv, err = template.PrepareInvoice(app.db)
	inv.Draft = !params.Commit

	if !params.IssueDate.IsZero() {
		inv.IssueDate = params.IssueDate
	}
	if template.SaleDate == tmpl.LastDayOfMonth {
		inv.SetSaleDateToEndOfMonth()
	}

	if err != nil {
		return errors.Join(ErrCannotIssueInvoice, err)
	}

	var invoiceDb db.InvoicesInterface = app.db.Invoices()
	var entitiesDb db.EntitiesInterface = app.db.Entities()

	// call upsertOrCreate on recipient and buyer (if any) but only if it contains a NIP
	if template.Recipient.NIP != "" {
		if err = entitiesDb.UpdateOrCreate(*template.Recipient); err != nil {
			return errors.Join(ErrCannotIssueInvoice, err)
		}
	}
	if template.Buyer != nil && template.Buyer.NIP != "" {
		if err = entitiesDb.UpdateOrCreate(*template.Buyer); err != nil {
			return errors.Join(ErrCannotIssueInvoice, err)
		}
	}

	if params.Amend {
		// let's check if there's only a single invoice for this contractor
		var invoices []*invoice.Invoice
		if invoices, err = invoiceDb.FilterByRecipient(template.Recipient.NIP, inv.IssueDate, params.InvoiceId); err != nil {
			return err
		}
		if len(invoices) > 1 {
			return &ErrSepecifyInvoiceId{Invoices: invoices}
		}
		// since we're replacing the invoice, let's use the found one to override the ID
		// though only for commited invoices - we can replace the drafts as many times as
		// we wish
		if !invoices[0].Draft && invoices[0].Number != "" && !params.Confirm {
			return ErrCannotAmendCommitedInvoice
		}
		inv.Id = invoices[0].Id
		inv.Number = invoices[0].Number
	}

	if inv.Number == "" {
		inv.Number, err = app.GetNextInvoiceNumber(inv.SaleDate, !params.Commit)
		if err != nil {
			return err
		}
	}

	if err = invoiceDb.Save(inv); err != nil {
		return err
	}

	// now let's do the printing part.
	// first, we have to locate the template
	templateFile, err := template.Metadata.GetTemplateFile(template, app.config)
	if err != nil {
		return errors.Join(ErrCannotPrintInvoice, err)
	}

	// now we can pass it to the actual renderer
	r, err := renderer.Initialize(app.config.PDFRenderer)
	if err != nil {
		return errors.Join(ErrCannotPrintInvoice, err)
	}
	targetPDFName, err := app.GetTargetPDFName(template, inv.SaleDate)
	if err != nil {
		return errors.Join(ErrCannotPrintInvoice, err)
	}
	pdfDir := path.Dir(targetPDFName)
	if _, err = os.Stat(pdfDir); os.IsNotExist(err) {
		if !params.MakeDirs {
			return errors.Join(ErrCannotPrintInvoice, err)
		}
		if err = os.MkdirAll(pdfDir, 0755); err != nil {
			return errors.Join(ErrCannotPrintInvoice, err)
		}
	}
	return r.Render(templateFile, inv, targetPDFName)
}
