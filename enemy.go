package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Enemy struct {
	window  *glfw.Window
	plan    Plan
	texture uint32
	coords  Rect
	height  float32
	width   float32
	r2l     bool
	onFloor bool
	floor   Rect
}

func NewEnemy(window *glfw.Window, plan Plan, r2l bool) Entity {
	texture, bounds := NewTexture("assets/enemy.png")
	h, w := bounds.Right.Y, bounds.Right.X
	return &Enemy{
		window:  window,
		plan:    plan,
		texture: texture,
		coords: Rect{
			Left:  Point{-w, 8},
			Right: Point{w, 8 - h},
		},
		height: h,
		width:  w,
		r2l:    r2l,
	}
}

func (e *Enemy) Update() {
	checkBoundaries := func() (bool, Rect) {
		return CheckBoundaries(e.coords, e.plan.GetBoundaries()...)
	}
	if e.onFloor {
		if (e.coords.Left.X < e.floor.Left.X && e.coords.Right.X < e.floor.Left.X) ||
			(e.coords.Right.X > e.floor.Right.X && e.coords.Left.X > e.floor.Right.X) {
			e.onFloor = false
			e.floor = Rect{}
		}
	}
	if !e.onFloor {
		e.coords.Right.Y -= 0.04
		if in, bound := checkBoundaries(); in {
			e.coords.Right.Y = bound.Left.Y + 0.04
			e.onFloor = true
			e.floor = bound
		}
		e.coords.Left.Y = e.coords.Right.Y + 2*e.height
	} else {
		if e.r2l {
			e.coords.Left.X -= 0.07
		} else {
			e.coords.Right.X += 0.07 // 0.04 < dx < 0.09
		}
		if in, bound := checkBoundaries(); in {
			if e.r2l {
				e.coords.Left.X = bound.Right.X + 0.001
				e.coords.Right.X = e.coords.Left.X + 2*e.width
				e.r2l = false
			} else {
				e.coords.Right.X = bound.Left.X - 0.001
				e.coords.Left.X = e.coords.Right.X - 2*e.width
				e.r2l = true
			}
		} else {
			if e.r2l {
				e.coords.Right.X = e.coords.Left.X + 2*e.width
			} else {
				e.coords.Left.X = e.coords.Right.X - 2*e.width
			}
		}
	}
}

func (e *Enemy) GetCoords() Rect {
	return e.coords
}

func (e *Enemy) Render() {
	// TODO extract
	gl.PushMatrix()
	{
		gl.Translatef(e.coords.Left.X+e.width, e.coords.Right.Y+e.height, 0)

		var rect Rect
		if e.r2l {
			rect = Rect{
				Left:  Point{e.width, -e.height},
				Right: Point{-e.width, e.height},
			}
		} else {
			rect = Rect{
				Left:  Point{-e.width, -e.height},
				Right: Point{e.width, e.height},
			}
		}
		DrawTexture(e.texture, rect)
	}
	gl.PopMatrix()
}

func (e *Enemy) Unload() {
	gl.DeleteTextures(1, &e.texture)
}
