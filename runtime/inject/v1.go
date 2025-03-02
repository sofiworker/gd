package inject

import (
	"fmt"
	"go.uber.org/dig"
	"reflect"
)

const (
	injectTag    = "inject"
	singletonTag = "singleton"
	canNilTag    = "canNil"
	nilAbleTag   = "nilAble"
)

type StartAble interface {
	Start() error
}

type CloseAble interface {
	Close()
}

type Injectable interface {
	StartAble
	CloseAble
}

type Graph struct {
	Container *Container
	itemCount int
}

type Object struct {
	Name        string
	reflectType reflect.Type
	Value       interface{} `gd:"-"`
	closed      bool
}

type DigObjectWrap struct {
	dig.In
	Objects []*Object `group:"objects"`
}

func NewGraph() *Graph {
	container := New().SetInjectTag(injectTag)
	return &Graph{
		Container: container,
	}
}

func getTypeName(t reflect.Type) string {
	isPtr := false
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}
	pkg := t.PkgPath()
	name := t.Name()
	if pkg != "" {
		name = pkg + "." + t.Name()
	}
	if isPtr {
		name = "*" + name
	}
	return name
}

func (g *Graph) FindByType(t reflect.Type) (*Object, bool) {
	return g.findByType(t)
}

func (g *Graph) findByType(t reflect.Type) (*Object, bool) {
	n := getTypeName(t)
	return g.find(n)
}

func (g *Graph) Len() int {
	return g.itemCount
}

func (g *Graph) Find(name string) (*Object, bool) {
	return g.find(name)
}

func (g *Graph) find(name string) (*Object, bool) {
	var obj *Object
	err := g.Container.Invoke(func(wrap *DigObjectWrap) {
		for _, object := range wrap.Objects {
			if object.Name == name {
				obj = object
				break
			}
		}
	})
	if err != nil {
		return nil, false
	}
	if obj == nil {
		return nil, false
	}
	return obj, true
}

func (g *Graph) del(name string) {

}

func (g *Graph) set(name string, o *Object) {

}

func (g *Graph) setBoth(name string, o *Object) {

}

func (g *Graph) RegWithoutInjection(name string, value interface{}) interface{} {
	return g.RegisterOrFailNoFill(name, value)
}

func (g *Graph) RegisterOrFailNoFill(name string, value interface{}) interface{} {
	return nil

}

func (g *Graph) RegisterOrFail(name string, value interface{}) interface{} {
	v, err := g.Register(name, value)
	if err != nil {
		//if g.Logger != nil {
		//	g.Logger.Error(err)
		//}
		panic(fmt.Sprintf("reg fail,name=%v,err=%v", name, err.Error()))
	}
	return v
}

func (g *Graph) Register(name string, value interface{}) (interface{}, error) {
	return g.register(name, value, false, false)
}

