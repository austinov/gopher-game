package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Gopher struct {
	window  *glfw.Window
	plan    Plan
	texture uint32
	coords  Rect
	height  float32
	width   float32
	r2l     bool // moving direction: true when right to left, false when left to right
	onFloor bool
	floor   Rect
}

func NewGopher(window *glfw.Window, plan Plan) Entity {
	texture, bounds := NewTexture("assets/gopher.png")
	h, w := bounds.Right.Y, bounds.Right.X
	return &Gopher{
		window:  window,
		plan:    plan,
		texture: texture,
		height:  h,
		width:   w,
		coords: Rect{
			Left:  Point{-w, h},
			Right: Point{w, -h},
		},
	}
}

func (g *Gopher) Update() {
	isButtonPress := func(b glfw.Key) bool {
		return g.window.GetKey(b) == glfw.Press
	}
	isLeft := isButtonPress(glfw.KeyLeft)
	isRight := isButtonPress(glfw.KeyRight)
	isTop := isButtonPress(glfw.KeyUp)

	checkBoundaries := func() (bool, Rect) {
		return CheckBoundaries(g.coords, g.plan.GetBoundaries()...)
	}
	if isLeft {
		g.coords.Left.X -= 0.1
		if in, bound := checkBoundaries(); in {
			g.coords.Left.X = bound.Right.X + 0.001
		}
		g.coords.Right.X = g.coords.Left.X + 2*g.width
		g.r2l = true
	}
	if isRight {
		g.coords.Right.X += 0.1
		if in, bound := checkBoundaries(); in {
			g.coords.Right.X = bound.Left.X - 0.001
		}
		g.coords.Left.X = g.coords.Right.X - 2*g.width
		g.r2l = false
	}
	if g.onFloor {
		if (g.coords.Left.X < g.floor.Left.X && g.coords.Right.X < g.floor.Left.X) ||
			(g.coords.Right.X > g.floor.Right.X && g.coords.Left.X > g.floor.Right.X) {
			g.onFloor = false
			g.floor = Rect{}
		}
	}
	if isTop && g.onFloor {
		g.coords.Left.Y += 6.6
		if in, bound := checkBoundaries(); in {
			g.coords.Left.Y = bound.Right.Y - 0.001
		}
		g.coords.Right.Y = g.coords.Left.Y - 2*g.height
		g.onFloor = false
	}
	if !g.onFloor {
		g.coords.Right.Y -= 0.05
		if in, bound := checkBoundaries(); in {
			g.coords.Right.Y = bound.Left.Y + 0.001
			g.onFloor = true
			g.floor = bound
		}
		g.coords.Left.Y = g.coords.Right.Y + 2*g.height
	}
}

func (g *Gopher) GetCoords() Rect {
	return g.coords
}

func (g *Gopher) Render() {
	// TODO extract
	gl.PushMatrix()
	{
		gl.Translatef(g.coords.Left.X+g.width, g.coords.Right.Y+g.height, 0)

		var rect Rect
		if g.r2l {
			rect = Rect{
				Left:  Point{g.width, -g.height},
				Right: Point{-g.width, g.height},
			}
		} else {
			rect = Rect{
				Left:  Point{-g.width, -g.height},
				Right: Point{g.width, g.height},
			}
		}
		DrawTexture(g.texture, rect)
	}
	gl.PopMatrix()
}

func (g *Gopher) Unload() {
	gl.DeleteTextures(1, &g.texture)
}
