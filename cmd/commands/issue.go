package commands

import (
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	appPkg "github.com/toudi/kwity/internal/app"
)

type IssueCommand struct {
	appPkg.IssueParams

	Issued string `short:"d" help:"override the issue date"`
}

func (i *IssueCommand) Run(ctx *Context) error {
	if !i.Commit && i.TemplateFile == "" {
		return errors.New("template is required when issuing a draft invoice")
	}
	if i.Commit && i.InvoiceId == "" && i.TemplateFile == "" {
		return errors.New(
			"you want to create a non-draft invoice; Please specify either the template file or a draft id",
		)
	}

	var err error

	i.IssueDate = time.Now().Local()

	if i.Issued != "" {
		// try to parse the date
		if i.IssueDate, err = time.ParseInLocation("2006-01-02", i.Issued, time.Local); err != nil {
			return errors.Join(errors.New("unable to parse date"), err)
		}
	}

	if err := ctx.App.IssueInvoice(i.IssueParams); err != nil {
		var specifyInvoiceId *appPkg.ErrSepecifyInvoiceId
		if errors.As(err, &specifyInvoiceId) {
			fmt.Printf(
				"Cannot uniquely identify invoice; please use -id switch to specify the ID\n",
			)
			fmt.Printf("Available invoices:\n")
			for _, invoice := range specifyInvoiceId.Invoices {
				invoice.CalculateTotalAmount()
				fmt.Printf(
					"%s\t%.2f gross\n",
					invoice.Id,
					float64(
						invoice.TotalAmount.Gross,
					)/math.Pow10(
						invoice.TotalAmount.Multiplier,
					),
				)
			}
			os.Exit(-1)
		}
		return errors.Join(errors.New("unable to issue invoice"), err)
	}

	return nil
}
