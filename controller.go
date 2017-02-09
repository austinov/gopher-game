package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"math/rand"
	"time"
)

const (
	maxDelayGiftSec  = 10
	maxDelayEnemySec = 5
	maxGifts         = 5
	maxEnemies       = 4
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Artefactory struct {
	window     *glfw.Window
	enemies    []Player
	gifts      []Gift
	lastGift   time.Time
	delayGift  int
	lastEnemy  time.Time
	delayEnemy int
}

func NewArtefactory(window *glfw.Window) *Artefactory {
	return &Artefactory{
		window:     window,
		enemies:    make([]Player, 0),
		gifts:      make([]Gift, 0),
		lastGift:   time.Unix(0, 0),
		delayGift:  maxDelayGiftSec,
		lastEnemy:  time.Unix(0, 0),
		delayEnemy: maxDelayEnemySec,
	}
}

func (a *Artefactory) Update(scene Scene, gopher Player) (caught bool, deadEnemies, collectedGifts int) {
	caught = a.isGopherCaught(gopher)
	if caught {
		return true, deadEnemies, collectedGifts
	}
	deadEnemies = a.updateEnemies(scene)
	collectedGifts = a.updateGifts(scene, gopher)
	return false, deadEnemies, collectedGifts
}

func (a *Artefactory) Render() {
	for _, enemy := range a.enemies {
		enemy.Render()
	}
	for _, gift := range a.gifts {
		gift.Render()
	}
}

func (a *Artefactory) Unload() {
	for _, enemy := range a.enemies {
		enemy.Unload()
	}
	for _, gift := range a.gifts {
		gift.Unload()
	}
}

func (a *Artefactory) updateEnemies(scene Scene) (deadEnemies int) {
	// TODO more intelligence
	deleted := 0
	for i := range a.enemies {
		j := i - deleted
		if v, _ := CheckBoundaries(a.enemies[j].GetCoords(), scene.GetBlackHole()); v {
			// enemy crossed the hole, delete it and create new one
			a.enemies[j].Unload()
			a.enemies = append(a.enemies[:j], a.enemies[j+1:]...)
			deleted++
			//c.score += enemyPoints
			deadEnemies++
		}
	}
	if time.Now().After(a.lastEnemy.Add(time.Duration(a.delayEnemy)*time.Second)) || a.lastEnemy.Unix() == 0 {
		if len(a.enemies) < rnd.Intn(maxEnemies) {
			a.lastEnemy = time.Now()
			a.delayEnemy = rnd.Intn(maxDelayEnemySec)
			a.enemies = append(a.enemies, NewEnemy(rnd.Intn(2) == 0))
		}
	}
	for _, enemy := range a.enemies {
		enemy.Update(a.window, scene)
	}
	return deadEnemies
}

func (a *Artefactory) updateGifts(scene Scene, gopher Player) (collectedGifts int) {
	deleted := 0
	for i := range a.gifts {
		j := i - deleted
		if v, _ := CheckBoundaries(gopher.GetCoords(), a.gifts[j].GetCoords()); v {
			// gopher picked up a gift, delete gift
			a.gifts[j].Unload()
			a.gifts = append(a.gifts[:j], a.gifts[j+1:]...)
			deleted++
			//g.score += giftPoints
			collectedGifts++
		}
	}
	if deleted > 0 {
		a.lastGift = time.Now()
		a.delayGift = rnd.Intn(maxDelayGiftSec)
		return collectedGifts
	}
	if a.lastGift.Unix() == 0 {
		a.lastGift = time.Now()
	}
	if time.Now().After(a.lastGift.Add(time.Duration(a.delayGift) * time.Second)) {
		needGift := rnd.Intn(maxGifts)
		for {
			if len(a.gifts) >= needGift {
				break
			}
			nextPoint := a.nextGiftPoint(scene, gopher)
			a.gifts = append(a.gifts, NewNutGift(nextPoint))
		}
	}
	return collectedGifts
}

func (a *Artefactory) isGopherCaught(gopher Player) bool {
	caught := false
	for _, enemy := range a.enemies {
		// gopher was caught by enemy, delete it and create new one
		if v, _ := CheckBoundaries(gopher.GetCoords(), enemy.GetCoords()); v {
			caught = true
			break
		}
	}
	if caught {
		for _, enemy := range a.enemies {
			enemy.Unload()
		}
		a.enemies = make([]Player, 0)
	}
	return caught
}

func (a *Artefactory) nextGiftPoint(scene Scene, gopher Player) Point {
	area := scene.GetArea()
	bounds := make([]Rect, 0)
	for _, b := range scene.GetBoundaries() {
		if b.Left.Y < area.Right.Y {
			bounds = append(bounds, b)
		}
	}
	bound := bounds[rnd.Intn(len(bounds))]
	dx := 0
	for dx == 0 {
		dx = rnd.Intn(int(bound.Right.X - bound.Left.X))
	}
	grect := gopher.GetCoords()
	dy := abs(abs(grect.Right.Y)-abs(grect.Left.Y)) / 2
	return Point{X: bound.Left.X + float32(dx), Y: bound.Left.Y + dy}
}

func abs(f float32) float32 {
	if f < 0.0 {
		return -f
	}
	return f
}
