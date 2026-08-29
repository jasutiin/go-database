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

func (engine *Engine) Stop() error {
	fmt.Println("Stop() not implemented")
	return nil
}

func (engine *Engine) Get(key []byte) error {
	fmt.Println("Get() not implemented")
	return nil
}

func (engine *Engine) Put(key []byte, val []byte) error {
	err := engine.writeAheadLog.Insert(key, val, false)
	if err != nil {
		return err
	}

	err = engine.table.Insert(key, val, false)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) Delete(key []byte) error {
	x := make([]byte, 3) // TODO: replace this later
	err := engine.writeAheadLog.Insert(key, x, true)
	if err != nil {
		return err
	}

	err = engine.table.Insert(key, x, true)
	if err != nil {
		return err
	}

	return nil
}
