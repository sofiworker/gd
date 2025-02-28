package inject

import (
	"github.com/chuck1024/gd/v2/logger"
	"reflect"
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

func InitDefault() {
	//g = NewGraph()
}

func Close() {
	//g.Close()
}

func SetLogger(logger logger.Logger) {
	//g.Logger = logger
}

func RegisterOrFailNoFill(name string, value interface{}) interface{} {
	//return g.RegisterOrFailNoFill(name, value)
	return nil
}

func RegWithoutInjection(name string, value interface{}) interface{} {
	//return g.RegWithoutInjection(name, value)
	return nil
}

func Reg(name string, value interface{}) interface{} {
	return RegisterOrFail(name, value)
}

func RegisterOrFail(name string, value interface{}) interface{} {
	//return g.RegisterOrFail(name, value)
	return nil
}

func Register(name string, value interface{}) (interface{}, error) {
	//return g.Register(name, value)
	return nil, nil
}

func RegisterOrFailSingleNoFill(name string, value interface{}) interface{} {
	//return g.RegisterOrFailSingleNoFill(name, value)
	return nil
}

func RegisterOrFailSingle(name string, value interface{}) interface{} {
	//return g.RegisterOrFailSingle(name, value)
	return nil
}

func RegisterSingle(name string, value interface{}) (interface{}, error) {
	//return g.RegisterSingle(name, value)
	return nil, nil
}

func FindByType(t reflect.Type) (interface{}, bool) {
	//o, ok := g.FindByType(t)
	//if !ok || o == nil || o.Value == nil {
	//	return nil, false
	//}
	//return o.Value, ok
	return nil, false
}

func Find(name string) (interface{}, bool) {
	//o, ok := g.Find(name)
	//if !ok || o == nil || o.Value == nil {
	//	return nil, false
	//}
	//return o.Value, ok
	return nil, false
}

func GraphLen() int {
	return 0
}

func GraphPrint() string {
	return ""
}
