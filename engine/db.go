package engine

import (
	"github.com/jasutiin/go-database/engine/btree"
	"github.com/jasutiin/go-database/engine/lsm"
)

type storageEngine interface {
	Stop() error
	Get(key []byte) error
	Put(key, value []byte) error
	Delete(key []byte) error
}

type DB struct {
	engine storageEngine
}

func Startup(opts *Options) (*DB, error) {
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

	return &DB{engine: selectedEngine}, nil
}

func (db *DB) Stop() error {
	return db.engine.Stop()
}

func (db *DB) Get(key []byte) error {
	return db.engine.Get(key)
}

func (db *DB) Put(key, value []byte) error {
	return db.engine.Put(key, value)
}

func (db *DB) Delete(key []byte) error {
	return db.engine.Delete(key)
}
