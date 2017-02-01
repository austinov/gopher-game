package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Gopher struct {
	window  *glfw.Window
	texture uint32
	plan    Plan
	curr    Rect
	height  float32
	width   float32
	r2l     bool
	onFloor bool
	floor   Rect
}

func NewGopher(window *glfw.Window, plan Plan) GameElement {
	h, w := float32(0.8), float32(0.8)
	return &Gopher{
		window:  window,
		texture: NewTexture("assets/gopher.png"),
		plan:    plan,
		height:  h,
		width:   w,
		curr: Rect{
			Left:  Coord{-w, h},
			Right: Coord{w, -h},
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
		return CheckBoundaries(g.curr, g.plan.Boundaries())
	}
	if isLeft {
		g.curr.Left.X -= 0.1
		if in, bound := checkBoundaries(); in {
			g.curr.Left.X = bound.Right.X + 0.001
		}
		g.curr.Right.X = g.curr.Left.X + 2*g.width
		g.r2l = true
	}
	if isRight {
		g.curr.Right.X += 0.1
		if in, bound := checkBoundaries(); in {
			g.curr.Right.X = bound.Left.X - 0.001
		}
		g.curr.Left.X = g.curr.Right.X - 2*g.width
		g.r2l = false
	}
	if g.onFloor {
		if (g.curr.Left.X < g.floor.Left.X && g.curr.Right.X < g.floor.Left.X) ||
			(g.curr.Right.X > g.floor.Right.X && g.curr.Left.X > g.floor.Right.X) {
			g.onFloor = false
			g.floor = Rect{}
		}
	}
	if isTop && g.onFloor {
		g.curr.Left.Y += 6.6
		if in, bound := checkBoundaries(); in {
			g.curr.Left.Y = bound.Right.Y - 0.001
		}
		g.curr.Right.Y = g.curr.Left.Y - 2*g.height
		g.onFloor = false
	}
	if !g.onFloor {
		g.curr.Right.Y -= 0.05
		if in, bound := checkBoundaries(); in {
			g.curr.Right.Y = bound.Left.Y + 0.001
			g.onFloor = true
			g.floor = bound
		}
		g.curr.Left.Y = g.curr.Right.Y + 2*g.height
	}
}

func (g *Gopher) Render() {
	gl.PushMatrix()
	{
		gl.Translatef(g.curr.Left.X+g.width, g.curr.Right.Y+g.height, 0)

		var rect Rect
		if g.r2l {
			rect = Rect{
				Left:  Coord{g.width, -g.height},
				Right: Coord{-g.width, g.height},
			}
		} else {
			rect = Rect{
				Left:  Coord{-g.width, -g.height},
				Right: Coord{g.width, g.height},
			}
		}
		DrawTexture(g.texture, rect)
	}
	gl.PopMatrix()
}

func (g *Gopher) Unload() {
	gl.DeleteTextures(1, &g.texture)
}
