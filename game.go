package main

import (
	"math/rand"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	giftPoints      = 10
	enemyPoints     = 5
	reserveLives    = 3
	maxDelayGiftSec = 10
	maxGifts        = 5
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Game struct {
	window    *glfw.Window
	scene     Scene
	gopher    Player
	enemies   []Player
	gifts     []Gift
	board     Board
	score     int
	lives     int
	lastGift  time.Time
	delayGift int
}

func NewGame(window *glfw.Window) *Game {
	rand.Seed(time.Now().UTC().UnixNano())
	scene := NewScene(window)
	return &Game{
		window:    window,
		scene:     scene,
		gopher:    NewGopher(window, scene),
		enemies:   make([]Player, 0), // TODO extract into controller
		gifts:     make([]Gift, 0),   // TODO extract into controller
		board:     NewScoreBoard(scene.GetArea()),
		lives:     reserveLives,
		lastGift:  time.Now(),      // TODO extract into controller
		delayGift: maxDelayGiftSec, // TODO extract into controller
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
	g.updateGifts()
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
	if g.isGopherCaught() {
		g.lives--
		if g.lives <= 0 {
			return
		}
		g.gopher.Unload()
		g.gopher = NewGopher(g.window, g.scene)
	}
}

func (g *Game) isGopherCaught() bool {
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
	}
	return caught
}

func (g *Game) updateEnemies() {
	// TODO more intelligence
	deleted := 0
	for i := range g.enemies {
		j := i - deleted
		if v, _ := CheckBoundaries(g.enemies[j].GetCoords(), g.scene.GetBlackHole()); v {
			// enemy crossed the hole, delete it and create new one
			g.enemies[j].Unload()
			g.enemies = append(g.enemies[:j], g.enemies[j+1:]...)
			deleted++
			g.score += enemyPoints
		}
	}
	if len(g.enemies) == 0 {
		g.enemies = append(g.enemies, NewEnemy(rand.Intn(2) == 0))
	}
	for _, enemy := range g.enemies {
		enemy.Update(g.window, g.scene)
	}
}

// TODO extract
func (g *Game) updateGifts() {
	deleted := 0
	for i := range g.gifts {
		j := i - deleted
		if v, _ := CheckBoundaries(g.gopher.GetCoords(), g.gifts[j].GetCoords()); v {
			// gopher picked up a gift, delete gift
			g.gifts[j].Unload()
			g.gifts = append(g.gifts[:j], g.gifts[j+1:]...)
			deleted++
			g.score += giftPoints
		}
	}
	if deleted > 0 {
		g.lastGift = time.Now()
		g.delayGift = rnd.Intn(maxDelayGiftSec)
		return
	}
	if time.Now().After(g.lastGift.Add(time.Duration(g.delayGift) * time.Second)) {
		needGift := rnd.Intn(maxGifts)
		for {
			if len(g.gifts) >= needGift {
				break
			}
			nextPoint := g.nextGiftPoint()
			g.gifts = append(g.gifts, NewNutGift(nextPoint))
		}
	}
}

func (g *Game) nextGiftPoint() Point {
	area := g.scene.GetArea()
	bounds := make([]Rect, 0)
	for _, b := range g.scene.GetBoundaries() {
		if b.Left.Y < area.Right.Y {
			bounds = append(bounds, b)
		}
	}
	bound := bounds[rnd.Intn(len(bounds))]
	dx := 0
	for dx == 0 {
		dx = rnd.Intn(int(bound.Right.X - bound.Left.X))
	}
	grect := g.gopher.GetCoords()
	dy := abs(abs(grect.Right.Y)-abs(grect.Left.Y)) / 2
	return Point{X: bound.Left.X + float32(dx), Y: bound.Left.Y + dy}
}

func abs(f float32) float32 {
	if f < 0.0 {
		return -f
	}
	return f
}

func (g *Game) Render() {
	g.scene.Render()
	g.gopher.Render()
	g.board.Show(g.score, g.lives)
	for _, enemy := range g.enemies {
		enemy.Render()
	}
	for _, gift := range g.gifts {
		gift.Render()
	}
}

func (g *Game) Unload() {
	g.scene.Unload()
	g.gopher.Unload()
	for _, enemy := range g.enemies {
		enemy.Unload()
	}
	for _, gift := range g.gifts {
		gift.Unload()
	}
	g.board.Unload()
}
