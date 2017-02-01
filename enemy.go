package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Enemy struct {
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

func NewEnemy(window *glfw.Window, plan Plan) GameElement {
	h, w := float32(0.8), float32(0.8)
	return &Enemy{
		window:  window,
		texture: NewTexture("assets/enemy.png"),
		plan:    plan,
		height:  h,
		width:   w,
		curr: Rect{
			Left:  Coord{-w, 8},
			Right: Coord{w, 8 - h},
		},
	}
}

func (e *Enemy) Update() {
	checkBoundaries := func() (bool, Rect) {
		return CheckBoundaries(e.curr, e.plan.Boundaries())
	}
	if e.onFloor {
		if (e.curr.Left.X < e.floor.Left.X && e.curr.Right.X < e.floor.Left.X) ||
			(e.curr.Right.X > e.floor.Right.X && e.curr.Left.X > e.floor.Right.X) {
			e.onFloor = false
			e.floor = Rect{}
		}
	}
	if !e.onFloor {
		e.curr.Right.Y -= 0.03
		if in, bound := checkBoundaries(); in {
			e.curr.Right.Y = bound.Left.Y + 0.03
			e.onFloor = true
			e.floor = bound
		}
		e.curr.Left.Y = e.curr.Right.Y + 2*e.height
	} else {
		if e.r2l {
			e.curr.Left.X -= 0.07
		} else {
			e.curr.Right.X += 0.07 // 0.04 < dx < 0.09
		}
		if in, bound := checkBoundaries(); in {
			if e.r2l {
				e.curr.Left.X = bound.Right.X + 0.001
				e.curr.Right.X = e.curr.Left.X + 2*e.width
				e.r2l = false
			} else {
				e.curr.Right.X = bound.Left.X - 0.001
				e.curr.Left.X = e.curr.Right.X - 2*e.width
				e.r2l = true
			}
		} else {
			if e.r2l {
				e.curr.Right.X = e.curr.Left.X + 2*e.width
			} else {
				e.curr.Left.X = e.curr.Right.X - 2*e.width
			}
		}
	}
}

func (e *Enemy) Render() {
	gl.PushMatrix()
	{
		gl.Translatef(e.curr.Left.X+e.width, e.curr.Right.Y+e.height, 0)

		rect := Rect{
			Left:  Coord{-e.width, -e.height},
			Right: Coord{e.width, e.height},
		}
		DrawTexture(e.texture, rect)
	}
	gl.PopMatrix()
}

func (e *Enemy) Unload() {
	gl.DeleteTextures(1, &e.texture)
}
