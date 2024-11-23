package main

import (
	"fmt"
	"math"

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

	nearestP, _ := findClosestPointOnPaths(e.pos, game.currentmap.paths)

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

func (e *enemy) patrol(nearestP pos) {
	var target pos
	if Distance(target, e.pos) < 0 {
		n, _ := findNearestNode(nearestP)
		target = randomPointWithinRange(*n, 10)
	} else {
		fmt.Println(target)
		fmt.Println(e.pos)
		e.moveTowards(target)
	}

	e.moveTowards(target)
}

func (e *enemy) updateAiState() {
	nearestP, distanceToNearest := findClosestPointOnPaths(e.pos, game.currentmap.paths)
	switch e.aiState {
	case 0: // roaming
		if distanceToNearest > 40 {
			e.moveTowards(nearestP)
		} else {
			e.patrol(nearestP)
		}
	}
}
