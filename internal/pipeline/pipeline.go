package pipeline

import (
	"github.com/toudi/kwity/internal/interfaces"
	"github.com/toudi/kwity/internal/invoice"
)

type Pipeline interface {
	Process(
		app interfaces.App,
		invoice *invoice.Invoice,
		config map[string]interface{},
		context map[string]interface{},
	) error
}
