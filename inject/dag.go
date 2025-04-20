package inject

type Dependencies struct {
	nodes map[string]*BeanDefinition
	edges map[string][]string
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		nodes: make(map[string]*BeanDefinition),
		edges: make(map[string][]string),
	}
}

func (d *Dependencies) Add(name string, def *BeanDefinition) {
	d.nodes[name] = def
}
