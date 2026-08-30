package engine

import "github.com/jasutiin/go-database/engine/lsm"

type storageEngine interface {
	Stop() error
	Get(key []byte) error
	Put(key, value []byte) error
	Delete(key []byte) error
}

type Engine struct {
	storage storageEngine
}

func Startup(opts *Options) (*Engine, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	storage, err := lsm.Startup(&lsm.Options{
		DbName:           opts.DbName,
		SkipListMaxLevel: opts.SkipListMaxLevel,
		SkipListMaxSize:  opts.SkipListMaxSize,
	})
	if err != nil {
		return nil, err
	}

	return &Engine{storage: storage}, nil
}

func (engine *Engine) Stop() error {
	return engine.storage.Stop()
}

func (engine *Engine) Get(key []byte) error {
	return engine.storage.Get(key)
}

func (engine *Engine) Put(key, value []byte) error {
	return engine.storage.Put(key, value)
}

func (engine *Engine) Delete(key []byte) error {
	return engine.storage.Delete(key)
}
