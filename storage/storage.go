package storage

import (
	btreeengine "github.com/jasutiin/go-database/storage/btree"
	lsmengine "github.com/jasutiin/go-database/storage/lsm"
)

type storageEngine interface {
	Stop() error
	Get(key []byte) error
	Insert(key, value []byte) error
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
		selectedEngine, err = lsmengine.Startup(&lsmengine.Options{
			DbName:           opts.DbName,
			SkipListMaxLevel: opts.SkipListMaxLevel,
			SkipListMaxSize:  opts.SkipListMaxSize,
		})

	case StorageBTree:
		selectedEngine, err = btreeengine.Startup(&btreeengine.Options{
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

func (engine *Engine) Stop() error {
	return engine.engine.Stop()
}

func (engine *Engine) Get(key []byte) error {
	return engine.engine.Get(key)
}

func (engine *Engine) Insert(key, value []byte) error {
	return engine.engine.Insert(key, value)
}

func (engine *Engine) Delete(key []byte) error {
	return engine.engine.Delete(key)
}
