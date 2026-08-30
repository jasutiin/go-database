package btree

type Engine struct{}

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
