package config

const EnginePuppeteer string = "puppeteer"

type PuppeteerConfig struct {
	Command []string `mapstructure:"command"`
}

type GotenbergConfig struct {
	Host string `mapstructure:"host"`
}

type Renderer struct {
	Puppeteer *PuppeteerConfig `mapstructure:"puppeteer"`
	Gotenberg *GotenbergConfig `mapstructure:"gotenberg"`
}
