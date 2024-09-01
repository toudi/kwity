package app

import (
	"errors"
	"os"

	"github.com/toudi/kwity/internal/config"
	"gopkg.in/yaml.v3"
)

func NewApp() (*App, error) {
	var app *App
	var config *config.ConfigType

	configFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, errors.New("cannot open config file")
	}
	defer configFile.Close()
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, errors.Join(errors.New("cannot decode config"), err)
	}

	if app, err = Init(config); err != nil {
		return nil, errors.Join(errors.New("error initializing app"), err)
	}

	return app, nil
}
