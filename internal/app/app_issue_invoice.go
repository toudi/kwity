package app

import (
	"errors"
	"time"

	"github.com/phuslu/log"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
	tmpl "github.com/toudi/kwity/internal/template"
)

type IssueParams struct {
	IssueDate    time.Time `kong:"-"`
	TemplateFile string    `         help:"template to use"                                              type:"existingfile" arg:""`
	Ammend       bool      `         help:"regenerate the invoice"`
	Commit       bool      `         help:"commit the invoice (create a regular invoice based on draft)"                            short:"c"`
	Confirm      bool      `         help:"confirm ammending a non-draft invoice that has a number"`
	InvoiceId    string    `         help:"specify invoice ID"                                                                      short:"i"`
	MakeDirs     bool      `         help:"create missing directories on the output path"                                           short:"m"`
}

var (
	ErrInvalidTemplate            = errors.New("invalid template")
	ErrCannotIssueInvoice         = errors.New("cannot issue invoice")
	ErrCannotAmendCommitedInvoice = errors.New(
		"cannot amend already commited invoice. If you know what you're doing, use --confirm flag",
	)
	ErrCannotConvertADraft = errors.New(
		"either the specified ID is not a draft, or it does not exist or it already has a matching invoice",
	)
	ErrCannotAmendADraft = errors.New("cannot ammend the draft")
)

type ErrSepecifyInvoiceId struct {
	Invoices []*invoice.Invoice
}

func (e *ErrSepecifyInvoiceId) Error() string {
	return "cannot uniquely identify invoice"
}

func (app *App) IssueInvoice(params IssueParams) error {
	var template *tmpl.Template
	var inv *invoice.Invoice
	var err error
	var invoiceDb db.InvoicesInterface = app.db.Invoices()
	var entitiesDb db.EntitiesInterface = app.db.Entities()

	if template, err = tmpl.LoadFromFile(params.TemplateFile); err != nil {
		return errors.Join(ErrInvalidTemplate, err)
	}

	if params.InvoiceId != "" {
		// we're either generating a draft or ammending existing invoice.
		if params.Ammend {
			// ok, we're ammending it
			// let's check if there's only a single invoice for this contractor
			var invoices []*invoice.Invoice
			if invoices, err = invoiceDb.FilterByRecipient(template.Recipient.NIP, params.IssueDate, params.InvoiceId); err != nil {
				return err
			}

			// is there even a single invoice?
			if len(invoices) > 0 {
				if len(invoices) > 1 {
					return &ErrSepecifyInvoiceId{Invoices: invoices}
				}

				// since we're replacing the invoice, let's use the found one to override the ID
				// though only for commited invoices - we can replace the drafts as many times as
				// we wish
				if !invoices[0].Draft && invoices[0].Number != "" && !params.Confirm {
					return ErrCannotAmendCommitedInvoice
				}

				inv = invoices[0]
				inv, err = template.PrepareInvoice(app.db, tmpl.PrepareInvoiceOptions{
					IssueDate: params.IssueDate,
				}, inv)
				if err != nil {
					return errors.Join(ErrCannotAmendADraft, err)
				}
			}
		} else {
			// we want to generate an invoice based on a draft.
			inv, err = invoiceDb.GenerateInvoiceFromDraft(params.InvoiceId, params.IssueDate)
			if err != nil {
				return errors.Join(ErrCannotConvertADraft, err)
			}
			// inv.DraftId = params.InvoiceId
		}
	} else {
		inv, err = template.PrepareInvoice(app.db, tmpl.PrepareInvoiceOptions{
			IssueDate: params.IssueDate,
		}, nil)
		inv.Draft = !params.Commit

		if err != nil {
			return errors.Join(ErrCannotIssueInvoice, err)
		}
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
	}

	if inv.Number == "" {
		inv.Number, err = app.GetNextInvoiceNumber(inv.SaleDate, !params.Commit)
		if err != nil {
			return err
		}
	}

	if err = invoiceDb.Save(inv); err != nil {
		log.Error().Err(err).Msg("error saving invoice")
		return err
	}

	log.Debug().Msg("invoice saved successfully")

	// process the invoice trough pipeline(s)
	var backupInvoice *invoice.Invoice = &invoice.Invoice{}
	// back up the invoice so that the pipeline step(s) do not
	// modify it
	*backupInvoice = *inv

	var context = make(map[string]interface{})

	for _, step := range template.Metadata.PipelineSteps {
		if err = app.executePipelineStep(step, template, inv, params, context); err != nil {
			return err
		}
	}

	// restore invoice from backup.
	inv = backupInvoice
	// finally, print the invoice.
	_, err = app.RenderInvoice(template, inv, params)
	return err
}
