package fs

import (
	"path/filepath"

	"github.com/phuslu/log"
	"github.com/toudi/kwity/internal/config"
)

type CloseHook func() error
type FSDb struct {
	config     *config.FileDBConfig
	closeHooks []CloseHook
	tables     map[string]interface{}
}

func NewFSDb(config *config.FileDBConfig) (*FSDb, error) {
	return &FSDb{config: config, tables: make(map[string]interface{})}, nil
}

func (f *FSDb) Close() {
	var err error

	for _, hook := range f.closeHooks {
		if err = hook(); err != nil {
			log.Error().Err(err).Msgf("error calling %v", hook)
		}
	}
}

func (f *FSDb) OnClose(hook CloseHook) {
	f.closeHooks = append(f.closeHooks, hook)
}

type TableInstance interface {
	Close() error
}

func getTable[T TableInstance](
	f *FSDb,
	filename string,
	getInstance func(filename string) (T, error),
) T {
	fullPath := filepath.Join(f.config.Root, filename)
	_, exists := f.tables[fullPath]
	if !exists {
		log.Trace().Str("fullPath", fullPath).Msg("getTable.getInstance")
		instance, err := getInstance(fullPath)
		if err != nil {
			panic(err)
		}
		f.tables[fullPath] = instance
		f.OnClose(instance.Close)
	}

	return f.tables[fullPath].(T)
}
