package fs

import (
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
	"github.com/toudi/yti"
)

const entitiesDbFilename = "entities.yaml"
const entityIndexId = "id"
const entityIndexNip = "nip"

type EntitiesDB struct {
	*yti.Table[common.Entity]
}

func (f *FSDb) Entities() db.EntitiesInterface {
	return getTable(f, entitiesDbFilename, func(filename string) (*EntitiesDB, error) {
		instance, err := yti.OpenFile[common.Entity](filename, &yti.TableOptions[common.Entity]{
			Indices: map[string]yti.Indexer[common.Entity]{
				entityIndexId: func(item common.Entity) interface{} {
					return item.Id
				},
				entityIndexNip: func(item common.Entity) interface{} {
					return item.NIP
				},
			},
		})

		return &EntitiesDB{
			Table: instance,
		}, err
	})
}

func (e *EntitiesDB) UpdateOrCreate(entity common.Entity) error {
	if entity.Id == "" {
		entity.Id = entity.Name
	}

	return e.UpdateOrCreateByIndexValue(entityIndexNip, entity.NIP, entity)
}

func (e *EntitiesDB) GetByNIP(nip string) (common.Entity, error) {
	return e.GetByIndex(entityIndexNip, nip)
}
