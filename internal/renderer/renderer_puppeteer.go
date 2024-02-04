package renderer

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/toudi/kwity/internal/config"
	"github.com/toudi/kwity/internal/invoice"
)

type PuppeteerRenderer struct {
	config config.PuppeteerConfig
}

func (p PuppeteerRenderer) Render(
	templateFile string,
	invoice *invoice.Invoice,
	pdfFileName string,
) error {
	prerender, err := PrerenderInvoice(templateFile, invoice)
	if err != nil {
		return errors.Join(ErrRendering, err)
	}

	command := p.config.Command[0]
	commandArgs := p.config.Command[1:]
	commandArgs = append(commandArgs, prerender, pdfFileName)

	fmt.Printf("executing %v %v\n", command, commandArgs)

	cmd := exec.Command(command, commandArgs...)
	b := new(strings.Builder)
	cmd.Stderr = b

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error from stderr: %v\n", b)
		return fmt.Errorf("error during conversion to PDF: %v", err)
	}

	return nil
}
