package main

import (
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

func newTexture(file string) uint32 {
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

	var texture uint32

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

	return texture
}

func drawTexture(texture uint32, dst Rect) {
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
