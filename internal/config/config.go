package config

import "github.com/toudi/kwity/internal/common"

type ConfigType struct {
	OrderNumberFormat      string           `yaml:"numbering"`
	OrderNumberFormatDraft string           `yaml:"numbering-draft"`
	Output                 string           `yaml:"output"`
	PDFRenderer            *Renderer        `yaml:"renderer"`
	PDFName                string           `yaml:"pdf-name"`
	Issuer                 *Issuer          `yaml:"issuer"`
	Templates              string           `yaml:"templates"`
	TemplateName           string           `yaml:"template-name"`
	Database               DatabaseConfig   `yaml:"database"`
	VatRates               []common.VatRate `yaml:"vat-rates"`
}
