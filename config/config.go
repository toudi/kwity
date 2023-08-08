package config

type ConfigType struct {
	OrderNumberFormat string    `mapstructure:"numbering"`
	Output            string    `mapstructure:"output"`
	PDFRenderer       *Renderer `mapstructure:"renderer"`
	Templates         string    `mapstructure:"templates"`
}

var Config ConfigType
