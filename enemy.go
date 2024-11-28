package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ENEMYNORMALSPEED = 50
	ENEMYALLERTSPEED = 100
)

type enemy struct {
	pos     pos
	texture *ebiten.Image
	id      int

	speed float32

	offsetForAnimation int
	animationState     int

	aiState int

	inPatrol   bool
	target     pos
	route      []pos
	routeIndex int
	chasing    bool

	sleeping   bool
	sinceSleep float64

	hp       float32
	hit      bool
	sinceHit float64
}

func enemiesInRange(pos pos, distance float32) []*enemy {
	var es []*enemy
	for i := 0; i < len(game.currentmap.enemies); i++ {
		if Distance(pos, game.currentmap.enemies[i].pos) < distance {
			es = append(es, game.currentmap.enemies[i])
		}
	}
	return es
}

func createEnemy(pos pos) {
	var e enemy
	e.pos = pos
	e.speed = ENEMYNORMALSPEED
	e.hp = 60
	e.offsetForAnimation = rand.Intn(5)
	game.currentmap.enemies = append(game.currentmap.enemies, &e)
	drawables = append(drawables, &e)
}

func (e *enemy) todoEnemy() {
	e.updateAnimation()
	e.updateState()
	e.checkHp()

}

func (e *enemy) updateAnimation() {
	e.texture = enemyAnimations[e.animationState][(animationCycle+e.offsetForAnimation)%6]
	e.animationState = 0
}

func (e *enemy) moveTowards(target pos) {
	e.animationState = 1
	dx := target.float_x - e.pos.float_x
	dy := target.float_y - e.pos.float_y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length > 0 {
		dx /= length
		dy /= length

		e.pos.float_x += dx * e.speed * float32(game.deltatime)
		e.pos.float_y += dy * e.speed * float32(game.deltatime)
	}

}

func (e *enemy) patrol() {

	e.sinceSleep += game.deltatime

	if !e.inPatrol {
		e.inPatrol = true
		e.sinceSleep = 0
		e.routeIndex = 0
		e.route = findShortestPathPositions(findClosestNode(e.pos).id, findClosestNode(randomPointWithinRange(findClosestNode(e.pos), 6)).id)
		e.target = e.route[e.routeIndex]

	} else {
		e.moveTowards(e.target)
	}

	if Distance(e.target, e.pos) < 30 {
		if e.routeIndex == len(e.route)-1 {
			e.inPatrol = false
			e.sinceSleep = 0
			e.routeIndex = 0
		} else {
			e.routeIndex++
			e.target = e.route[e.routeIndex]
		}
	}
}

func (e *enemy) chase() {
	nearestP := nearestCharacter(e.pos)
	e.moveTowards(nearestP.pos)
}

func (e *enemy) checkCollision(posA, posB pos) bool {
	if Distance(posA, posB) < 50 {
		return true
	}
	return false
}

func (e *enemy) hurt(c *character) {
	e.hp -= 0.3
	c.hp -= 0.1
}

func (e *enemy) checkHp() {
	if e.hp < 1 {
		removeAtID(e.id, drawables)
	}
}

func (e *enemy) updateState() {

	e.sinceHit -= game.deltatime

	if e.sinceHit < 0 {

		e.hit = false
	}

	if Distance(e.pos, nearestCharacter(e.pos).pos) > 100 && !e.chasing {
		e.speed = ENEMYNORMALSPEED
		nearestP, distanceToNearest := findClosestPointOnPaths(e.pos)
		switch e.aiState {
		case 0: // move towards nearest path
			if distanceToNearest > 40 {
				e.moveTowards(nearestP)
			} else {
				e.aiState = 1
			}
		case 1:
			if e.sinceSleep > 1 {
				e.patrol()
			} else {
				e.sinceSleep += game.deltatime
			}
		}
	} else {

		for i := 0; i < len(game.currentmap.players); i++ {
			c := game.currentmap.players[i]
			if checkCollision(e.pos, c.pos) {
				e.hurt(c)
			}
		}

		e.chasing = true
		e.speed = ENEMYALLERTSPEED
		e.inPatrol = false
		e.chase()

		if Distance(e.pos, nearestCharacter(e.pos).pos) > 400 {
			e.chasing = false
		}
	}
}
