package commands

import (
	"errors"
	"fmt"

	"github.com/toudi/kwity/internal/app"
	"github.com/toudi/kwity/internal/invoice"
	tmpl "github.com/toudi/kwity/internal/template"
)

type RenderCommand struct {
	InvoiceId    string `help:"Specify invoice ID" short:"i"`
	TemplateFile string `help:"template to use" type:"existingfile" arg:""`
}

func (r *RenderCommand) Run(ctx *Context) error {
	var template *tmpl.Template
	var inv *invoice.Invoice
	var err error

	if template, err = tmpl.LoadFromFile(r.TemplateFile); err != nil {
		return errors.Join(app.ErrInvalidTemplate, err)
	}

	if inv, err = ctx.App.LoadInvoice(r.InvoiceId); err != nil {
		return errors.New("invalid invoice ID")
	}

	fmt.Printf("metadata: %+v\n", inv.Metadata)

	pdfName, err := ctx.App.RenderInvoice(
		template,
		inv,
		app.IssueParams{},
	)

	if err != nil {
		return err
	}

	fmt.Printf("rendered invoice as %s\n", pdfName)
	return nil
}
