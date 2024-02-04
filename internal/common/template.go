package common

import "github.com/flosch/pongo2"

func ExecutePongo2Template(source string, context pongo2.Context) (string, error) {
	template, err := pongo2.FromString(source)
	if err != nil {
		return "", err
	}
	return template.Execute(context)
}
