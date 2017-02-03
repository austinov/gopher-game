package main

import (
	"math/rand"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	window  *glfw.Window
	scene   Scene
	gopher  Player
	enemies []Player
	board   *Scoreboard
	score   int32
	lives   int
}

func NewGame(window *glfw.Window) *Game {
	rand.Seed(time.Now().UTC().UnixNano())
	scene := NewScene(window)
	return &Game{
		window:  window,
		scene:   scene,
		gopher:  NewGopher(window, scene),
		enemies: make([]Player, 0),
		board:   NewScoreBoard(scene.GetArea()),
		lives:   3,
	}
}

func (g *Game) Update() {
	g.scene.Update()
	g.checkLives()
	if g.lives <= 0 {
		return
	}
	g.gopher.Update(g.window, g.scene)
	g.updateEnemies()
}

func (g *Game) checkLives() {
	if g.lives <= 0 {
		return
	}
	if v, _ := CheckBoundaries(g.gopher.GetCoords(), g.scene.GetHole()); v {
		// gopher crossed the hole, delete it and create new one
		g.lives--
		if g.lives <= 0 {
			return
		}
		g.gopher.Unload()
		g.gopher = NewGopher(g.window, g.scene)
	}
	caught := false
	for _, enemy := range g.enemies {
		// gopher was caught by enemy, delete it and create new one
		if v, _ := CheckBoundaries(g.gopher.GetCoords(), enemy.GetCoords()); v {
			caught = true
			break
		}
	}
	if caught {
		for _, enemy := range g.enemies {
			enemy.Unload()
		}
		g.enemies = make([]Player, 0)
		g.lives--
		if g.lives <= 0 {
			return
		}
		g.gopher.Unload()
		g.gopher = NewGopher(g.window, g.scene)
	}
}

func (g *Game) updateEnemies() {
	// TODO
	deleted := 0
	for i := range g.enemies {
		j := i - deleted
		if v, _ := CheckBoundaries(g.enemies[j].GetCoords(), g.scene.GetHole()); v {
			// enemy crossed the hole, delete it and create new one
			g.enemies[j].Unload()
			g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
			deleted++
		}
	}
	if len(g.enemies) == 0 {
		g.enemies = append(g.enemies, NewEnemy(rand.Intn(2) == 0))
	}
	for _, enemy := range g.enemies {
		enemy.Update(g.window, g.scene)
	}
}

func (g *Game) Render() {
	g.scene.Render()
	g.gopher.Render()
	g.board.Show(g.score, g.lives)
	for _, enemy := range g.enemies {
		enemy.Render()
	}
}

func (g *Game) Unload() {
	g.scene.Unload()
	g.gopher.Unload()
	for _, enemy := range g.enemies {
		enemy.Unload()
	}
	g.board.Unload()
}
