package inject

type ContainerOptionFunc func(*ContainerOption)

type ContainerOption struct {
	injectTag string
	name      string
}

type Container struct {
	opts *ContainerOption
	root *Module
}

type ContainerNode struct {
	Name string
}

func WithInjectTag(tag string) ContainerOptionFunc {
	return func(opts *ContainerOption) {
		opts.injectTag = tag
	}
}

func WithName(name string) ContainerOptionFunc {
	return func(opts *ContainerOption) {
		opts.name = name
	}
}

func New(opts ...ContainerOptionFunc) *Container {
	c := &Container{root: &Module{Name: "root"}}

	//o := &ContainerOption{
	//	injectTag: "gd",
	//}
	//for _, opt := range opts {
	//	opt(o)
	//}
	//
	//c.opts = o

	return c
}

func (c *Container) Module(name string, opts ...ModuleOptionFunc) *Module {
	o := &ModuleOption{
		injectTag: c.opts.injectTag,
	}
	for _, opt := range opts {
		opt(o)
	}
	return &Module{
		Name:   name,
		parent: c.root,
		//digContainer: dig.New(dig.RecoverFromPanics()),
		opts: o,
	}
}

func (c *Container) Provide(f interface{}) error {
	//return c.root.Provide(f)
	return nil
}
