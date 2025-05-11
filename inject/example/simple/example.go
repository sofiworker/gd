package main

import "github.com/chuck1024/gd/v2/inject"

type Wheel struct {
}

type Engine struct {
}

type Car struct {
	Engine *Engine
	Wheels []*Wheel
}

func (c *Car) Start() {
	println("vroooom")
}

func main() {
	container := inject.New()
	root := container.Root()
	err := root.Provide("car", func() *Car {
		return nil
	})
	if err != nil {
		panic(err)
	}
}
