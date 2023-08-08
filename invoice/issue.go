package invoice

import (
	"fmt"
	"time"

	"github.com/toudi/kwity/workdays"
)

func Issue(template *Template, srcFile string) error {
	var err error

	if err = ParseTemplate(srcFile, template); err != nil {
		return fmt.Errorf("error parsing template %v", err)
	}

	invoice := &Invoice{}
	invoice.Issued = time.Now()
	if template.CurrentDate != "" {
		invoice.Issued, err = time.ParseInLocation("2006-01-02", template.CurrentDate, time.Local)
		if err != nil {
			return fmt.Errorf("unable to parse issued date: %v", err)
		}
	}
	beginningOfLastMonth := invoice.Issued.AddDate(0, -1, -invoice.Issued.Day()+1)
	beginningOfCurrentMonth := invoice.Issued.AddDate(0, 0, -invoice.Issued.Day()+1)
	endOfLastMonth := invoice.Issued.AddDate(0, 0, -invoice.Issued.Day())
	invoice.Workdays = workdays.CalculateWorkingDays(beginningOfLastMonth, beginningOfCurrentMonth)

	var invoiceItem *Item

	var contractor *Contractor = template.Contractor

	invoice.Contractor = contractor

	if contractor.IssueDate == IssueDateEOM {
		invoice.Issued = endOfLastMonth
	}

	if invoice.Number, err = getNextInvoiceNumber(invoice.Issued, contractor.Id); err != nil {
		return fmt.Errorf("unable to obtain next invoice number: %v", err)
	}

	invoice.Items = make([]*Item, 0)
	for _, templateItem := range template.Items {
		if string(templateItem.Quantity) == quantityNumWorkdays {
			daysWithRate, err := contractor.WorkdaysWithRates(beginningOfLastMonth, endOfLastMonth)
			if err != nil {
				return fmt.Errorf("could not obtain list of workdays: %v", err)
			}
			for i, daysWithRateItem := range daysWithRate {
				invoiceItem = &Item{Name: templateItem.Name}

				invoiceItem.Quantity = float32(daysWithRateItem.Workdays)
				// yes, I get it. it's a little inconvenient that we always subtract the free days
				// from the last item because if your rate changes (hopefully increases) near the end
				// of the month it means you'll lose more. I was also considering to add a parameter
				// that would allow you to specify specific days you were not working, but for the time
				// being I won't bother with this.
				if templateItem.DaysFree > 0 && i == len(daysWithRate)-1 {
					invoiceItem.Quantity -= float32(templateItem.DaysFree)
				}
				invoiceItem.UnitPrice.Price = daysWithRateItem.Rate.Rate
				invoiceItem.UnitPrice.Gross = daysWithRateItem.Rate.IsGross
				invoiceItem.Vat = daysWithRateItem.Rate.Vat
				invoice.addItem(invoiceItem)
			}
		}
		if templateItem.QuantityNumber > 0 {
			invoiceItem = &Item{
				Name:     templateItem.Name,
				Quantity: float32(templateItem.QuantityNumber),
				UnitPrice: UnitPrice{
					Price: float32(templateItem.UnitPrice.Price),
					Gross: templateItem.UnitPrice.Gross,
				},
				Vat: templateItem.Vat,
			}
			invoice.addItem(invoiceItem)
		}
	}

	return invoice.Render()
}
