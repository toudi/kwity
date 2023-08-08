package invoice

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

const quantityNumWorkdays = "\"month-workdays\""

type TemplateItem struct {
	Name           string          `json:"name"`
	Quantity       json.RawMessage `json:"quantity"`
	QuantityNumber float64
	UnitPrice      UnitPrice `json:"unit_price"`
	DaysFree       int       `json:"days_free"`
	Vat            *Vat      `json:"vat"`
}

type RateDefinition struct {
	StartDate string  `json:"start_date"`
	Rate      float64 `json:"rate"`
}

type Template struct {
	Contractor  *Contractor     `json:"contractor"`
	Items       []*TemplateItem `json:"items"`
	CurrentDate string
}

func ParseTemplate(srcFile string, dest *Template) error {
	templateFile, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("error opening template: %v", err)
	}
	defer templateFile.Close()
	jsonReader := json.NewDecoder(templateFile)
	if err = jsonReader.Decode(&dest); err != nil {
		return fmt.Errorf("could not parse template file: %v", err)
	}
	// let's check if the template is defined
	htmlTemplate := dest.Contractor.HTMLTemplate()
	if _, err = os.Stat(htmlTemplate); htmlTemplate == "" || os.IsNotExist(err) {
		return errInvoiceTemplateNotDefined
	}
	// let's check if the quantity is ok
	for i, item := range dest.Items {
		if string(item.Quantity) != quantityNumWorkdays {
			// let's check if it's a valid number.
			if item.QuantityNumber, err = strconv.ParseFloat(string(item.Quantity), 64); err != nil {
				return fmt.Errorf("quantity at item %d is neither a known string or a float: %v", i, err)
			}
		}
	}
	// parse rates
	if err = dest.Contractor.parseRates(); err != nil {
		fmt.Printf("error parsing rates: %v\n", err)
		return errInvalidRates
	}
	return nil
}
