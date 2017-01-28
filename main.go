package main

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Coord struct {
	X, Y float32
}

type Rect struct {
	Left, Right Coord
}

var (
	textureScene  uint32
	textureGopher uint32
	textureEnemy  uint32
	delta         Coord
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

	window, err := glfw.CreateWindow(800, 600, "Gopher Game", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetPos(100, 100)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	textureScene = newTexture("assets/scene.png")
	defer gl.DeleteTextures(1, &textureScene)

	textureGopher = newTexture("assets/gopher.png")
	defer gl.DeleteTextures(1, &textureGopher)

	textureEnemy = newTexture("assets/enemy.png")
	defer gl.DeleteTextures(1, &textureEnemy)

	s := newScene(window)

	for !window.ShouldClose() {
		s.Update()
		s.Render()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
