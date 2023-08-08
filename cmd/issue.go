package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toudi/kwity/invoice"
)

var template invoice.Template

var issueCmd = &cobra.Command{
	Use:   "issue [flags] template-file",
	Short: "issue a new invoice based on input JSON spec file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := invoice.Issue(&template, args[0]); err != nil {
			fmt.Printf("error issuing invoice: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	issueCmd.Flags().StringVarP(&template.CurrentDate, "date", "d", "", "Override running date")
	issueCmd.Flags().BoolVarP(&invoice.RenderOptions.HTML, "html", "", false, "Render to HTML")
	issueCmd.Flags().BoolVarP(&invoice.RenderOptions.JSON, "json", "", false, "Render to JSON")
	rootCmd.AddCommand(issueCmd)
}
