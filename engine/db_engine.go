package engine

type Engine struct {
	table         *memTable
	writeAheadLog *wal
	sstables      []*sst
}

func Startup() (*Engine, error) {
	return &Engine{}, nil
}

func (engine *Engine) Stop() {

}

func (engine *Engine) Get() {

}

func (engine *Engine) Put() {

}

func (engine *Engine) Delete() {

}
