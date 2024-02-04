package template

import (
	"errors"

	"github.com/toudi/kwity/internal/config"
)

type InvoiceMetadata struct {
	SEI          bool   `yaml:"sei"`
	TemplateFile string `yaml:"template"`
	PDFName      string `yaml:"pdf-name"`
}

var (
	ErrTemplateFileUnset = errors.New("template file cannot be determined")
)

func (im InvoiceMetadata) GetTemplateFile(t *Template, config *config.ConfigType) (string, error) {
	if im.TemplateFile != "" {
		return im.TemplateFile, nil
	}
	if config.TemplateName != "" {
		return config.TemplateName, nil
	}

	return "", ErrTemplateFileUnset
}
