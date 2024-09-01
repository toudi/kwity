package app

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/phuslu/log"
	"github.com/toudi/kwity/internal"
	"github.com/toudi/kwity/internal/invoice"
	tmpl "github.com/toudi/kwity/internal/template"
)

func (app *App) executePipelineStep(step interface{}, template *tmpl.Template,
	inv *invoice.Invoice,
	params IssueParams,
	context map[string]interface{},
) error {
	var stepName string
	var ok bool
	var stepMap map[string]interface{}

	if stepName, ok = step.(string); !ok {
		if stepMap, ok = step.(map[string]interface{}); ok {
			stepName = stepMap["name"].(string)
		} else {
			return fmt.Errorf("step is neither a string nor map[string]interface{}")
		}
	}

	if stepName == "render-pdf" {
		pdfName, err := app.RenderInvoice(template, inv, params)
		context["pdf-name"] = pdfName
		return err
	}
	// try to load the plugin
	pipelinePlugin, err := plugin.Open(
		filepath.Join("pipelines", strings.Join([]string{stepName, ".so"}, "")),
	)

	log.Trace().Err(err).Msg("plugin load")

	if err == nil {
		pipelineInstance, err := pipelinePlugin.Lookup("Pipeline")
		log.Trace().Err(err).Msg("plugin lookup")

		if err == nil {
			var pipeline internal.Pipeline
			pipeline, ok := pipelineInstance.(internal.Pipeline)
			log.Trace().Bool("ok", ok).Msg("type cast")

			if ok {
				processErr := pipeline.Process(inv, stepMap, context)
				log.Trace().Err(processErr).Msg("process")

				if processErr != nil {
					return processErr
				}
			}
		}
	}
	return err
}
