package main

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

type Scoreboard struct {
	life   uint32
	digits []uint32
	keys   uint32
	gover  uint32
	area   Rect
	lh, lw float32 // height, width of life image
	dh, dw float32 // height, width of digit image
	kh, kw float32 // height, width of keys image
	gh, gw float32 // height, width of gover image
}

func NewScoreBoard(area Rect) Board {
	life, bounds := NewTexture("assets/life.png")
	lh, lw := bounds.Right.Y, bounds.Right.X

	var dh, dw float32
	digits := make([]uint32, 10)
	for i := 0; i < 10; i++ {
		digit, bounds := NewTexture(fmt.Sprintf("assets/%d.png", i))
		digits[i] = digit
		dh, dw = bounds.Right.Y, bounds.Right.X
	}

	keys, bounds := NewTexture("assets/keys.png")
	kh, kw := bounds.Right.Y, bounds.Right.X

	gover, bounds := NewTexture("assets/gameover.png")
	gh, gw := bounds.Right.Y, bounds.Right.X

	return &Scoreboard{
		life:   life,
		digits: digits,
		keys:   keys,
		gover:  gover,
		area:   area,
		lh:     lh,
		lw:     lw,
		dh:     dh,
		dw:     dw,
		kh:     kh,
		kw:     kw,
		gh:     gh,
		gw:     gw,
	}
}

func (s *Scoreboard) Show(score, lives int) {
	const space = 0.3
	rect := Rect{
		Left:  Point{s.lw, s.lh},
		Right: Point{s.lw * 2, s.lh * 2},
	}
	for i := 0; i < lives; i++ {
		gl.PushMatrix()
		{
			gl.Translatef(s.area.Left.X+(float32(i)*(s.lw+space)), -s.area.Left.Y-s.lh+space, 0)
			DrawTexture(s.life, rect)
		}
		gl.PopMatrix()
	}
	rect = Rect{
		Left:  Point{s.dw, s.dh},
		Right: Point{s.dw * 2, s.dh * 2},
	}
	digits := IntToDigits(int64(score))
	for i, digit := range digits {
		gl.PushMatrix()
		{
			gl.Translatef(s.area.Right.X/2+(float32(i)*(s.dw+space)), -s.area.Left.Y-s.dh+space, 0)
			DrawTexture(s.digits[digit], rect)
		}
		gl.PopMatrix()
	}
	rect = Rect{
		Left:  Point{-s.kw, -s.kh},
		Right: Point{s.kw, s.kh},
	}
	gl.PushMatrix()
	{
		gl.Translatef(0, s.area.Right.Y+s.kh*1.5, 0)
		DrawTexture(s.keys, rect)
	}
	gl.PopMatrix()
	if lives <= 0 {
		rect = Rect{
			Left:  Point{-s.gw, -s.gh},
			Right: Point{s.gw, s.gh},
		}
		gl.PushMatrix()
		{
			gl.Translatef(0, 0, 0)
			DrawTexture(s.gover, rect)
		}
		gl.PopMatrix()
	}
}

func (s *Scoreboard) Unload() {
	gl.DeleteTextures(1, &s.life)
	for _, digit := range s.digits {
		gl.DeleteTextures(1, &digit)
	}
}
