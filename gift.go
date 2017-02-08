package main

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type NutGift struct {
	texture uint32
	coords  Rect
	height  float32
	width   float32
}

func NewNutGift(point Point) Gift {
	texture, bounds := NewTexture("assets/nut.png")
	h, w := bounds.Right.Y, bounds.Right.X
	return &NutGift{
		texture: texture,
		coords: Rect{
			Left:  Point{point.X - w/2, point.Y + h/2},
			Right: Point{point.X + w/2, point.Y - h/2},
		},
		height: bounds.Right.Y,
		width:  bounds.Right.X,
	}
}

func (e *NutGift) GetCoords() Rect {
	return e.coords
}

func (e *NutGift) Render() {
	gl.PushMatrix()
	{
		gl.Translatef(e.coords.Left.X+e.width, e.coords.Right.Y+e.height, 0)

		rect := Rect{
			Left:  Point{e.width, -e.height},
			Right: Point{-e.width, e.height},
		}
		DrawTexture(e.texture, rect)
	}
	gl.PopMatrix()
}

func (e *NutGift) Unload() {
	gl.DeleteTextures(1, &e.texture)
}
