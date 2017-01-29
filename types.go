package main

type Plan interface {
	Boundaries() []Rect
}

type GameElement interface {
	Update()
	Render()
	Unload()
}

type Scene interface {
	GameElement
	Plan
}

type Coord struct {
	X, Y float32
}

type Rect struct {
	Left, Right Coord
}
