package storage

import (
	"github.com/jasutiin/go-database/storage/btree"
	"github.com/jasutiin/go-database/storage/lsm"
)

type storageEngine interface {
	Stop() error
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
}

type Engine struct {
	engine storageEngine
}

func Startup(opts *Options) (*Engine, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	var selectedEngine storageEngine
	var err error

	switch opts.StorageType {
	case StorageLSM:
		selectedEngine, err = lsm.Startup(&lsm.Options{
			DbName:           opts.DbName,
			SkipListMaxLevel: opts.SkipListMaxLevel,
			SkipListMaxSize:  opts.SkipListMaxSize,
		})

	case StorageBTree:
		selectedEngine, err = btree.Startup(&btree.Options{
			DbName: opts.DbName,
		})

	default:
		return nil, ErrUnknownStorageType
	}

	if err != nil {
		return nil, err
	}

	return &Engine{engine: selectedEngine}, nil
}

func (storage *Engine) Stop() error {
	return storage.engine.Stop()
}

func (storage *Engine) Get(key []byte) ([]byte, error) {
	return storage.engine.Get(key)
}

func (storage *Engine) Put(key, value []byte) error {
	return storage.engine.Put(key, value)
}

func (storage *Engine) Delete(key []byte) error {
	return storage.engine.Delete(key)
}
