/**
 * Copyright 2020 gd Author. All rights reserved.
 * Author: Chuck1024
 */

package inject

import (
	"reflect"
)

var g *Graph

func InitDefault() {
	g = NewGraph()
}

//func Close() {
//	g.Close()
//}
//
//func SetLogger(logger Logger) {
//	g.Logger = logger
//}

func RegisterOrFailNoFill(name string, value interface{}) interface{} {
	return g.RegisterOrFailNoFill(name, value)
}

func RegWithoutInjection(name string, value interface{}) interface{} {
	return g.RegWithoutInjection(name, value)
}

func Reg(name string, value interface{}) interface{} {
	return RegisterOrFail(name, value)
}

func RegisterOrFail(name string, value interface{}) interface{} {
	return g.RegisterOrFail(name, value)
}

func Register(name string, value interface{}) (interface{}, error) {
	return g.Register(name, value)
}

func RegisterOrFailSingleNoFill(name string, value interface{}) interface{} {
	return g.RegisterOrFailSingleNoFill(name, value)
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
	o, ok := g.FindByType(t)
	if !ok || o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, ok
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
	//return g.Len()
	return 1
}

func GraphPrint() string {
	//return g.SPrint()
	return ""
}
