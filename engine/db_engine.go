package engine

import (
	"fmt"
)

type Engine struct {
	table         *memTable
	writeAheadLog *wal
	sstables      []*sst
}

func Startup() (*Engine, error) {
	engine := &Engine{}
	err := engine.recover()

	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (engine *Engine) recover() error {
	log, err := LoadWAL()

	if err != nil {
		return err
	}

	table, err := LoadMemTable()

	if err != nil {
		return err
	}

	engine.table = table
	engine.writeAheadLog = log
	return nil
}

func (engine *Engine) Stop() {

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
