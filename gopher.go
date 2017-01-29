package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// TODO refactoring
var b = Rect{
	Left:  Coord{X: -16, Y: 10},
	Right: Coord{X: 16, Y: 10},
}

// TODO refactoring
var (
	isLeft, isRight, isTop, isBottom bool
)

type Gopher struct {
	window  *glfw.Window
	texture uint32
	plan    Plan
	curr    Rect
	r2l     bool
}

func NewGopher(window *glfw.Window, plan Plan) GameElement {
	return &Gopher{
		window:  window,
		texture: NewTexture("assets/gopher.png"),
		plan:    plan,
		curr: Rect{
			Left:  Coord{-0.8, -0.8},
			Right: Coord{0.8, 0.8},
		},
		//x: 0.0, // TODO center
		//y: -3, // TODO center
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
		r2l := g.r2l
		var dx, dy float32
		if isLeft {
			dx -= 0.1
			r2l = true
		}
		if isRight {
			dx += 0.1
			r2l = false
		}
		if isTop {
			dy += 0.1
		}
		if isBottom {
			dy -= 0.1
		}
		g.curr.Left.X += dx
		g.curr.Right.X += dx
		g.curr.Left.Y += dy
		g.curr.Right.Y += dy
		for _, bound := range g.plan.Boundaries() {
			if ((g.curr.Left.X > bound.Left.X && g.curr.Left.X < bound.Right.X) ||
				(g.curr.Right.X > bound.Left.X && g.curr.Right.X < bound.Right.X) ||
				(g.curr.Left.X < bound.Left.X && g.curr.Right.X > bound.Right.X)) &&
				((g.curr.Left.Y > bound.Right.Y && g.curr.Left.Y < bound.Left.Y) ||
					(g.curr.Right.Y > bound.Right.Y && g.curr.Right.Y < bound.Left.Y) ||
					(g.curr.Left.Y > bound.Left.Y && g.curr.Right.Y < bound.Right.Y)) {
				if isLeft {
					//g.x = bound.Right.X + 0.8
					g.curr.Left.X = bound.Right.X
					g.curr.Right.X = g.curr.Left.X + 0.8 // TODO extract 0.8

				}
				if isRight {
					//g.x = bound.Left.X - 0.8
					g.curr.Right.X = bound.Left.X
					g.curr.Left.X = g.curr.Right.X - 0.8
				}
				if isTop {
					//g.y = bound.Right.Y - 0.8
					g.curr.Left.Y = bound.Right.Y
					g.curr.Right.Y = g.curr.Left.Y - 0.8
				}
				if isBottom {
					//g.y = bound.Left.Y + 0.8
					g.curr.Right.Y = bound.Left.Y
					g.curr.Left.Y = g.curr.Right.Y - 0.8
				}
				// boundary violation
			}
		}

		//gl.Translatef(g.x, g.y, 0)
		gl.Translatef(g.curr.Left.X+0.8, g.curr.Right.Y+0.8, 0)

		rect := Rect{
			Left:  Coord{-0.8, -0.8},
			Right: Coord{0.8, 0.8},
		}
		// set previouse direction
		if g.r2l {
			rect.Left.X *= float32(-1.0)
			rect.Right.X *= float32(-1.0)
		}
		if r2l != g.r2l {
			// set new direction
			g.r2l = r2l
			rect.Left.X *= float32(-1.0)
			rect.Right.X *= float32(-1.0)
		}
		DrawTexture(g.texture, rect)
		if isTop {
			//g.y -= 3
		}
	}
	gl.PopMatrix()
}

func (g *Gopher) Unload() {
	gl.DeleteTextures(1, &g.texture)
}
