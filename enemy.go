package main

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	hp int
}

func createEnemy(pos pos) {
	var e enemy
	e.pos = pos
	e.speed = ENEMYNORMALSPEED
	e.hp = 60
	drawables = append(drawables, &e)

}

func (e *enemy) todoEnemy() {
	e.updateAnimation()
	e.updateState()
	e.checkHp()

	nearestP, _ := findClosestPointOnPaths(e.pos)

	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(nearestP.float_x)), float64(offsetsy(nearestP.float_y)), 8, uilightgray2)
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

	c := findClosestNode(e.pos)
	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(e.target.float_x)), float64(offsetsy(e.target.float_y)), 20, uilightred)

	if len(e.route) != 0 {
		ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(e.route[len(e.route)-1].float_x)), float64(offsetsy(e.route[len(e.route)-1].float_y)), 20, uidarkred)

	}

	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(e.pos.float_x)), float64(offsetsy(e.pos.float_y)), 20, uilightred)
	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(c.pos.float_x)), float64(offsetsy(c.pos.float_y)), 20, uilightred)

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

	ebitenutil.DebugPrintAt(screenGlobal, strconv.Itoa(int(Distance(e.target, e.pos))), 0, 50)

	if Distance(e.target, e.pos) < 10 {
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
	nearestP := nearestPlayer(e.pos)
	e.moveTowards(nearestP.pos)
}

func (e *enemy) checkCollision(posA, posB pos) bool {
	if Distance(posA, posB) < 50 {
		return true
	}
	return false
}

func (e *enemy) hurt(c *character) {
	e.hp -= 30
	c.hp -= 30
}

func (e *enemy) checkHp() {
	if e.hp < 1 {
		removeAtID(e.id, drawables)
	}
}

func (e *enemy) updateState() {
	if Distance(e.pos, nearestPlayer(e.pos).pos) > 100 && !e.chasing {
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

		if Distance(e.pos, nearestPlayer(e.pos).pos) > 500 {
			e.chasing = false
		}
	}
}
