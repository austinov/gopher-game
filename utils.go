package main

import (
	"image"
	"image/draw"
	"log"
	"math"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/paulmach/go.geo"
)

// NewTexture loads texture from file.
func NewTexture(file string) (texture uint32, bounds Rect) {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	//var texture uint32

	// enable server-side GL capabilities
	gl.Enable(gl.TEXTURE_2D)
	// generate texture names
	gl.GenTextures(1, &texture)
	// bind a named texture to a texturing target
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// set texture parameters
	// texture minifying function
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	// texture magnification function
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// set the wrap parameter for texture coordinate s (x dimension)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	// set the wrap parameter for texture coordinate t (y dimension)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// set the wrap parameter for texture coordinate r (z dimension)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	// specify a two-dimensional texture image
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	// disable server-side GL capabilities
	gl.Disable(gl.TEXTURE_2D)

	// translate image size into game units
	imgWidth := float32(img.Bounds().Max.X) / 100.0
	imgHeight := float32(img.Bounds().Max.Y) / 100.0
	bounds = Rect{
		Left:  Point{X: float32(0), Y: float32(0)},
		Right: Point{X: imgWidth, Y: imgHeight},
	}
	return texture, bounds
}

// DrawTexture draws texture into rectange.
func DrawTexture(texture uint32, dst Rect) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)
	gl.Begin(gl.QUADS)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(dst.Left.X, dst.Left.Y)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(dst.Right.X, dst.Left.Y)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(dst.Right.X, dst.Right.Y)

	gl.TexCoord2f(0, 0)
	gl.Vertex2f(dst.Left.X, dst.Right.Y)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

// CheckBoundaries checks violation of boundaries.
func CheckBoundaries(rect Rect, boundaries ...Rect) (violated bool, violatedBound Rect) {
	for _, bound := range boundaries {
		p1 := geo.NewPath()
		pp := []geo.Point{
			geo.Point{float64(rect.Left.X), float64(rect.Left.Y)}, geo.Point{float64(rect.Right.X), float64(rect.Left.Y)}, // top
			geo.Point{float64(rect.Right.X), float64(rect.Left.Y)}, geo.Point{float64(rect.Right.X), float64(rect.Right.Y)}, // right
			geo.Point{float64(rect.Right.X), float64(rect.Right.Y)}, geo.Point{float64(rect.Left.X), float64(rect.Right.Y)}, // bottom
			geo.Point{float64(rect.Left.X), float64(rect.Right.Y)}, geo.Point{float64(rect.Left.X), float64(rect.Left.Y)}, // left
		}
		p1.SetPoints(pp)
		p2 := geo.NewPath()
		pp = []geo.Point{
			geo.Point{float64(bound.Left.X), float64(bound.Left.Y)}, geo.Point{float64(bound.Right.X), float64(bound.Left.Y)}, // top
			geo.Point{float64(bound.Right.X), float64(bound.Left.Y)}, geo.Point{float64(bound.Right.X), float64(bound.Right.Y)}, // right
			geo.Point{float64(bound.Right.X), float64(bound.Right.Y)}, geo.Point{float64(bound.Left.X), float64(bound.Right.Y)}, // bottom
			geo.Point{float64(bound.Left.X), float64(bound.Right.Y)}, geo.Point{float64(bound.Left.X), float64(bound.Left.Y)}, // left
		}
		p2.SetPoints(pp)

		if p1.Intersects(p2) {
			return true, bound
		}
	}
	return false, Rect{}
}

// CountOfDigits returns count of digits in number
func CountOfDigits(number int64) int {
	if number == 0 {
		return 1
	}
	if number < 0 {
		number = -number
	}
	return int(math.Ceil(math.Log10(math.Abs(float64(number)) + 0.5)))
}
