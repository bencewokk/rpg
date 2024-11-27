package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARSPEED  = 200
	DASHSPEED  = 700
	BOOSTSPEED = 300
)

type character struct {
	pos     pos
	texture *ebiten.Image
	id      int

	speed              float32
	offsetForAnimation int
	animationState     int

	running     bool
	facingNorth int // between 0 and 1

	untilNewDash   float64
	untilEndOfDash float64
	dashing        bool

	untilEndOfBoost float64

	hp int
}

func createCharacter() {
	var c character

	c.pos = createPos(screenWidth/2, screenHeight/2)
	c.speed = CHARSPEED

	c.offsetForAnimation = rand.IntN(5)

	drawables = append(drawables, &c)
	game.currentmap.players = append(game.currentmap.players, &c)
}

func nearestPlayer(pos pos) *character {
	var closest int
	var closestDistance float32

	for i := 0; i < len(game.currentmap.players); i++ {
		if closestDistance > Distance(pos, game.currentmap.players[i].pos) {
			closestDistance = Distance(pos, game.currentmap.players[i].pos)
			closest = i
		}
	}

	return game.currentmap.players[closest]
}

func playersInRange(pos pos, distance float32) []*character {
	var cs []*character
	for i := 0; i < len(game.currentmap.players); i++ {
		if Distance(pos, game.currentmap.players[i].pos) > distance {
			cs = append(cs, game.currentmap.players[i])
		}
	}
	return cs
}

func (c *character) updateCamera() {
	game.camera.pos = c.pos
}

func (c *character) updateAnimation() {
	if c.running {
		c.animationState = 2
	} else {
		c.animationState = 0
	}
	c.texture = characterAnimations[c.animationState+c.facingNorth][(animationCycle+c.offsetForAnimation)%6]
	c.running = false
}

func (c *character) todoCharacter() {
	c.updateCamera()
	c.checkMovement()
	c.updateAnimation()

}
