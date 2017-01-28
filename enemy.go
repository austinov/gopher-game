package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Enemy struct {
	window     *glfw.Window
	h, w, x, y float32
}

func newEnemy(window *glfw.Window) *Enemy {
	return &Enemy{
		window: window,
		h:      0.8,
		w:      0.8,
		x:      0.0, // TODO
		y:      0.0, // TODO
	}
}

func (e *Enemy) Update() {
}

func (e *Enemy) Render() {
	gl.PushMatrix()
	{
		gl.Translatef(e.x, e.y, 0)

		rect := Rect{
			Left:  Coord{-e.w, -e.h},
			Right: Coord{e.w, e.h},
		}
		drawTexture(textureEnemy, rect)
	}
	gl.PopMatrix()
}
