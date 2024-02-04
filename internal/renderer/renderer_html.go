package renderer

import (
	"errors"
	"fmt"

	"github.com/toudi/kwity/internal/invoice"
)

type HTMLRenderer struct{}

func (h HTMLRenderer) Render(templateFile string, i *invoice.Invoice, pdfFileName string) error {
	// let's start by pre-rendering the invoice to HTML form which can then be
	prerender, err := PrerenderInvoice(templateFile, i)
	if err != nil {
		return errors.Join(ErrRendering, err)
	}
	// now we can pass this to the actual renderer.
	fmt.Printf("rendering saved as %s\n", prerender)
	return nil
}
