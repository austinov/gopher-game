package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Enemy struct {
	window  *glfw.Window
	texture uint32
	h, w    float32
	x, y    float32
}

func NewEnemy(window *glfw.Window) GameElement {
	return &Enemy{
		window:  window,
		texture: newTexture("assets/enemy.png"),
		h:       0.8,
		w:       0.8,
		x:       0.0, // TODO
		y:       0.0, // TODO
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
		drawTexture(e.texture, rect)
	}
	gl.PopMatrix()
}

func (e *Enemy) Unload() {
	gl.DeleteTextures(1, &e.texture)
}
