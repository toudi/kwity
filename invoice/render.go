package invoice

import "fmt"

var RenderOptions struct {
	HTML bool
	JSON bool
}

func (i *Invoice) Render() error {
	if RenderOptions.JSON {
		return i.renderJSON()
	} else {
		renderedHtmlFile, err := i.renderHTML()
		fmt.Printf("result of renderToHTML: %v, %v\n", renderedHtmlFile, err)
		if RenderOptions.HTML {
			return err
		}
		return i.renderPDF(renderedHtmlFile)
	}
}
