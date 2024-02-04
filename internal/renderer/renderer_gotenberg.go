package renderer

import (
	"github.com/toudi/kwity/internal/config"
	"github.com/toudi/kwity/internal/invoice"
)

type GotenbergRenderer struct {
	config config.GotenbergConfig
}

func (g GotenbergRenderer) Render(invoice *invoice.Invoice) error {
	return nil
}
