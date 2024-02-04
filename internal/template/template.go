package template

import (
	"errors"

	"github.com/toudi/kwity/internal/common"
)

var ErrTemplateParseError = errors.New("unable to parse template")
var ErrInvalidQuantity = errors.New("quantity is an unknown")

type RateDefinition struct {
	StartDate string `yaml:"start-date" json:"start-date"`
	common.Rate
}

type Template struct {
	Recipient *common.Entity    `yaml:"recipient"`
	Buyer     *common.Entity    `yaml:"buyer"`
	Items     []*TemplateItem   `yaml:"items"`
	Rates     []*RateDefinition `yaml:"rates"`
	Sei       bool              `yaml:"sei"`
	SaleDate  string            `yaml:"sale-date,omitempty"`
	Metadata  InvoiceMetadata   `yaml:"metadata,omitempty"`
}
