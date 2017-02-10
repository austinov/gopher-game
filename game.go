package main

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	giftPoints   = 10
	enemyPoints  = 5
	reserveLives = 3
)

type Game struct {
	window *glfw.Window
	scene  Scene
	gopher Player
	arts   *Artefactory
	board  Board
	score  int
	lives  int
}

func NewGame(window *glfw.Window) *Game {
	scene := NewScene(window)
	g := &Game{
		window: window,
		scene:  scene,
		arts:   NewArtefactory(window),
		board:  NewScoreBoard(scene.GetArea()),
	}
	g.init()
	return g
}

func (g *Game) Update() {
	g.scene.Update()
	g.checkLives()
	if g.lives <= 0 {
		if g.window.GetKey(glfw.KeySpace) == glfw.Press {
			g.init()
		} else {
			return
		}
	}
	g.gopher.Update(g.window, g.scene)
	caught, enemies, gifts := g.arts.Update(g.scene, g.gopher)
	if caught {
		g.lives--
		if g.lives <= 0 {
			return
		}
		g.gopher.Unload()
		g.gopher = NewGopher(g.window, g.scene)
	}
	g.score += enemies*enemyPoints + gifts*giftPoints
}

func (g *Game) init() {
	g.gopher = NewGopher(g.window, g.scene)
	g.lives = reserveLives
	g.score = 0
}

func (g *Game) checkLives() {
	if g.lives <= 0 {
		return
	}
	if v, _ := CheckBoundaries(g.gopher.GetCoords(), g.scene.GetBlackHole()); v {
		// gopher crossed the hole, delete it and create new one
		g.lives--
		if g.lives <= 0 {
			return
		}
		g.gopher.Unload()
		g.gopher = NewGopher(g.window, g.scene)
	}
}

func (g *Game) Render() {
	g.scene.Render()
	g.gopher.Render()
	g.board.Show(g.score, g.lives)
	g.arts.Render()
}

func (g *Game) Unload() {
	g.scene.Unload()
	g.gopher.Unload()
	g.board.Unload()
	g.arts.Unload()
}
