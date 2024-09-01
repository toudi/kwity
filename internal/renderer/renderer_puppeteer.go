package renderer

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/phuslu/log"
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

	log.Trace().Msgf("executing %s %v", command, commandArgs)

	cmd := exec.Command(command, commandArgs...)
	b := new(strings.Builder)
	cmd.Stderr = b

	if err := cmd.Run(); err != nil {
		log.Trace().Any("error from stderr", b).Msg("")
		return errors.Join(errors.New("error during conversion to PDF"), err)
	}

	return nil
}
