package invoice

import (
	"errors"
	"path"

	"github.com/toudi/kwity/config"
)

var errInvoiceTemplateNotDefined = errors.New("template file not defined or does not exist")

func (c *Contractor) HTMLTemplate() string {
	// if a contractor has a template field populated - let's use that
	if c.Template != "" {
		return c.Template
	}
	// otherwise, let's see if there's a file within config's templates directory
	// + contractor ID + invoice.html
	return path.Join(config.Config.Templates, c.Id, "invoice.html")
}
