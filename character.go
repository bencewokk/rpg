package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARSPEED  = 200
	DASHSPEED  = 600
	BOOSTSPEED = 300
)

type character struct {
	id      int
	pos     pos
	texture *ebiten.Image

	speed              float32
	offsetForAnimation int
	animationState     int

	running     bool
	facingNorth int // between 0 and 1

	untilNewDash   float64
	untilEndOfDash float64
	dashing        bool

	untilEndOfBoost float64
}

func createCharacter() {
	var c character

	c.pos = createPos(screenWidth/2, screenHeight/2)
	c.speed = CHARSPEED

	c.offsetForAnimation = rand.IntN(5)

	drawables = append(drawables, &c)
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
	c.checkMovement()
	c.updateCamera()
	c.updateAnimation()
}
