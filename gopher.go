package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var b = Rect{
	Left:  Coord{X: -16, Y: 10},
	Right: Coord{X: 16, Y: 10},
}

var (
	isLeft, isRight, isTop, isBottom bool
)

type Gopher struct {
	window     *glfw.Window
	h, w, x, y float32
}

func newGopher(window *glfw.Window) *Gopher {
	return &Gopher{
		window: window,
		h:      0.8,
		w:      0.8,
		x:      0.0, // TODO
		y:      0.0, // TODO
	}
}

func (g *Gopher) Update() {
	isButtonPress := func(b glfw.Key) bool {
		return g.window.GetKey(b) == glfw.Press
	}
	isLeft = isButtonPress(glfw.KeyLeft)
	isRight = isButtonPress(glfw.KeyRight)
	isTop = isButtonPress(glfw.KeyUp)
	isBottom = isButtonPress(glfw.KeyDown)
}

func (g *Gopher) Render() {
	gl.PushMatrix()
	{
		if isLeft {
			g.x -= 0.3
			if g.x < b.Left.X+g.w {
				g.x = b.Left.X + g.w
			}
		}
		if isRight {
			g.x += 0.3
			if g.x > b.Right.X-g.w {
				g.x = b.Right.X - g.w
			}
		}
		if isTop {
			g.y += 3
		}
		if isBottom {
		}
		gl.Translatef(g.x, g.y, 0)

		rect := Rect{
			Left:  Coord{-g.w, -g.h},
			Right: Coord{g.w, g.h},
		}
		drawTexture(textureGopher, rect)
		if isTop {
			g.y -= 3
		}
	}
	gl.PopMatrix()
}
