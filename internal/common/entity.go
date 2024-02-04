package common

type Entity struct {
	Id      string   `yaml:"id"`
	Name    string   `yaml:"name"`
	NIP     string   `yaml:"nip"`
	Address []string `yaml:"address"`
}
