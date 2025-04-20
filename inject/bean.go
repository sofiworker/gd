package inject

type Scope int

const (
	Singleton Scope = iota
	Prototype
	Request
)

type BeanDefinition struct {
	Name  string
	New   func() interface{}
	Scope Scope
}
