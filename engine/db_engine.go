package engine

import (
	"fmt"
)

type Engine struct {
	table         *memTable
	writeAheadLog *wal
	sstables      []*sst
}

func Startup(opts *Options) (*Engine, error) {
	err := opts.Validate()

	if err != nil {
		return nil, err
	}

	engine := &Engine{}
	err = engine.recover(opts)

	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (engine *Engine) recover(opts *Options) error {
	log, err := LoadWAL(opts)

	if err != nil {
		return err
	}

	table, err := LoadMemTable(opts, log)

	if err != nil {
		return err
	}

	engine.table = table
	engine.writeAheadLog = log
	return nil
}

func (engine *Engine) Stop() {
	fmt.Println("Stop() not implemented")
}

func (engine *Engine) Get(key []byte) {
	fmt.Println("Get() not implemented")
}

func (engine *Engine) Put(key []byte, val []byte) {
	fmt.Println("Put() not implemented")
}

func (engine *Engine) Delete(key []byte) {
	fmt.Println("Delete() not implemented")
}
