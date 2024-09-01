package config

const EnginePuppeteer string = "puppeteer"

type PuppeteerConfig struct {
	Command []string `yaml:"command"`
}

type GotenbergConfig struct {
	Host string `yaml:"host"`
}

type Renderer struct {
	Puppeteer *PuppeteerConfig `yaml:"puppeteer"`
	Gotenberg *GotenbergConfig `yaml:"gotenberg"`
}
