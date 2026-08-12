package engine

type Engine struct {
	table         *memTable
	writeAheadLog *wal
	sstables      []*sst
}

func Startup() (*Engine, error) {
	table := NewMemTable()
	return &Engine{table: table}, nil
}

func (engine *Engine) Stop() {

}

func (engine *Engine) Get() {

}

func (engine *Engine) Put() {

}

func (engine *Engine) Delete() {

}
