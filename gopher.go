package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/paulmach/go.geo"
)

// TODO refactoring
var (
	isLeft, isRight, isTop, isBottom bool
)

type Gopher struct {
	window  *glfw.Window
	texture uint32
	plan    Plan
	curr    Rect
	h, w    float32
	r2l     bool
	onFloor bool
}

func NewGopher(window *glfw.Window, plan Plan) GameElement {
	h, w := float32(0.8), float32(0.8)
	return &Gopher{
		window:  window,
		texture: NewTexture("assets/gopher.png"),
		plan:    plan,
		h:       h,
		w:       w,
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
	isLeft = isButtonPress(glfw.KeyLeft)
	isRight = isButtonPress(glfw.KeyRight)
	isTop = isButtonPress(glfw.KeyUp)
	isBottom = isButtonPress(glfw.KeyDown)

	if isLeft {
		g.curr.Left.X -= 0.1
		if in, bound := g.intersects(g.curr); in {
			g.curr.Left.X = bound.Right.X + 0.001
		}
		g.curr.Right.X = g.curr.Left.X + 2*g.w
	}
	if isRight {
		g.curr.Right.X += 0.1
		if in, bound := g.intersects(g.curr); in {
			g.curr.Right.X = bound.Left.X - 0.001
		}
		g.curr.Left.X = g.curr.Right.X - 2*g.w
	}
	if g.onFloor {
		if (g.curr.Left.X < floor.Left.X && g.curr.Right.X-g.w/2 < floor.Left.X) ||
			(g.curr.Right.X > floor.Right.X && g.curr.Left.X+g.w/2 > floor.Right.X) {
			g.onFloor = false
			floor = Rect{}
		}
	}
	if isTop && g.onFloor {
		g.curr.Left.Y += 6.6
		if in, bound := g.intersects(g.curr); in {
			g.curr.Left.Y = bound.Right.Y - 0.001
		}
		g.curr.Right.Y = g.curr.Left.Y - 2*g.h
		g.onFloor = false
	}
	if !g.onFloor {
		g.curr.Right.Y -= 0.05
		if in, bound := g.intersects(g.curr); in {
			g.curr.Right.Y = bound.Left.Y + 0.001
			g.onFloor = true
			floor = bound
		}
		g.curr.Left.Y = g.curr.Right.Y + 2*g.h
	}
}

var floor Rect

func (g *Gopher) Render() {
	gl.PushMatrix()
	{
		gl.Translatef(g.curr.Left.X+g.w, g.curr.Right.Y+g.h, 0)

		rect := Rect{
			Left:  Coord{-g.w, -g.h},
			Right: Coord{g.w, g.h},
		}
		r2l := g.r2l
		if isLeft {
			r2l = true
		}
		if isRight {
			r2l = false
		}
		// set previouse direction
		if g.r2l {
			rect.Left.X *= float32(-1.0)
			rect.Right.X *= float32(-1.0)
		}
		if r2l != g.r2l {
			// direction was changed, set new one
			rect.Left.X *= float32(-1.0)
			rect.Right.X *= float32(-1.0)
			g.r2l = r2l
		}
		DrawTexture(g.texture, rect)
	}
	gl.PopMatrix()
}

func (g *Gopher) Unload() {
	gl.DeleteTextures(1, &g.texture)
}

func (g *Gopher) intersects(r1 Rect) (bool, Rect) {
	for _, r2 := range g.plan.Boundaries() {
		p1 := geo.NewPath()
		pp := []geo.Point{
			geo.Point{float64(r1.Left.X), float64(r1.Left.Y)}, geo.Point{float64(r1.Right.X), float64(r1.Left.Y)}, // top
			geo.Point{float64(r1.Right.X), float64(r1.Left.Y)}, geo.Point{float64(r1.Right.X), float64(r1.Right.Y)}, // right
			geo.Point{float64(r1.Right.X), float64(r1.Right.Y)}, geo.Point{float64(r1.Left.X), float64(r1.Right.Y)}, // bottom
			geo.Point{float64(r1.Left.X), float64(r1.Right.Y)}, geo.Point{float64(r1.Left.X), float64(r1.Left.Y)}, // left
		}
		p1.SetPoints(pp)
		p2 := geo.NewPath()
		pp = []geo.Point{
			geo.Point{float64(r2.Left.X), float64(r2.Left.Y)}, geo.Point{float64(r2.Right.X), float64(r2.Left.Y)}, // top
			geo.Point{float64(r2.Right.X), float64(r2.Left.Y)}, geo.Point{float64(r2.Right.X), float64(r2.Right.Y)}, // right
			geo.Point{float64(r2.Right.X), float64(r2.Right.Y)}, geo.Point{float64(r2.Left.X), float64(r2.Right.Y)}, // bottom
			geo.Point{float64(r2.Left.X), float64(r2.Right.Y)}, geo.Point{float64(r2.Left.X), float64(r2.Left.Y)}, // left
		}
		p2.SetPoints(pp)

		if p1.Intersects(p2) {
			return true, r2
		}
	}
	return false, Rect{}
}
