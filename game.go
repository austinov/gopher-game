package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	window  *glfw.Window
	scene   Scene
	gopher  Player
	enemies []Player
	score   uint32
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
		lives:   3,
	}
}

func (g *Game) Update() {
	g.scene.Update()
	g.showScore()
	if g.lives <= 0 {
		return
	}
	g.calcLives()
	g.gopher.Update(g.window, g.scene)
	g.updateEnemies()
}

func (g *Game) calcLives() {
	if g.lives <= 0 {
		return
	}
	if v, _ := CheckBoundaries(g.gopher.GetCoords(), g.scene.GetHole()); v {
		g.lives--
		g.gopher.Unload()
		if g.lives <= 0 {
			return
		}
		g.gopher = NewGopher(g.window, g.scene)
	}
	for _, enemy := range g.enemies {
		if v, _ := CheckBoundaries(g.gopher.GetCoords(), enemy.GetCoords()); v {
			g.lives--
			g.gopher.Unload()
			if g.lives <= 0 {
				return
			}
			g.gopher = NewGopher(g.window, g.scene)
		}
	}
}

func (g *Game) showScore() {
	// TODO
	for i := 0; i < g.lives; i++ {
		//star.Render()
	}
	// score.Render()
	if g.lives <= 0 {
		fmt.Printf("*** GAME OVER ***\n")
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
}
