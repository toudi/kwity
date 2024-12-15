package template

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/kwity/internal/invoice"
	"github.com/toudi/kwity/internal/workdays"
)

var ErrContractorIdMissing = errors.New("contractor ID missing")
var ErrUnableToGetVATRate = errors.New("unable to get VAT rate")
var ErrUnableToGetUnitPrice = errors.New("unable to get unit price")
var ErrUnableToParseQuantity = errors.New("unable to parse quantity")

const LastDayOfMonth = "last-day-of-month"

type PrepareInvoiceOptions struct {
	IssueDate time.Time
}

func (t *Template) PrepareInvoice(
	_db db.Database,
	params PrepareInvoiceOptions,
	draft *invoice.Invoice,
) (*invoice.Invoice, error) {
	if t.Recipient.NIP == "" {
		return nil, ErrContractorIdMissing
	}
	var _invoice *invoice.Invoice = draft
	if _invoice == nil {
		_invoice = &invoice.Invoice{
			Recipient:   t.Recipient,
			Buyer:       t.Buyer,
			RecipientId: t.Recipient.NIP,
			IssueDate:   time.Now(),
			SaleDate:    time.Now(),
			Items:       make([]*invoice.Item, 0, len(t.Items)),
		}
	}
	_invoice.Items = make([]*invoice.Item, 0, len(t.Items))

	if !params.IssueDate.IsZero() {
		_invoice.IssueDate = params.IssueDate
		_invoice.SaleDate = params.IssueDate
	}

	if t.Buyer != nil {
		_invoice.BuyerId = t.Buyer.NIP
	}

	if t.SaleDate == LastDayOfMonth {
		// take the issue date and set the sale date to last day of it's month.
		_invoice.SetSaleDateToEndOfMonth()
	} else {
		// let's check if the sale date is some number in range of 1 .. 31
		day, err := strconv.Atoi(t.SaleDate)
		fmt.Printf("parsed day: %d, err=%v\n", day, err)
		if err == nil && day >= 1 && day <= 31 {
			fmt.Printf("set last day of month")
			_invoice.SetSaleDateToDayOfMonth(day)
		}
	}

	workingDays, bankHolidays := workdays.CalculateWorkingDays(
		_invoice.SaleDate.AddDate(0, 0, -_invoice.SaleDate.Day()+1), // beginning of the month
		_invoice.SaleDate,
		lo.SliceToMap(
			_db.Holidays().GetBankHolidays(),
			func(holiday string) (string, bool) { return holiday, true },
		),
	)

	var accumulatorItem *invoice.Item

	for _, item := range t.Items {
		unitPrice, err := item.GetUnitPrice()
		if err != nil {
			return nil, errors.Join(ErrUnableToGetUnitPrice, err)
		}
		vatRate, err := _db.VATRates().GetByID(item.VatRateId)
		if err != nil {
			return nil, errors.Join(ErrUnableToGetVATRate, err)
		}
		quantity, err := item.GetQuantity(workingDays, bankHolidays)
		if err != nil {
			return nil, errors.Join(ErrUnableToParseQuantity, err)
		}
		invoiceItem := &invoice.Item{
			Name:      item.Name,
			UnitPrice: unitPrice,
			VatRate:   vatRate,
			Quantity:  quantity,
		}

		_invoice.AddItem(invoiceItem)

		if item.Accumulate {
			accumulatorItem = invoiceItem
		}
	}

	if accumulatorItem != nil {
		// we want to keep the source items for the plugins, but aggregate items to a single
		// output for post-processing the invoice.
		_invoice.SourceItems = make([]*invoice.Item, len(_invoice.Items))
		copy(_invoice.SourceItems, _invoice.Items)

		accumulatedAmount := _invoice.TotalAmount

		netAmount := common.UnitPrice{
			IsGross: false,
			Price: common.PriceNormalized{
				Price:      accumulatedAmount.Net,
				Multiplier: accumulatedAmount.Multiplier,
			},
		}

		_invoice.Items = []*invoice.Item{{
			Name:      accumulatorItem.Name,
			UnitPrice: netAmount,
			VatRate:   accumulatorItem.VatRate,
			Quantity:  common.PriceNormalized{Price: 1},
		}}
	}

	return _invoice, nil
}
