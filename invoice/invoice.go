package invoice

import (
	"fmt"
	"strings"
	"time"

	"github.com/toudi/kwity/config"
)

type Invoice struct {
	Contractor        *Contractor       `json:"-"`
	Workdays          int               `json:"-"`
	Issued            time.Time         `json:"issued"`
	Number            string            `json:"number"`
	Items             []*Item           `json:"items"`
	TotalAmount       Amount            `json:"total"`
	TotalAmountPerVAT map[string]Amount `json:"total-per-vat"`
}

type UnitPrice struct {
	Price float32 `json:"price"`
	Gross bool    `json:"gross"`
}

type Item struct {
	Name      string    `json:"name"`
	Quantity  float32   `json:"quantity"`
	UnitPrice UnitPrice `json:"unit_price"`
	Amount    Amount    `json:"amount"`
	Vat       *Vat      `json:"vat"`
}

func (i *Invoice) targetPDFName() string {
	filenameTemplate := i.Contractor.PDFTemplateName

	replacer := strings.NewReplacer(
		"{{ year }}", fmt.Sprint(i.Issued.Year()),
		"{{ month }}", fmt.Sprintf("%02d", i.Issued.Month()),
	)

	fileNameBase := replacer.Replace(filenameTemplate)

	replacer = strings.NewReplacer(
		"{{ year }}", fmt.Sprint(i.Issued.Year()),
		"{{ month }}", fmt.Sprintf("%02d", i.Issued.Month()),
		"{{ invoice }}", fileNameBase,
	)

	return replacer.Replace(config.Config.Output)
}
