package fs

import (
	"fmt"

	"github.com/toudi/kwity/internal/config"
)

type CloseHook func() error
type FSDb struct {
	config     *config.FileDBConfig
	closeHooks []CloseHook
}

func NewFSDb(config *config.FileDBConfig) (*FSDb, error) {
	return &FSDb{config: config}, nil
}

func (f *FSDb) Close() {
	var err error

	for _, hook := range f.closeHooks {
		if err = hook(); err != nil {
			fmt.Printf("error calling %v: %v\n", hook, err)
		}
	}
}

func (f *FSDb) OnClose(hook CloseHook) {
	f.closeHooks = append(f.closeHooks, hook)
}
