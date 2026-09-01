package lsm

import errs "github.com/jasutiin/go-database/engine/errors"

type Engine struct {
	table         *memTable
	writeAheadLog *wal
	sstables      []*sst
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
	log, err := LoadWAL(opts)
	if err != nil {
		return err
	}

	table, err := LoadMemTable(opts, log)
	if err != nil {
		log.file.Close()
		return err
	}

	engine.table = table
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
	value, tombstone, found := engine.table.Get(key)
	if !found || tombstone {
		return nil, errs.ErrKeyNotFound
	}

	return value, nil
}

func (engine *Engine) Put(key, value []byte) error {
	if err := engine.writeAheadLog.Insert(key, value, false); err != nil {
		return err
	}

	return engine.table.Insert(key, value, false)
}

func (engine *Engine) Delete(key []byte) error {
	if err := engine.writeAheadLog.Insert(key, nil, true); err != nil {
		return err
	}

	return engine.table.Insert(key, nil, true)
}
