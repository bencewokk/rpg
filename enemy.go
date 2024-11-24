package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type enemy struct {
	pos     pos
	texture *ebiten.Image

	speed float32

	offsetForAnimation int
	animationState     int

	aiState int

	inPatrol   bool
	target     pos
	route      []pos
	routeIndex int
}

func createEnemy(pos pos) {
	var e enemy
	e.pos = pos
	e.speed = 100
	drawables = append(drawables, &e)
}

func (e *enemy) todoEnemy() {
	e.updateAnimation()
	e.updateAiState()

	nearestP, _ := findClosestPointOnPaths(e.pos)

	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(nearestP.float_x)), float64(offsetsy(nearestP.float_y)), 8, uilightgray2)
}

func (e *enemy) updateAnimation() {
	e.texture = enemyAnimations[0][(animationCycle+e.offsetForAnimation)%6]
}

func (e *enemy) moveTowards(target pos) {
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
	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(e.pos.float_x)), float64(offsetsy(e.pos.float_y)), 20, uilightred)
	ebitenutil.DrawCircle(screenGlobal, float64(offsetsx(c.pos.float_x)), float64(offsetsy(c.pos.float_y)), 20, uilightred)

	if !e.inPatrol {
		e.inPatrol = true
		//TODO: add calc route with A*

		fmt.Println(e.route)
		e.target = e.route[e.routeIndex]

	} else {
		e.moveTowards(e.target)
	}

	ebitenutil.DebugPrintAt(screenGlobal, strconv.Itoa(int(Distance(e.target, e.pos))), 0, 50)

	if Distance(e.target, e.pos) < 10 {

		fmt.Println(e.routeIndex, len(e.route)-1)

		if e.routeIndex == len(e.route)-1 {
			e.inPatrol = false
			e.routeIndex = 0
		} else {
			e.routeIndex++

			e.target = e.route[e.routeIndex]
		}

	}
}

func (e *enemy) updateAiState() {
	nearestP, distanceToNearest := findClosestPointOnPaths(e.pos)
	switch e.aiState {
	case 0: // roaming
		if distanceToNearest > 40 {
			e.moveTowards(nearestP)
		} else {
			e.aiState = 1
		}
	case 1:
		e.patrol()
	}
}
