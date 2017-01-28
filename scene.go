package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type scene struct {
	window  *glfw.Window
	texture uint32
	rect    Rect
	bounds  Rect
	gopher  *Gopher
	enemy   *Enemy
	mmap    []Rect
}

func NewScene(window *glfw.Window) Scene {
	// wich part of the window will be used for rendering
	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	screenRatio := float32(height) / float32(width)
	sceneRect := Rect{
		Left:  Coord{-16, -10},
		Right: Coord{16, 10},
	}
	sceneSize := Coord{
		X: sceneRect.Right.X - sceneRect.Left.X,
		Y: sceneRect.Right.Y - sceneRect.Left.Y,
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

	bounds := Rect{
		Left: Coord{
			X: -screenSize.X / 2, // left
			Y: -screenSize.Y / 2, // right
		},
		Right: Coord{
			X: screenSize.X / 2, // bottom
			Y: screenSize.Y / 2, // top
		},
	}
	return &scene{
		window:  window,
		texture: newTexture("assets/scene.png"),
		rect:    sceneRect, // TODO refactoring
		bounds:  bounds,    // TODO refactoring
		mmap:    loadMap(),
	}
}

func (s *scene) Boundaries() []Rect {
	return s.mmap
}

func (s *scene) Update() {
}

func (s *scene) Render() {
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

	// multiply the current matrix with an orthographic matrix
	gl.Ortho(
		float64(s.bounds.Left.X),  // left
		float64(s.bounds.Right.X), // right
		float64(s.bounds.Left.Y),  // bottom
		float64(s.bounds.Right.Y), // top
		10, -10) // distances to the nearer and farther depth clipping planes

	/*
		x, y, z := float32(0.0), float32(0.0), float32(0.0)
		gl.Translatef(x, y, z)
	*/

	gl.Enable(gl.MULTISAMPLE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.ALPHA_TEST)

	drawTexture(s.texture, s.rect)
}

func (s *scene) Unload() {
	gl.DeleteTextures(1, &s.texture)
}

func loadMap() []Rect {
	bounds := make([]Rect, 0)
	// walls
	bounds = append(bounds, Rect{
		Left:  Coord{-16, 10},
		Right: Coord{-16, -10},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{-16, 10},
		Right: Coord{16, 10},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{16, 10},
		Right: Coord{16, -10},
	})
	// platforms
	bounds = append(bounds, Rect{
		Left:  Coord{-16, -4},
		Right: Coord{-13, -4},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{13, -4},
		Right: Coord{16, -4},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{-14, -1.5},
		Right: Coord{14, -1.5},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{-16, 1},
		Right: Coord{-13, 1},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{13, 1},
		Right: Coord{16, 1},
	})
	bounds = append(bounds, Rect{
		Left:  Coord{-14, 3.5},
		Right: Coord{14, 3.5},
	})
	return bounds
}
