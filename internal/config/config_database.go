package config

type FileDBConfig struct {
	Root string `mapstructure:"path"`
}

type DatabaseConfig struct {
	FileDB *FileDBConfig `mapstructure:"fs"`
}
