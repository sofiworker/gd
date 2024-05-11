package reflectx

import (
	"fmt"
	"reflect"

	"github.com/chuck1024/gd/v2/gerr"
)

// New only can decode ptr, if use ptr to ptr will panic
func New(v interface{}) *Reflectx {
	if CheckPtrDepth(v) > 1 {
		panic(gerr.NotAllowMultiLayerPointer)
	}
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	var fillValue interface{}
	fieldMap := make(map[string]*Fieldx)
	methodMap := make(map[string]*Methodx)
	switch typeOf.Kind() {
	case reflect.Bool:
		b := reflect.New(typeOf).Bool()
		fillValue = b
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
		fieldMap[typeOf.Name()] = &Fieldx{}
	case reflect.Array, reflect.Slice:
	case reflect.Chan:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.String:
	case reflect.Struct:
		// handle struct field
		nums := typeOf.NumField()
		for i := 0; i < nums; i++ {
			f := typeOf.Field(i)
			v := valueOf.Field(i)
			field := &Fieldx{
				Name:      f.Name,
				Tag:       f.Tag,
				Kind:      f.Type.Kind(),
				PkgPath:   f.PkgPath,
				Anonymous: f.Anonymous,
			}
			if v.CanInterface() {
				field.Value = v.Interface()
			}
			fieldMap[f.Name] = field
		}
		// handle struct method
		typeOf = reflect.PointerTo(typeOf)
		nums = typeOf.NumMethod()
		for i := 0; i < nums; i++ {
			m := typeOf.Method(i)
			methodx := &Methodx{
				Name:       m.Name,
				PkgPath:    m.PkgPath,
				Index:      m.Index,
				IsExported: m.IsExported(),
			}
			methodType := m.Type
			methodx.InParams, methodx.OutParams = GetFuncOrMethodInParams(methodType), GetFuncOrMethodOutParams(methodType)
			methodMap[m.Name] = methodx
		}

	case reflect.Func:
		methodx := &Methodx{
			Name:    typeOf.Name(),
			PkgPath: typeOf.PkgPath(),
		}
		methodx.InParams, methodx.OutParams = GetFuncOrMethodInParams(typeOf), GetFuncOrMethodOutParams(typeOf)
		methodMap[typeOf.Name()] = methodx
	}
	return &Reflectx{
		originValue: v,
		fillValue:   fillValue,
		FieldMap:    fieldMap,
		MethodMap:   methodMap,
		PkgPath:     typeOf.PkgPath(),
		Kind:        typeOf.Kind(),
	}
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
		if funcType.Kind() == reflect.Pointer {
			funcType = funcType.Elem()
		}
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
		if funcType.Kind() == reflect.Pointer {
			funcType = funcType.Elem()
		}
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
