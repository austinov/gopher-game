package main

type GameElement interface {
	Update()
	Render()
	Unload()
}

type Scene interface {
	Boundaries() []Rect
	GameElement
}

type Coord struct {
	X, Y float32
}

type Rect struct {
	Left, Right Coord
}
