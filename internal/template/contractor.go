package template

import "github.com/toudi/kwity/internal/common"

type Entity struct {
	common.Entity
	Template  string `yaml:"template"`
	IssueDate string `yaml:"issue-date"`
	// Name            string         `yaml:"name"`
	// NIP             string         `yaml:"nip"`
	Rates           []*common.Rate `yaml:"rate"`
	PDFTemplateName string         `yaml:"pdf-name"`
}
