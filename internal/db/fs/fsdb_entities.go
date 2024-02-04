package fs

import (
	"errors"
	"path"

	"github.com/jaevor/go-nanoid"
	"github.com/toudi/kwity/internal/common"
	"github.com/toudi/kwity/internal/db"
)

type EntitiesDB struct {
	itemsView *FSDBView[common.Entity]
}

func (f *FSDb) Entities() db.EntitiesInterface {
	dbview, err := FSDBView_init[common.Entity](
		path.Join(f.config.Root, "entities.yaml"),
		&FSDBViewParams[common.Entity]{
			indexer: func(document common.Entity) []DocIndex {
				return []DocIndex{
					{Name: "id", Value: document.Id},
				}
			},
		},
	)
	if err != nil {
		panic(err)
	}
	instance := &EntitiesDB{
		itemsView: dbview,
	}
	f.OnClose(instance.itemsView.save)
	return instance
}

func (e *EntitiesDB) UpdateOrCreate(entity common.Entity) error {
	if entity.Id == "" {
		generator, err := nanoid.Standard(4)
		if err != nil {
			return errors.Join(ErrInstantiatingIDGenerator, err)
		}
		var contains = true
		for contains {
			entity.Id = generator()
			contains, err = e.itemsView.IndexContainsValue("id", entity.Id)
			if err != nil {
				return err
			}
		}
	}

	return e.itemsView.UpsertDocument(DocIndex{Name: "id", Value: entity.Id}, entity)
}
