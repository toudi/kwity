package config

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func CastConfig(src map[string]interface{}, dest interface{}) error {
	var buffer bytes.Buffer

	if err := yaml.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return yaml.NewDecoder(&buffer).Decode(dest)
}
