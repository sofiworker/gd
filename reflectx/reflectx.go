package reflectx

import (
	"fmt"
	"reflect"

	"github.com/chuck1024/gd/v2/gerr"
)

var (
	NilValue    = fmt.Errorf("nil value")
	UnknownType = fmt.Errorf("unknown type")
)

// New return a reflectx instance
func New(v interface{}) *Reflectx {
	if IsNil(v) {
		panic(NilValue)
	}

	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)

	instance := &Reflectx{
		originValue: v,
		FieldMap:    make(map[string]*Fieldx),
		MethodMap:   make(map[string]*Methodx),
		PkgPath:     typeOf.PkgPath(),
		Kind:        typeOf.Kind(),
	}
PTR:
	switch typeOf.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.String:
	case reflect.Bool:
	case reflect.Array:
	case reflect.Slice:
	case reflect.Map:
	case reflect.Chan:
	case reflect.Ptr:
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
		goto PTR
	case reflect.Struct:
		for i := 0; i < typeOf.NumField(); i++ {
			structField := typeOf.Field(i)
			field := &Fieldx{
				Name:      structField.Name,
				Tag:       structField.Tag,
				Kind:      structField.Type.Kind(),
				PkgPath:   structField.PkgPath,
				Anonymous: structField.Anonymous,
			}
			instance.FieldMap[field.Name] = field
		}
		for i := 0; i < typeOf.NumMethod(); i++ {
			m := typeOf.Method(i)
			methodx := &Methodx{
				Name:       m.Name,
				PkgPath:    m.PkgPath,
				Index:      m.Index,
				IsExported: m.IsExported(),
				Kind:       PtrMethod,
			}
			methodType := m.Type
			methodx.InParams, methodx.OutParams = GetFuncOrMethodInParams(methodType), GetFuncOrMethodOutParams(methodType)
			instance.MethodMap[m.Name] = methodx
		}
	case reflect.Func:
		methodx := &Methodx{
			Name:    typeOf.Name(),
			PkgPath: typeOf.PkgPath(),
		}
		methodx.InParams, methodx.OutParams = GetFuncOrMethodInParams(typeOf), GetFuncOrMethodOutParams(typeOf)
		instance.MethodMap[typeOf.Name()] = methodx
	default:
		panic(UnknownType)
	}
	return instance
}

func CheckPtrDepth(v interface{}) int {
	typeOf := reflect.TypeOf(v)
	ret := 0
	for typeOf.Kind() == reflect.Pointer {
		ret++
		typeOf = typeOf.Elem()
	}
	return ret
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// Default will return `v` with filling default value
func (r *Reflectx) Default() interface{} {
	return r.DefaultByTag("")
}

func (r *Reflectx) DefaultAgain() interface{} {
	return r.DefaultByTagAgain("")
}

func (r *Reflectx) DefaultByTag(tagName string) interface{} {
	if r.fillValue != nil {
		return r.fillValue
	}
	return r.DefaultByTagAgain(tagName)
}

func (r *Reflectx) DefaultByTagAgain(tagName string) interface{} {
	return 1
}

func GetFuncOrMethodInParams(valueOf interface{}) []*Fieldx {
	var (
		funcType reflect.Type
	)
	switch v := valueOf.(type) {
	case reflect.Value:
		funcType = v.Type()
	case reflect.Type:
		funcType = v
	default:
		funcType = reflect.TypeOf(valueOf)
		//if funcType.Kind() == reflect.Pointer {
		//	funcType = funcType.Elem()
		//}
	}

	if funcType.Kind() != reflect.Func {
		panic(gerr.ParamsMustBeFunc)
	}

	ret := make([]*Fieldx, 0)
	for i := 0; i < funcType.NumIn(); i++ {
		f := &Fieldx{
			Name:        funcType.In(i).Name(),
			Tag:         "",
			Kind:        funcType.In(i).Kind(),
			PkgPath:     funcType.In(i).PkgPath(),
			ReflectType: funcType.In(i),
		}
		ret = append(ret, f)
	}
	return ret
}

func GetFuncOrMethodOutParams(valueOf interface{}) []*Fieldx {
	var (
		funcType reflect.Type
	)
	switch v := valueOf.(type) {
	case reflect.Value:
		funcType = v.Type()
	case reflect.Type:
		funcType = v
	default:
		funcType = reflect.TypeOf(valueOf)
		//if funcType.Kind() == reflect.Pointer {
		//	funcType = funcType.Elem()
		//}
	}

	if funcType.Kind() != reflect.Func {
		panic(gerr.ParamsMustBeFunc)
	}

	ret := make([]*Fieldx, 0)
	for i := 0; i < funcType.NumOut(); i++ {
		f := &Fieldx{
			Name:        funcType.Out(i).Name(),
			Tag:         "",
			Kind:        funcType.Out(i).Kind(),
			PkgPath:     funcType.Out(i).PkgPath(),
			ReflectType: funcType.Out(i),
		}
		ret = append(ret, f)
	}
	return ret
}

func GetFuncMap(v, t interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	var (
		valueOf reflect.Value
		typeOf  reflect.Type
	)
	if vv, ok := v.(reflect.Value); ok {
		valueOf = vv
	} else {
		valueOf = reflect.ValueOf(v)
		if valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
	}
	if tt, ok := t.(reflect.Type); ok {
		typeOf = tt
	} else {
		typeOf = reflect.TypeOf(t)
		if typeOf.Kind() == reflect.Pointer {
			typeOf = typeOf.Elem()
		}
	}
	return m
}

func (r *Reflectx) CallMethodByName(name string, args ...interface{}) (interface{}, error) {
	// if m, ok := r.MethodMap[name]; ok {

	// }
	return nil, gerr.NotFoundMethod
}

func (r *Reflectx) CallMethodByIndex(index int, args ...interface{}) {

}

func (r *Reflectx) IsMethodOrFunc() bool {
	return r.Kind == reflect.Func
}

func (r *Reflectx) GetNumIn(index int) int {
	//return r.
	return 0
}

func (r *Reflectx) GetNumOut() int {
	return 0
}

func (r *Reflectx) Call(args ...interface{}) ([]*Fieldx, error) {
	// f := m.(reflect.Value)
	// vals := make([]reflect.Value, len(args))
	// for i := 0; i < len(args); i++ {
	// 	vals[i] = reflect.ValueOf(args[i])
	// }
	// return f.Call(vals), nil
	return nil, nil
}

func (r *Reflectx) IsStructOrStructPtr() bool {
	return r.Kind == reflect.Struct
}

func (r *Reflectx) IsPtr() bool {
	return r.Kind == reflect.Pointer
}

func (r *Reflectx) String() string {
	return fmt.Sprintf("%+v", *r)
}

func (f *Fieldx) IsError() bool {
	return false
}

func IsBasicType(data interface{}) bool {
	t := reflect.TypeOf(data)
	if t == nil {
		return false
	}
	// 解引用指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 检查基础类型
	switch t.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}
