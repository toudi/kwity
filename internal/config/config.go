package config

import "github.com/toudi/kwity/internal/common"

type ConfigType struct {
	OrderNumberFormat      string           `mapstructure:"numbering"`
	OrderNumberFormatDraft string           `mapstructure:"numbering-draft"`
	Output                 string           `mapstructure:"output"`
	PDFRenderer            *Renderer        `mapstructure:"renderer"`
	PDFName                string           `mapstructure:"pdf-name"`
	Issuer                 *Issuer          `mapstructure:"issuer"`
	Templates              string           `mapstructure:"templates"`
	TemplateName           string           `mapstructure:"template-name"`
	Database               DatabaseConfig   `mapstructure:"database"`
	VatRates               []common.VatRate `mapstructure:"vat-rates"`
}
