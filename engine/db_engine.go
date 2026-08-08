package engine

type Engine struct {
	memTable MemTable
	wal      WAL
}

// TODO: consider not having this act on engine
func (engine *Engine) Startup() {

}

func (engine *Engine) Stop() {

}

func (engine *Engine) Get() {

}

func (engine *Engine) Put() {

}

func (engine *Engine) Delete() {

}
