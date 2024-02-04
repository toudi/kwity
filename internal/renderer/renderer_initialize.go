package renderer

import (
	"errors"

	"github.com/toudi/kwity/internal/config"
)

var ErrUnknownRenderer = errors.New("unknown renderer")

func Initialize(config *config.Renderer) (PDFRenderer, error) {
	if config.Puppeteer != nil {
		return PuppeteerRenderer{*config.Puppeteer}, nil
	}

	return nil, ErrUnknownRenderer
}
