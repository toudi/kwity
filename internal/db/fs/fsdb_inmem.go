package fs

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Index map[interface{}]int
type DocIndex struct {
	Name  string
	Value interface{}
}
type IndexingFunc[T any] func(document T) []DocIndex
type SaveHookFunc[T any] func(document T)
type LoadHookFunc[T any] func(document T)

type FSDBView[T any] struct {
	items    []T
	filename string
	params   *FSDBViewParams[T]
	docIndex map[string]Index
	docHash  map[string]bool
	_dirty   bool
	_empty   bool
}

type FSDBViewParams[T any] struct {
	indexer  IndexingFunc[T]
	saveHook SaveHookFunc[T]
	loadHook LoadHookFunc[T]
}

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrUnknownIndex    = errors.New("unknown index")
	ErrUnexpectedError = errors.New("unexpected error during opening")
	ErrUnableToLoadDB  = errors.New("unable to load database")
	ErrUnableToHash    = errors.New("unable to hash document")
	ErrUnableToAdd     = errors.New("unable to add document")
)

func hashDocument(document interface{}) (string, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(document); err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(buffer.Bytes())
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func FSDBView_init[T any](filename string, params *FSDBViewParams[T]) (*FSDBView[T], error) {
	db := &FSDBView[T]{
		filename: filename,
		items:    make([]T, 0),
		docIndex: make(map[string]Index),
		docHash:  make(map[string]bool),
		params:   params,
	}
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist. no worries, we can continue.
			db._empty = true
			return db, nil
		}
		// file potentially exist or it's a different kind of error
		if !os.IsNotExist(err) {
			return nil, errors.Join(ErrUnexpectedError, err)
		}
	}
	// file opened fine. let's continue.
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&db.items)
	if err != nil {
		return nil, errors.Join(ErrUnableToLoadDB, err)
	}

	for index, item := range db.items {
		docHash, err := hashDocument(item)
		if err != nil {
			return nil, errors.Join(ErrUnableToHash, err)
		}
		db.indexDocument(item, docHash, index)
	}

	return db, nil
}

func (v *FSDBView[T]) indexDocument(document T, hash string, docIndex int) {
	if v.params == nil || v.params.indexer == nil {
		return
	}

	v.docHash[hash] = true
	for _, index := range v.params.indexer(document) {
		if _, exists := v.docIndex[index.Name]; !exists {
			v.docIndex[index.Name] = make(Index)
		}
		v.docIndex[index.Name][index.Value] = docIndex
	}
}

func (v *FSDBView[T]) IndexContainsValue(indexName string, value interface{}) (bool, error) {
	var exists bool
	var err error
	if _, exists = v.docIndex[indexName]; !exists {
		err = ErrUnknownIndex
		if v._empty {
			err = nil
		}
		return false, err
	}
	_, exists = v.docIndex[indexName][value]
	return exists, nil
}

func (v *FSDBView[T]) AddDocument(document T) (bool, error) {
	if v.params != nil && v.params.saveHook != nil {
		v.params.saveHook(document)
	}
	// let's try to hash the document
	docHash, err := hashDocument(document)
	if err != nil {
		return false, errors.Join(ErrUnableToHash, err)
	}
	if v.docHash[docHash] {
		// the document already exists
		return false, nil
	}
	// add the document
	index := len(v.items)
	v.items = append(v.items, document)
	v.indexDocument(document, docHash, index)
	v._dirty = true
	v._empty = false

	return true, nil
}

func (v *FSDBView[T]) removeDocument(document T) error {
	hash, err := hashDocument(document)
	if err != nil {
		return err
	}
	delete(v.docHash, hash)
	var docIndex *int = nil

	if v.params != nil && v.params.indexer != nil {
		for _, index := range v.params.indexer(document) {
			if docIndex == nil {
				docIndex = new(int)
				*docIndex = v.docIndex[index.Name][index.Value]
			}
			delete(v.docIndex[index.Name], index.Value)
		}
		if docIndex != nil {
			newItems := make([]T, 0, len(v.items)-1)
			for index, document := range v.items {
				if index != *docIndex {
					newItems = append(newItems, document)
				}
			}
			v.items = newItems
		}
	}
	return nil
}

func (v *FSDBView[T]) UpsertDocument(lookup DocIndex, document T) error {
	// first, let's check if the document exists.
	doc, err := v.GetByIndex(lookup.Name, lookup.Value)
	if errors.Is(err, ErrItemNotFound) {
		// no worries, we can simply add the document.
		_, err = v.AddDocument(document)
		return err
	}
	// ok this means that the document is found and we have to replace it.
	// remove existing document
	if err = v.removeDocument(doc); err != nil {
		return err
	}
	// and add new one.
	_, err = v.AddDocument(document)
	return err
}

func (v *FSDBView[T]) GetByIndex(fieldName string, fieldValue interface{}) (T, error) {
	var empty T
	// first, let's check if the index is even in the db:
	fieldIndex, exists := v.docIndex[fieldName]
	if !exists {
		return empty, ErrItemNotFound
	}
	// ok, it does. now let's try to query it
	docIndex, exists := fieldIndex[fieldValue]
	if !exists {
		// no luck this time.
		return empty, ErrItemNotFound
	}
	// we have found the document
	return v.items[docIndex], nil
}

func (v *FSDBView[T]) GetByID(id string) (T, error) {
	return v.GetByIndex("id", id)
}

func (v *FSDBView[T]) ForEach(inspector func(document T)) {
	for _, document := range v.items {
		inspector(document)
	}
}

func (v *FSDBView[T]) save() error {
	if !v._dirty {
		return nil
	}
	if err := os.MkdirAll(path.Dir(v.filename), 0755); err != nil {
		return err
	}
	file, err := os.Create(v.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	v._dirty = false
	return yaml.NewEncoder(file).Encode(v.items)
}
