package renderer

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/toudi/kwity/internal/invoice"
	"gopkg.in/yaml.v3"
)

type DumpParams struct {
	JSON   bool
	YAML   bool
	Output string
}

type DumpRenderer struct {
	params DumpParams
}

func (d DumpRenderer) Render(templateFile string, i *invoice.Invoice, pdfFileName string) error {
	var output io.Writer = os.Stdout

	if d.params.Output != "" {
		output, err := os.Create(d.params.Output)
		if err != nil {
			return errors.Join(ErrRendering, err)
		}
		defer output.Close()
	}

	if d.params.JSON {
		return json.NewEncoder(output).Encode(i)
	} else if d.params.YAML {
		return yaml.NewEncoder(output).Encode(i)
	}
	return nil
}
