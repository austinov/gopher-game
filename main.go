package main

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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
	// set position on the screen
	window.SetPos(100, 100)
	// set icon
	if err := setIcon(window); err != nil {
		log.Printf("unable to set window icon: %v\n", err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	game := NewGame(window)

	for !window.ShouldClose() {
		game.Update()
		game.Render()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func setIcon(window *glfw.Window) error {
	imgFile, err := os.Open("assets/gopher.png")
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	window.SetIcon([]image.Image{img})
	return nil
}
