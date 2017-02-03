package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type Scoreboard struct {
	life   uint32
	digits []uint32
	area   Rect
	lh     float32 // height of life image
	lw     float32 // width of life image
	dh     float32 // height of digit image
	dw     float32 // width of digit image
}

func NewScoreBoard(area Rect) *Scoreboard {
	life, bounds := NewTexture("assets/life.png")
	lh, lw := bounds.Right.Y, bounds.Right.X

	var dh, dw float32
	digits := make([]uint32, 10)
	for i := 0; i < 10; i++ {
		digit, bounds := NewTexture(fmt.Sprintf("assets/%d.png", i))
		digits[i] = digit
		dh, dw = bounds.Right.Y, bounds.Right.X
	}
	return &Scoreboard{
		life:   life,
		digits: digits,
		area:   area,
		lh:     lh,
		lw:     lw,
		dh:     dh,
		dw:     dw,
	}
}

func (s *Scoreboard) Show(score int32, lives int) {
	rect := Rect{
		Left:  Point{s.lw, s.lh},
		Right: Point{s.lw * 2, s.lh * 2},
	}
	const space = 0.3
	for i := 0; i < lives; i++ {
		gl.PushMatrix()
		{
			gl.Translatef(s.area.Left.X+(float32(i)*(s.lw+space)), -s.area.Left.Y-s.lh+space, 0)
			DrawTexture(s.life, rect)
		}
		gl.PopMatrix()
	}
	score = 125
	rect = Rect{
		Left:  Point{s.dw, s.dh},
		Right: Point{s.dw * 2, s.dh * 2},
	}
	count := CountOfDigits(int64(score))
	for i := 1; i <= count; i++ {
		// TODO
		digit := score % int32(math.Pow(10, float64(i)))
		fmt.Printf("count=%d, digit=%d, pow=%v\n", count, digit, math.Pow(10, float64(i)))
		gl.PushMatrix()
		{
			gl.Translatef(s.area.Right.Y-(float32(i)*(s.dw+space)), -s.area.Left.Y-s.dh+space, 0)
			DrawTexture(s.digits[digit], rect)
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