func (g *Graph) register(name string, value interface{}, singleton bool, noFill bool) (interface{}, error) {

	t := reflect.TypeOf(value)

	if isStructPtr(t) {
		if name == "" {
			name = getTypeName(t)
		}
	} else {
		if name == "" {
			return nil, fmt.Errorf("name can not be empty,name=%s,type=%v", name, t)
		}
	}

	o := &Object{
		Name:        name,
		reflectType: t,
	}

	//already registered
	found, ok := g.find(name)
	if ok {
		return nil, fmt.Errorf("already registered,name=%s,type=%v,found=%v", name, t, found)
	}

	valueOf := reflect.ValueOf(value)
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
		valueOf = valueOf.Elem()
	}

	if isStructPtr(o.reflectType) {
		//t := reflectType.Elem()
		//var v reflect.Value
		//created := false
		//if isNil(value) {
		//	created = true
		//	v = reflect.New(t)
		//} else {
		//	v = reflect.ValueOf(value)
		//}
		//
		//for i := 0; i < t.NumField(); i++ {
		//	if !created && noFill {
		//		continue
		//	}
		//
		//	f := t.Field(i)
		//	vfe := v.Elem()
		//	vf := vfe.Field(i)
		//
		//	tag, ok := f.Tag.Lookup(injectTag)
		//
		//	if !ok {
		//		continue
		//	}
		//
		//	if vf.CanInterface() {
		//		if reflect.ValueOf(vf.Interface()).Kind() == reflect.Struct {
		//			return nil, fmt.Errorf("inject a struct field is not supported,field=%v,type=%v", f.Name, t.Name())
		//		}
		//
		//		if !isZeroOfUnderlyingType(vf.Interface()) {
		//			continue
		//		}
		//	}
		//
		//	if f.Anonymous || !vf.CanSet() {
		//		return nil, fmt.Errorf("inject tag must on a public field!field=%s,type=%s", f.Name, t.Name())
		//	}
		//
		//	tag, _ = f.Tag.Lookup(singletonTag)
		//	singletonTag := false
		//	if tag == "true" {
		//		singletonTag = true
		//	}
		//	canNilStr, _ := f.Tag.Lookup(canNilTag)
		//	nilAbleStr, _ := f.Tag.Lookup(nilAbleTag)
		//	canNil := false
		//	if canNilStr == "true" || nilAbleStr == "true" {
		//		canNil = true
		//	}
		//
		//	var found *Object
		//	if tag != "" {
		//		//due to default singleton of struct ptr injections
		//		//we should first find by name,then find by type
		//		found, ok = g.find(tag)
		//		if singletonTag && !ok && isStructPtr(f.Type) {
		//			found, ok = g.findByType(f.Type)
		//		}
		//	} else {
		//		found, ok = g.findByType(f.Type)
		//	}
		//
		//	if !ok || found == nil {
		//		if canNil {
		//			continue
		//		}
		//		if isStructPtr(f.Type) {
		//			_, err := g.register(tag, reflect.NewAt(f.Type.Elem(), nil).Interface(), singletonTag, noFill)
		//			if err != nil {
		//				return nil, err
		//			}
		//		} else {
		//			var implFound reflect.Type
		//			//impls := Get(tag)
		//			//for _, impl := range impls {
		//			//	if impl == nil {
		//			//		continue
		//			//	}
		//			//	if impl.AssignableTo(f.Type) {
		//			//		implFound = impl
		//			//		break
		//			//	}
		//			//
		//			//}
		//
		//			if implFound != nil {
		//				_, err := g.register(tag, reflect.NewAt(implFound.Elem(), nil).Interface(), singletonTag, noFill)
		//				if err != nil {
		//					return nil, err
		//				}
		//			} else {
		//				return nil, fmt.Errorf("dependency field=%s,tag=%s not found in object %s:%v", f.Name, tag, name, reflectType)
		//			}
		//		}
		//
		//		if tag != "" {
		//			found, ok = g.find(tag)
		//			if !ok && singleton {
		//				found, ok = g.findByType(f.Type)
		//			}
		//		} else {
		//			found, ok = g.findByType(f.Type)
		//		}
		//	}
		//
		//	if !ok || found == nil {
		//		return nil, fmt.Errorf("dependency %s not found in object %s:%v", f.Name, name, reflectType)
		//	}
		//
		//	reflectFoundValue := reflect.ValueOf(found.Value)
		//	if !found.reflectType.AssignableTo(f.Type) {
		//		switch reflectFoundValue.Kind() {
		//		case reflect.Int:
		//			fallthrough
		//		case reflect.Int8:
		//			fallthrough
		//		case reflect.Int16:
		//			fallthrough
		//		case reflect.Int32:
		//			fallthrough
		//		case reflect.Int64:
		//			iv := reflectFoundValue.Int()
		//			switch f.Type.Kind() {
		//			case reflect.Int:
		//				fallthrough
		//			case reflect.Int8:
		//				fallthrough
		//			case reflect.Int16:
		//				fallthrough
		//			case reflect.Int32:
		//				fallthrough
		//			case reflect.Int64:
		//				vf.SetInt(iv)
		//			default:
		//				return nil, fmt.Errorf("dependency name=%s,type=%v not valid in object %s:%v", f.Name, f.Type, name, reflectType)
		//			}
		//		case reflect.Float32:
		//			fallthrough
		//		case reflect.Float64:
		//			fv := reflectFoundValue.Float()
		//			switch f.Type.Kind() {
		//			case reflect.Float32:
		//				fallthrough
		//			case reflect.Float64:
		//				vf.SetFloat(fv)
		//			default:
		//				return nil, fmt.Errorf("dependency name=%s,type=%v not valid in object %s:%v", f.Name, f.Type, name, reflectType)
		//			}
		//		default:
		//			return nil, fmt.Errorf("dependency name=%s,type=%v not valid in object %s:%v", f.Name, f.Type, name, reflectType)
		//		}
		//	} else {
		//		vf.Set(reflectFoundValue)
		//	}
		//}
		//o.Value = v.Interface()
	} else {
		if canNil(value) && isNil(value) {
			return nil, fmt.Errorf("register nil on name=%s, val=%v", name, value)
		}
		o.Value = value
	}

	err := g.Container.Provide(o, WithGroup("objects"))
	if err != nil {
		return nil, fmt.Errorf("provide object fail,name=%v,err=%v", name, err)
	}

	// dependency resolved, init the object
	canStart, ok := o.Value.(StartAble)
	if ok {
		err := canStart.Start()
		if err != nil {
			return nil, fmt.Errorf("start object fail,name=%v,err=%v", name, err)
		}
	}

	//set to graph
	//if isStructPtr(reflectType) && singleton {
	//	g.setBoth(name, o)
	//} else {
	//	g.set(name, o)
	//}
	return o.Value, nil
}

func (g *Graph) RegisterOrFailSingleNoFill(name string, value interface{}) interface{} {
	v, err := g.RegisterSingleNoFill(name, value)
	if err != nil {
		//if g.Logger != nil {
		//	g.Logger.Error(err)
		//}
		panic(fmt.Sprintf("reg fail,name=%v,err=%v", name, err.Error()))
	}
	return v
}

func (g *Graph) RegisterSingleNoFill(name string, value interface{}) (interface{}, error) {
	return g.register(name, value, true, true)
}

func (g *Graph) Close() {
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func canNil(v interface{}) bool {
	k := reflect.ValueOf(v).Kind()
	return k == reflect.Ptr || k == reflect.Interface
}

func isNil(v interface{}) bool {
	return reflect.ValueOf(v).IsNil()
}

func isZeroOfUnderlyingType(x interface{}) bool {
	if x == nil {
		return true
	}
	rv := reflect.ValueOf(x)
	k := rv.Kind()

	if k == reflect.Func {
		return rv.IsNil()
	}

	if (k == reflect.Ptr || k == reflect.Interface || k == reflect.Chan || k == reflect.Map || k == reflect.Slice) && rv.IsNil() {
		return true
	}

	switch k {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if rv.Len() <= 0 {
			return true
		} else {
			return false
		}
	}
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
