package config

type FileDBConfig struct {
	Root string `yaml:"path"`
}

type DatabaseConfig struct {
	FileDB *FileDBConfig `yaml:"fs"`
}
