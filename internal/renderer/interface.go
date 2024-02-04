package renderer

import (
	"errors"

	"github.com/toudi/kwity/internal/invoice"
)

type PDFRenderer interface {
	Render(templateFile string, invoice *invoice.Invoice, pdfFileName string) error
}

var (
	ErrRendering = errors.New("cannot render invoice")
)
