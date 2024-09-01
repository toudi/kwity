package app

import (
	"errors"
	"os"
	"path"

	"github.com/toudi/kwity/internal/invoice"
	"github.com/toudi/kwity/internal/renderer"
	tmpl "github.com/toudi/kwity/internal/template"
)

var (
	ErrCannotPrintInvoice = errors.New(
		"cannot print invoice",
	)
	ErrDirectoryDoesNotExistMissingMkdirsOption = errors.New(
		"the target directory for PDF does not exist and you did not specify -m / --make-dirs option",
	)
)

func (app *App) RenderInvoice(
	template *tmpl.Template,
	inv *invoice.Invoice,
	params IssueParams,
) (string, error) {
	if app.rendered {
		return "", nil
	}
	// first, we have to locate the template
	templateFile, err := template.Metadata.GetTemplateFile(template, app.config)
	if err != nil {
		return "", errors.Join(ErrCannotPrintInvoice, err)
	}

	// now we can pass it to the actual renderer
	r, err := renderer.Initialize(app.config.PDFRenderer)
	if err != nil {
		return "", errors.Join(ErrCannotPrintInvoice, err)
	}
	targetPDFName, err := app.GetTargetPDFName(template, inv.SaleDate)
	if err != nil {
		return "", errors.Join(ErrCannotPrintInvoice, err)
	}
	pdfDir := path.Dir(targetPDFName)
	if _, err = os.Stat(pdfDir); os.IsNotExist(err) {
		if !params.MakeDirs {
			return "", errors.Join(
				ErrCannotPrintInvoice,
				ErrDirectoryDoesNotExistMissingMkdirsOption,
				err,
			)
		}
		if err = os.MkdirAll(pdfDir, 0755); err != nil {
			return "", errors.Join(ErrCannotPrintInvoice, err)
		}
	}
	app.rendered = true
	return targetPDFName, r.Render(templateFile, inv, targetPDFName)
}
