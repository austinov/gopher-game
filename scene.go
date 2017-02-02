package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type scene struct {
	texture uint32
	rect    Rect
	bounds  Rect
	plan    []Rect
}

func NewScene(window *glfw.Window) Scene {
	// wich part of the window will be used for rendering
	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	screenRatio := float32(height) / float32(width)
	sceneRect := Rect{
		Left:  Point{-16, -10},
		Right: Point{16, 10},
	}
	sceneSize := Point{
		X: sceneRect.Right.X - sceneRect.Left.X,
		Y: sceneRect.Right.Y - sceneRect.Left.Y,
	}

	var screenSize Point
	sceneRatio := sceneSize.Y / sceneSize.X

	if screenRatio < sceneRatio {
		screenSize.Y = sceneSize.Y + 2
		screenSize.X = screenSize.Y / screenRatio
	} else {
		screenSize.X = sceneSize.X + 2
		screenSize.Y = screenSize.X * screenRatio
	}

	bounds := Rect{
		Left: Point{
			X: -screenSize.X / 2, // left
			Y: -screenSize.Y / 2, // right
		},
		Right: Point{
			X: screenSize.X / 2, // bottom
			Y: screenSize.Y / 2, // top
		},
	}
	return &scene{
		texture: NewTexture("assets/scene.png"),
		rect:    sceneRect, // TODO refactoring
		bounds:  bounds,    // TODO refactoring
		plan:    loadMap(),
	}
}

func (s *scene) Boundaries() []Rect {
	return s.plan
}

func (s *scene) Update() {
}

func (s *scene) Render() {
	// specify clear values for the color buffers (r,g,b,a)
	gl.ClearColor(0, 0, 0, 1)
	// clear buffers to preset values
	// COLOR_BUFFER_BIT indicates the buffers coordsently enabled for color writing
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// specify which matrix is the coordsent matrix
	// GL_MODELVIEW applies subsequent matrix operations to the modelview matrix stack
	gl.MatrixMode(gl.MODELVIEW)
	// replace the coordsent matrix with the identity matrix
	gl.LoadIdentity()

	// multiply the coordsent matrix with an orthographic matrix
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

	DrawTexture(s.texture, s.rect)
}

func (s *scene) Unload() {
	gl.DeleteTextures(1, &s.texture)
}

func loadMap() []Rect {
	bounds := make([]Rect, 0)
	// platforms
	// left first from bottom
	bounds = append(bounds, Rect{
		Left:  Point{-16, -8},
		Right: Point{-10, -8.4},
	})
	// right first bottom
	bounds = append(bounds, Rect{
		Left:  Point{10, -8},
		Right: Point{16, -8.4},
	})
	// middle second bottom
	bounds = append(bounds, Rect{
		Left:  Point{-12, -3},
		Right: Point{12, -3.4},
	})
	// left middle
	bounds = append(bounds, Rect{
		Left:  Point{-16, 2},
		Right: Point{-10, 1.6},
	})
	// right middle
	bounds = append(bounds, Rect{
		Left:  Point{10, 2},
		Right: Point{16, 1.6},
	})
	// top middle
	bounds = append(bounds, Rect{
		Left:  Point{-12, 7},
		Right: Point{12, 6.6},
	})
	// walls
	// left
	bounds = append(bounds, Rect{
		Left:  Point{-16, 10},
		Right: Point{-15.6, -10},
	})
	// top
	bounds = append(bounds, Rect{
		Left:  Point{-16, 10},
		Right: Point{16, 9.6},
	})
	// right
	bounds = append(bounds, Rect{
		Left:  Point{15.6, 10},
		Right: Point{16, -10},
	})
	return bounds
}
