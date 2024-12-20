package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARSPEED   = 200
	DASHSPEED   = 700
	BOOSTSPEED  = 300
	ATTACKSPEED = 100
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

	hp float32

	attacking                bool
	sinceAttack              float64
	attackCooldown           float64
	offsetForAnimationAttack int
}

func createCharacter() {
	var c character

	c.hp = 100
	c.pos = createPos(screenWidth/2, screenHeight/2)
	c.speed = CHARSPEED

	c.offsetForAnimation = rand.IntN(5)

	drawables = append(drawables, &c)
	game.currentmap.players = append(game.currentmap.players, &c)
}

func nearestCharacter(pos pos) *character {
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

func charactersInRange(pos pos, distance float32) []*character {
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
	if !c.attacking {
		if c.running {
			c.animationState = 2
		} else {
			c.animationState = 0
		}

		c.texture = characterAnimations[c.animationState+c.facingNorth][(animationCycle+c.offsetForAnimation)%6]
		c.running = false
	} else {
		c.speed = ATTACKSPEED
		c.animationState = 4
		c.texture = characterAnimations[c.animationState+c.facingNorth][(animationCycle+c.offsetForAnimation-c.offsetForAnimationAttack)%4]
		if c.sinceAttack < 0 {
			c.attacking = false
			c.speed = CHARSPEED
			c.attackCooldown = 0.5
		}
	}
}

func PushAway(enemy *enemy, character *character, pushStrength float32) pos {
	// Calculate direction vector
	dx := enemy.pos.float_x - character.pos.float_x
	dy := enemy.pos.float_y - character.pos.float_y

	// Calculate the magnitude (distance)
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Normalize the direction vector and avoid division by zero
	if distance != 0 {
		dx /= distance
		dy /= distance
	}

	// Push enemy away by the normalized direction scaled by pushStrength
	enemy.pos.float_x += dx * pushStrength
	enemy.pos.float_y += dy * pushStrength

	return enemy.pos
}

func (c *character) attack() {
	c.attacking = true
	c.sinceAttack = 0.52
	c.offsetForAnimationAttack = animationCycle

	es := enemiesInRange(c.pos, 80)

	for i := 0; i < len(es); i++ {
		es[i].hp -= float32(5 / len(es) * 4)
		es[i].hit = true
		es[i].sinceHit = 0.2
		es[i].pos = PushAway(es[i], c, 30)
	}

}

func (c *character) checkHp() {
	if c.hp < 1 {
		removeAtID(c.id, drawables)
	}
}

func (c *character) todoCharacter() {
	c.checkHp()
	c.updateCamera()
	c.checkMovement()
	c.updateAnimation()
}
