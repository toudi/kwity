package renderer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/invoice"
)

// PrerenderInvoice populates the template with invoice properties and yields a
// HTML file that is ready to be converted to PDF
func PrerenderInvoice(templateFile string, i *invoice.Invoice) (string, error) {
	var err error
	htmlTemplateDir := path.Dir(templateFile)

	destFileName, err := filepath.Abs(path.Join(htmlTemplateDir, "render.html"))
	if err != nil {
		return "", err
	}
	dest, err := os.OpenFile(destFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("cannot open dest file: %v", err)
	}
	defer dest.Close()

	template, err := pongo2.FromFile(templateFile)
	if err != nil {
		return "", err
	}

	for _, item := range i.Items {
		amount := item.Amount()
		if item.Metadata == nil {
			item.Metadata = make(map[string]interface{})
		}
		item.Metadata["Amount"] = map[string]common.PriceNormalized{
			"Net": {
				Price:      amount.Net,
				Multiplier: amount.Multiplier,
			},
			"VAT": {
				Price:      amount.Vat,
				Multiplier: amount.Multiplier,
			},
			"Gross": {
				Price:      amount.Gross,
				Multiplier: amount.Multiplier,
			},
		}
	}

	err = template.ExecuteWriter(pongo2.Context{
		"invoice": map[string]interface{}{
			"draft":      i.Draft,
			"recipient":  i.Recipient,
			"buyer":      i.Buyer,
			"number":     i.Number,
			"issue_date": i.IssueDate,
			"sale_date":  i.SaleDate,
			"items":      i.Items,
			"total": map[string]common.PriceNormalized{
				"Gross": {
					Price:      i.TotalAmount.Gross,
					Multiplier: i.TotalAmount.Multiplier,
				},
				"Net": {
					Price:      i.TotalAmount.Net,
					Multiplier: i.TotalAmount.Multiplier,
				},
				"VAT": {
					Price:      i.TotalAmount.Vat,
					Multiplier: i.TotalAmount.Multiplier,
				},
			},
			"totalPerVATRate": i.GetTotalsPerVATRateArray(),
			"metadata":        i.Metadata,
		},
	}, dest)

	return destFileName, err
}

func renderPriceNormalized(
	value *pongo2.Value,
	param *pongo2.Value,
) (*pongo2.Value, *pongo2.Error) {
	price, ok := value.Interface().(common.PriceNormalized)
	if !ok {
		// check if it's a map that contains the values so we can convert it back
		// as common.PriceNormalized
		if tmpMap, ok := value.Interface().(map[string]interface{}); ok {
			price = common.PriceNormalized{}
			if priceAmt, ok := tmpMap["price"].(int); ok {
				price.Price = priceAmt
			}
			if multiplier, ok := tmpMap["multiplier"].(int); ok {
				price.Multiplier = multiplier
			}
		}
	}
	separator := param.String()
	return pongo2.AsValue(price.Format(separator)), nil
}

func init() {
	_ = pongo2.RegisterFilter("render_price", renderPriceNormalized)
}
