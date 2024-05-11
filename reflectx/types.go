package reflectx

import (
	"fmt"
	"reflect"
)



type Reflectx struct {
	// store the origin value pass to reflectx
	originValue interface{}
	// fill value by tag or use golang default value
	fillValue interface{}
	// store struct data field, the key is field name, the value use fieldx to wrap the field value
	FieldMap map[string]*Fieldx
	// to store the struct method or single func/method info
	MethodMap map[string]*Methodx
	// store struct pkg path only struct used
	PkgPath string
	Kind    reflect.Kind
}

type Fieldx struct {
	Name        string
	Value       interface{}
	Tag         reflect.StructTag
	Kind        reflect.Kind
	ReflectType reflect.Type
	PkgPath     string
	Anonymous   bool
}

func (f *Fieldx) String() string {
	return fmt.Sprintf("%+v", *f)
}

type Methodx struct {
	Name       string
	PkgPath    string
	Index      int
	IsExported bool
	InParams   []*Fieldx
	OutParams  []*Fieldx
}

func (f *Methodx) String() string {
	return fmt.Sprintf("%+v", *f)
}
