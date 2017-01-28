package main

import (
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Coord struct {
	X, Y float32
}

type Rect struct {
	Min, Max Coord
}

var (
	textureScene  uint32
	textureGopher uint32
	textureDragon uint32
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Samples, 4)

	//mode := glfw.GetPrimaryMonitor().GetVideoMode()

	window, err := glfw.CreateWindow(800, 600, "Gopher Game", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetPos(100, 100)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	textureScene = newTexture("scene.png")
	defer gl.DeleteTextures(1, &textureScene)

	textureGopher = newTexture("gopher.png")
	defer gl.DeleteTextures(1, &textureGopher)

	textureDragon = newTexture("dragon.png")
	defer gl.DeleteTextures(1, &textureDragon)

	for !window.ShouldClose() {
		drawScene(window)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

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

func drawScene(window *glfw.Window) {
	// specify clear values for the color buffers (r,g,b,a)
	gl.ClearColor(0, 0, 0, 1)
	// clear buffers to preset values
	// COLOR_BUFFER_BIT indicates the buffers currently enabled for color writing
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// specify which matrix is the current matrix
	// GL_MODELVIEW applies subsequent matrix operations to the modelview matrix stack
	gl.MatrixMode(gl.MODELVIEW)
	// replace the current matrix with the identity matrix
	gl.LoadIdentity()

	// wich part of the window will be used for rendering
	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	screenRatio := float32(height) / float32(width)
	// TODO extract
	sceneRect := Rect{
		Min: Coord{-22, -16},
		Max: Coord{22, 16},
	}
	sceneSize := Coord{
		X: sceneRect.Max.X - sceneRect.Min.X,
		Y: sceneRect.Max.Y - sceneRect.Min.Y,
	}

	var screenSize Coord
	sceneRatio := sceneSize.Y / sceneSize.X

	if screenRatio < sceneRatio {
		screenSize.Y = sceneSize.Y + 2
		screenSize.X = screenSize.Y / screenRatio
	} else {
		screenSize.X = sceneSize.X + 2
		screenSize.Y = screenSize.X * screenRatio
	}

	// multiply the current matrix with an orthographic matrix
	gl.Ortho(
		float64(-screenSize.X/2), // left
		float64(screenSize.X/2),  // right
		float64(-screenSize.Y/2), // bottom
		float64(screenSize.Y/2),  // top
		10, -10) // distances to the nearer and farther depth clipping planes

	/*
		x, y, z := float32(0.0), float32(0.0), float32(0.0)
		gl.Translatef(x, y, z)
	*/

	gl.Enable(gl.MULTISAMPLE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.ALPHA_TEST)

	// render game's parts
	drawTexture(textureScene, sceneRect)

	/*
		gopherRect := Rect{
			Min: V2{-0.8, -0.8},
			Max: V2{0.8, 0.8},
		}
		drawTexture(textureGopher, gopherRect)
	*/

	dragonRect := Rect{
		Min: Coord{-0.8, -0.8},
		Max: Coord{0.8, 0.8},
	}
	drawTexture(textureDragon, dragonRect)

	playerRender()
}

func drawTexture(texture uint32, dst Rect) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)
	gl.Begin(gl.QUADS)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(dst.Min.X, dst.Min.Y)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(dst.Max.X, dst.Min.Y)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(dst.Max.X, dst.Max.Y)

	gl.TexCoord2f(0, 0)
	gl.Vertex2f(dst.Min.X, dst.Max.Y)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

func playerRender() {
	gl.PushMatrix()
	{
		v := Coord{X: -8.614434, Y: 2.4074345}
		gl.Translatef(v.X, v.Y, 0)

		gopherRect := Rect{
			Min: Coord{-0.8, -0.8},
			Max: Coord{0.8, 0.8},
		}
		drawTexture(textureGopher, gopherRect)
	}
	gl.PopMatrix()
}
