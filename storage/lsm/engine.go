package lsm

import (
	"fmt"
	"sync"

	errs "github.com/jasutiin/go-database/storage/errors"
)

type Engine struct {
	mutex          sync.RWMutex
	immutableTable *memTable
	currentTable   *memTable
	writeAheadLog  *wal
	sstables       []*sst
	opts           *Options
}

func Startup(opts *Options) (*Engine, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	engine := &Engine{}
	if err := engine.recover(opts); err != nil {
		return nil, err
	}

	return engine, nil
}

func (engine *Engine) recover(opts *Options) error {
	engine.opts = opts
	log, err := LoadWAL(opts)
	if err != nil {
		return err
	}

	table, err := LoadMemTable(opts, log)
	if err != nil {
		log.file.Close()
		return err
	}

	engine.currentTable = table
	engine.writeAheadLog = log
	return nil
}

func (engine *Engine) Stop() error {
	if engine.writeAheadLog == nil || engine.writeAheadLog.file == nil {
		return nil
	}

	return engine.writeAheadLog.file.Close()
}

func (engine *Engine) Get(key []byte) ([]byte, error) {
	value, tombstone, found := engine.currentTable.Get(key)
	if !found || tombstone {
		return nil, errs.ErrKeyNotFound
	}

	return value, nil
}

func (engine *Engine) Put(key, value []byte) error {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	if err := engine.writeAheadLog.Insert(key, value, false); err != nil {
		return err
	}

	err := engine.currentTable.Insert(key, value, false)

	if err != nil {
		return err
	}

	if engine.currentTable.Size() == engine.currentTable.MaxSize() {
		engine.rotateMemTables()
		go engine.flush()
	}

	return nil
}

func (engine *Engine) Delete(key []byte) error {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	if err := engine.writeAheadLog.Insert(key, nil, true); err != nil {
		return err
	}

	return engine.currentTable.Insert(key, nil, true)
}

// this method returns (new table, old table)
func (engine *Engine) rotateMemTables() {
	newTable := CreateMemTable(engine.opts)
	engine.immutableTable = engine.currentTable
	engine.currentTable = newTable
}

// TODO: make it return an error. however idk how to handle goroutines that return errors so look into that too
func (engine *Engine) flush() {
	fmt.Printf("flush not implemented")
}
