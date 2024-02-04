package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/flosch/pongo2"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/config"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/db/fs"
	"github.com/toudi/kwity/internal/template"
)

type App struct {
	config *config.ConfigType
	db     db.Database
}

var ErrUnknownDatabaseType = errors.New("unknown database type")
var ErrCannotGenerateNumber = errors.New("cannot generate invoice sequence number")
var ErrPDFNameUndefined = errors.New("pdf-name not defined neither in config nor in invoice")
var ErrNumberingSchemeEmpty = errors.New("numbering scheme yielded empty string")

func Init(config *config.ConfigType) (*App, error) {
	app := &App{config: config}
	db, err := app.initDatabase()
	if err != nil {
		return nil, err
	}
	app.db = db

	if err = app.db.VATRates().Sync(app.config.VatRates); err != nil {
		return nil, err
	}

	return app, err
}

func (a *App) initDatabase() (db.Database, error) {

	if a.config.Database.FileDB != nil {
		return fs.NewFSDb(a.config.Database.FileDB)
	}

	return nil, ErrUnknownDatabaseType
}

func (a *App) GetNextInvoiceNumber(saleDate time.Time, draft bool) (string, error) {
	number, err := a.db.Invoices().GetNextSequenceNumber(saleDate, draft)
	if err != nil {
		return "", errors.Join(ErrCannotGenerateNumber, err)
	}
	var orderNumberTemplateSource string = a.config.OrderNumberFormat
	if draft {
		orderNumberTemplateSource = a.config.OrderNumberFormatDraft
	}
	if orderNumberTemplateSource == "" {
		return "", ErrNumberingSchemeEmpty
	}
	orderNumberTemplate, err := pongo2.FromString(orderNumberTemplateSource)
	if err != nil {
		return "", errors.Join(ErrCannotGenerateNumber, err)
	}

	out, err := orderNumberTemplate.Execute(pongo2.Context{
		"month_no": number.Month,
		"no":       number.Year,
		"year":     saleDate.Year(),
		"month":    saleDate.Month(),
	})
	if err != nil {
		return "", errors.Join(ErrCannotGenerateNumber, err)
	}

	return out, nil
}

func (a *App) GetTargetPDFName(t *template.Template, saleDate time.Time) (string, error) {
	pdfNameTemplate := t.Metadata.PDFName
	if pdfNameTemplate == "" {
		pdfNameTemplate = a.config.PDFName
	}
	if pdfNameTemplate == "" {
		return "", ErrPDFNameUndefined
	}
	pdfNameTemplate, err := common.ExecutePongo2Template(pdfNameTemplate, pongo2.Context{
		"year":     saleDate.Year(),
		"month":    saleDate.Month(),
		"month02d": fmt.Sprintf("%02d", saleDate.Month()),
	})
	if err != nil {
		return "", err
	}
	// that's just the basename though. now we can put it into the "full" path like so:
	return common.ExecutePongo2Template(a.config.Output, pongo2.Context{
		"year":    saleDate.Year(),
		"invoice": pdfNameTemplate,
	})
}

func (a *App) Close() {
	a.db.Close()
}
