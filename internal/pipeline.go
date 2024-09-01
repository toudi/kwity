package internal

import "github.com/toudi/kwity/internal/invoice"

type Pipeline interface {
	Process(
		invoice *invoice.Invoice,
		config map[string]interface{},
		context map[string]interface{},
	) error
}
