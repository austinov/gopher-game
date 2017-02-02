package main

import "github.com/go-gl/glfw/v3.2/glfw"

type Game struct {
	window  *glfw.Window
	scene   Scene
	gopher  Player
	enemies []Player
}

func NewGame(window *glfw.Window) *Game {
	scene := NewScene(window)
	return &Game{
		window:  window,
		scene:   scene,
		gopher:  NewGopher(window, scene),
		enemies: make([]Player, 0),
	}
}

func (g *Game) Update() {
	g.scene.Update()
	g.gopher.Update(g.window, g.scene)
	// TODO
	/*
		deleted := 0
		for i := range g.enemies {
			j := i - deleted
			if g.enemies[j].IsDead(g.scene.Boundaries()) {
				// delete
				g.enemies[j].Unload()
				g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
				deleted++
			}
		}
	*/
	if len(g.enemies) == 0 {
		g.enemies = append(g.enemies, NewEnemy())
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
