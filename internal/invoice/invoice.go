package invoice

import (
	"time"

	"github.com/toudi/kwity/internal/common"
)

type Invoice struct {
	Id                string                            `yaml:"id"`
	RecipientId       string                            `yaml:"recipient-id"`
	BuyerId           string                            `yaml:"buyer-id,omitempty"`
	Recipient         *common.Entity                    `yaml:"-"`
	Buyer             *common.Entity                    `yaml:"-"`
	Draft             bool                              `yaml:"draft"`
	IssueDate         time.Time                         `yaml:"issue-date"`
	SaleDate          time.Time                         `yaml:"sale-date"`
	Number            string                            `yaml:"number"`
	Items             []*Item                           `yaml:"items"`
	TotalAmount       *common.Amount                    `yaml:"-"`
	_aggregatesPerVAT map[common.VatRate]*common.Amount `yaml:"-"`
}

func (i *Invoice) targetPDFName() string {
	// filenameTemplate := i.Contractor.PDFTemplateName

	// replacer := strings.NewReplacer(
	// 	"{{ year }}", fmt.Sprint(i.Issued.Year()),
	// 	"{{ month }}", fmt.Sprintf("%02d", i.Issued.Month()),
	// )

	// fileNameBase := replacer.Replace(filenameTemplate)

	// replacer = strings.NewReplacer(
	// 	"{{ year }}", fmt.Sprint(i.Issued.Year()),
	// 	"{{ month }}", fmt.Sprintf("%02d", i.Issued.Month()),
	// 	"{{ invoice }}", fileNameBase,
	// )

	return ""

	// return replacer.Replace(config.Config.Output)
}

func (i *Invoice) SetSaleDateToEndOfMonth() {
	i.SaleDate = i.IssueDate.AddDate(0, 1, -i.IssueDate.Day())
}
