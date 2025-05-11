package inject

import (
	"sync"
)

type ProvideOption func()

type ModuleOption struct {
	injectTag string
}

type ModuleOptionFunc func(opts *ModuleOption)

type Module struct {
	Name        string
	Parent      *Module
	opts        *ModuleOption
	locker      sync.RWMutex
	definitions map[string]*BeanDefinition
	dag         *Dependencies
	modules     []*Module
}

//func NewModule(name string, opts ...ModuleOptionFunc) *Module {
//	o := &ModuleOption{
//		injectTag: "gd",
//	}
//	for _, opt := range opts {
//		opt(o)
//	}
//	return &Module{
//		Name:        name,
//		Parent:      nil,
//		opts:        o,
//		locker:      &sync.RWMutex{},
//		definitions: make(map[string]*BeanDefinition),
//		dag:         NewDependencies(),
//	}
//}

func (m *Module) SubModule(name string, opts ...ModuleOptionFunc) *Module {
	return &Module{}
}

func (m *Module) Provide(name string, newer interface{}, opts ...ProvideOption) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	def := &BeanDefinition{
		Name:  name,
		New:   newer.(func() interface{}),
		Scope: 0,
	}
	m.definitions[name] = def
	m.dag.Add(name, def)
	return nil
}

func (m *Module) ProvideWithName(name string, newer interface{}, opts ...ProvideOption) error {
	return nil
}

func (m *Module) Invoke() error {
	return nil
}

func (m *Module) MustInvoke() {
}

//func (m *Module) Provide(f interface{}, opts ...ProvideOption) error {
//	m.locker.Lock()
//	defer m.locker.Unlock()
//	hasName := false
//	hasGroup := false
//	for _, opt := range opts {
//		if _, ok := opt.(NameOption); ok {
//			hasName = true
//		}
//		if _, ok := opt.(GroupOption); ok {
//			hasGroup = true
//		}
//	}
//
//	if reflectx.IsBasicType(f) && !hasName && !hasGroup {
//		return BasicTypeShouldWithNameError
//	}
//
//	if hasName && hasGroup {
//		return NameOrGroupOnlyOneError
//	}
//
//	options := BuildDigProvideOption(opts...)
//
//	constructor, err := m.GenerateConstructor(f)
//	if err != nil {
//		return err
//	}
//	return m.digContainer.Provide(constructor, options...)
//}
//
//// GenerateConstructor generates a constructor function for the given instance using reflection
//func (m *Module) GenerateConstructor(instance interface{}) (interface{}, error) {
//	if instance == nil {
//		return nil, NilError
//	}
//
//	t := reflect.TypeOf(instance)
//
//	// Handle pointer types
//	isPtr := t.Kind() == reflect.Ptr
//	if isPtr {
//		t = t.Elem()
//		if t == nil {
//			return nil, NilPtrError
//		}
//	}
//
//	valueOf := reflect.ValueOf(instance)
//	if isPtr {
//		valueOf = valueOf.Elem()
//	}
//
//	switch t.Kind() {
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
//		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
//		reflect.Complex64, reflect.Complex128,
//		reflect.Float32, reflect.Float64,
//		reflect.String, reflect.Bool,
//		reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
//		ctorType := reflect.FuncOf(nil, []reflect.Type{t}, false)
//		fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
//			if isPtr {
//				return []reflect.Value{valueOf.Addr()}
//			}
//			return []reflect.Value{valueOf}
//		})
//		return fn.Interface(), nil
//	case reflect.Struct:
//	default:
//		return nil, TypeError
//	}
//
//	// Collect struct field types as constructor parameters
//	var paramTypes []reflect.Type
//	var fieldIndices []int
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		if field.PkgPath != "" {
//			continue
//		}
//		if tagVal := field.Tag.Get(m.opts.injectTag); tagVal == SkipInjectTag {
//			continue
//		}
//		if !isBasicType(field.Type.Kind()) {
//			paramTypes = append(paramTypes, field.Type)
//			fieldIndices = append(fieldIndices, i)
//		}
//	}
//
//	// Set return type
//	returnType := t
//	if isPtr {
//		returnType = reflect.PointerTo(t)
//	}
//
//	ctorType := reflect.FuncOf(paramTypes, []reflect.Type{returnType}, false)
//
//	fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
//		//newInstance := reflect.New(t).Elem()
//		//
//		//// Copy all fields from initial value
//		//if valueOf.IsValid() {
//		//	newInstance.Set(valueOf)
//		//}
//
//		// Set only injectable fields
//		for i, fieldIdx := range fieldIndices {
//			if i >= len(args) {
//				break
//			}
//			//newInstance.Field(fieldIdx).Set(args[i])
//			valueOf.Field(fieldIdx).Set(args[i])
//		}
//
//		if isPtr {
//			return []reflect.Value{valueOf.Addr()}
//		}
//		return []reflect.Value{valueOf}
//	})
//
//	return fn.Interface(), nil
//}
//
//func isBasicType(k reflect.Kind) bool {
//	switch k {
//	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
//		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
//		reflect.Float32, reflect.Float64, reflect.String:
//		return true
//	default:
//		return false
//	}
//}
//
//func (m *Module) Invoke(f interface{}) error {
//	return m.digContainer.Invoke(f)
//}
//
//func (m *Module) InvokeAll(f interface{}) error {
//	return nil
//}
//
//func (m *Module) Inject(f interface{}) error {
//	typeOf := reflect.TypeOf(f)
//
//	if typeOf.Kind() != reflect.Struct {
//		return fmt.Errorf("type error")
//	}
//
//	for i := 0; i < typeOf.NumField(); i++ {
//
//	}
//	return nil
//}
//
//func (m *Module) LookupByKey(key string) (interface{}, bool) {
//	return nil, false
//}
//
//func (m *Module) LookupByKeyGlobal(key string) (interface{}, bool) {
//	return nil, false
//}
//
//func BuildDigProvideOption(opts ...ProvideOption) []dig.ProvideOption {
//	var provideOpts ProvideOptions
//	for _, opt := range opts {
//		opt.ApplyProvideOption(&provideOpts)
//	}
//	ret := make([]dig.ProvideOption, 0)
//	if provideOpts.Name != "" {
//		ret = append(ret, dig.Name(provideOpts.Name))
//	}
//	if provideOpts.Group != "" {
//		ret = append(ret, dig.Group(provideOpts.Group))
//	}
//	if provideOpts.As != nil {
//		for _, as := range provideOpts.As {
//			ret = append(ret, dig.As(as))
//		}
//	}
//
//	return ret
//}
