package main

type Plan interface {
	GetArea() Rect
	GetBlackHole() Rect
	GetBoundaries() []Rect
}

type Entity interface {
	GetCoords() Rect
	Update()
	Render()
	Unload()
}

type Scene interface {
	Update()
	Render()
	Unload()
	Plan
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
