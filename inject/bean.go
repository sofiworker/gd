package inject

type Scope int

const (
	Singleton Scope = iota
	Prototype
	Request
)

type BeanDefinition struct {
	Name  string
	New   interface{}
	Scope Scope
}
