package invoice

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/toudi/kwity/config"
)

func (i *Invoice) renderPDF(renderedHtmlFile string) error {
	defer os.Remove(renderedHtmlFile)

	_config := config.Config

	targetPDFName := i.targetPDFName()

	if _, err := os.Stat(path.Dir(targetPDFName)); os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(targetPDFName), 0755); err != nil {
			return fmt.Errorf("cannot create target directroy for PDF's: %v", err)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not detect cwd: %v", err)
	}

	if _config.PDFRenderer.Engine == config.EnginePuppeteer {
		command := _config.PDFRenderer.Command[0]
		commandArgs := _config.PDFRenderer.Command[1:]

		// we could treat with pongo2 at later time but for time being
		// let's keep things simple - these are the only 2 supported
		// parameters.
		replacer := strings.NewReplacer(
			"{{ source_file }}", path.Join(cwd, renderedHtmlFile),
			"{{ target_pdf }}", targetPDFName,
		)

		for i, arg := range commandArgs {
			commandArgs[i] = replacer.Replace(arg)
		}

		fmt.Printf("executing %v %v\n", command, commandArgs)

		cmd := exec.Command(command, commandArgs...)
		b := new(strings.Builder)
		cmd.Stderr = b

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error from stderr: %v\n", b)
			return fmt.Errorf("error during conversion to PDF: %v", err)
		}
	}

	fmt.Printf("Rendered invoice as %s\n", targetPDFName)

	return nil
}
