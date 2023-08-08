package config

const EnginePuppeteer string = "puppeteer"

type Renderer struct {
	Engine  string   `mapstructure:"engine"`
	Command []string `mapstructure:"command"`
}
