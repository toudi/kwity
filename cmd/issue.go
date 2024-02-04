package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	appPkg "github.com/toudi/kwity/internal/app"

	"github.com/spf13/cobra"
)

var params appPkg.IssueParams
var issueDate string

var issueCmd = &cobra.Command{
	Use:   "issue [flags] template-file",
	Short: "issue a new invoice based on input JSON spec file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params.TemplateFile = args[0]
		var err error
		if issueDate != "" {
			// try to parse the date
			if params.IssueDate, err = time.ParseInLocation("2006-01-02", issueDate, time.Local); err != nil {
				fmt.Printf("unable to parse date: %v", err)
				os.Exit(1)
			}
		}
		if err := app.IssueInvoice(params); err != nil {
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
				os.Exit(1)
			}
			fmt.Printf("error issuing invoice: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	issueCmd.Flags().BoolVarP(&params.Amend, "amend", "a", false, "Re-generate the invoice")
	issueCmd.Flags().BoolVarP(&params.Commit, "commit", "c", false, "Commit the invoice")
	issueCmd.Flags().
		BoolVarP(&params.Confirm, "confirm", "", false, "Confirm amending invoice that has a number")
	issueCmd.Flags().
		StringVarP(&params.InvoiceId, "invoice-id", "i", "", "Specify the invoice ID for override")
	issueCmd.Flags().StringVarP(&issueDate, "date", "d", "", "override the issue date")
	issueCmd.Flags().
		BoolVarP(&params.MakeDirs, "mkdirs", "m", false, "Create missing directories on the output path")
	// issueCmd.Flags().StringVarP(&template.CurrentDate, "date", "d", "", "Override running date")
	// issueCmd.Flags().BoolVarP(&invoice.RenderOptions.HTML, "html", "", false, "Render to HTML")
	// issueCmd.Flags().BoolVarP(&invoice.RenderOptions.JSON, "json", "", false, "Render to JSON")
	rootCmd.AddCommand(issueCmd)
}
