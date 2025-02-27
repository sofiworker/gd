package gd

type Engine struct {
}

func Default() *Engine {
	return &Engine{}
}

func (e *Engine) Run() error {
	return nil
}
