package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Plan interface {
	GetArea() Rect
	GetBlackHole() Rect
	GetBoundaries() []Rect
}

type Player interface {
	Update(window *glfw.Window, plan Plan)
	GetCoords() Rect
	Render()
	Unload()
}

type Scene interface {
	Update()
	Render()
	Unload()
	Plan
}

type Gift interface {
	GetCoords() Rect
	Render()
	Unload()
}

type Board interface {
	Show(score, lives int)
	Unload()
}

type Point struct {
	X, Y float32
}

type Rect struct {
	Left  Point
	Right Point
}
