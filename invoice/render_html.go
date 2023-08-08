package invoice

import (
	"fmt"
	"os"
	"path"

	"github.com/flosch/pongo2"
)

func (i *Invoice) renderHTML() (string, error) {
	var err error
	htmlTemplate := i.Contractor.HTMLTemplate()
	htmlTemplateDir := path.Dir(htmlTemplate)

	destFileName := path.Join(htmlTemplateDir, "render.html")
	dest, err := os.OpenFile(destFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return "", fmt.Errorf("cannot open dest file: %v", err)
	}
	defer dest.Close()

	template := pongo2.Must(pongo2.FromFile(htmlTemplate))
	err = template.ExecuteWriter(pongo2.Context{
		"invoice": map[string]interface{}{
			"number":     i.Number,
			"issue_date": i.Issued,
			"items":      i.Items,
			"total":      i.TotalAmount,
		},
	}, dest)

	return destFileName, err
}
