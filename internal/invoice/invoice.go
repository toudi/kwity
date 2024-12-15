package invoice

import (
	"time"

	"github.com/toudi/kwity/internal/common"
)

type Invoice struct {
	Id                string                            `yaml:"id"`
	DraftId           string                            `yaml:"draft-id,omitempty"`
	RecipientId       string                            `yaml:"recipient-id"`
	BuyerId           string                            `yaml:"buyer-id,omitempty"`
	Recipient         *common.Entity                    `yaml:"-"`
	Buyer             *common.Entity                    `yaml:"-"`
	Draft             bool                              `yaml:"draft"`
	IssueDate         time.Time                         `yaml:"issue-date"`
	SaleDate          time.Time                         `yaml:"sale-date"`
	Number            string                            `yaml:"number"`
	Items             []*Item                           `yaml:"items"`
	SourceItems       []*Item                           `yaml:"-"`
	TotalAmount       *common.Amount                    `yaml:"-"`
	_aggregatesPerVAT map[common.VatRate]*common.Amount `yaml:"-"`
	Metadata          map[string]interface{}            `yaml:"metadata,omitempty"`
}

func (i *Invoice) SetSaleDateToEndOfMonth() {
	i.SaleDate = i.IssueDate.AddDate(0, 1, -i.IssueDate.Day())
}

func (i *Invoice) SetSaleDateToDayOfMonth(day int) {
	i.SaleDate = i.IssueDate.AddDate(0, 0, day-i.IssueDate.Day())
}

func (i *Invoice) SetMetadata(key string, value interface{}) {
	if i.Metadata == nil {
		i.Metadata = make(map[string]interface{})
	}
	i.Metadata[key] = value
}
