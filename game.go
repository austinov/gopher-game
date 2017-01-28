package main

import "github.com/go-gl/glfw/v3.2/glfw"

type Game struct {
	window *glfw.Window
	scene  Scene
	gopher GameElement
	enemy  GameElement
}

func NewGame(window *glfw.Window) *Game {
	return &Game{
		scene:  NewScene(window),
		gopher: NewGopher(window),
		enemy:  NewEnemy(window),
	}
}

func (g *Game) Update() {
	g.scene.Update()
	g.gopher.Update()
	g.enemy.Update()
}

func (g *Game) Render() {
	g.scene.Render()
	g.gopher.Render()
	g.enemy.Render()
}

func (g *Game) Unload() {
	g.scene.Unload()
	g.gopher.Unload()
	g.enemy.Unload()
}
