package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Plan interface {
	Boundaries() []Rect
}

type Player interface {
	Update(window *glfw.Window, plan Plan)
	Coords() Rect
	Render()
	Unload()
}

type Scene interface {
	Update()
	Render()
	Unload()
	Plan
}

type Point struct {
	X, Y float32
}

type Rect struct {
	Left  Point
	Right Point
}
