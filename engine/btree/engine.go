package btree

import "sync"

type Engine struct {
	mu       sync.RWMutex
	pager    *pager
	tree     *tree
	freelist *freelist
	closed   bool
}

func Startup(opts *Options) (*Engine, error) {
	return &Engine{}, nil
}

func (engine *Engine) Stop() error {
	return nil
}

func (engine *Engine) Get(key []byte) error {
	return nil
}

func (engine *Engine) Put(key, value []byte) error {
	return nil
}

func (engine *Engine) Delete(key []byte) error {
	return nil
}
