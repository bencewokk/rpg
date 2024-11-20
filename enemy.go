package main

import "github.com/hajimehoshi/ebiten/v2"

type enemy struct {
	pos     pos
	texture *ebiten.Image

	offsetForAnimation int
	animationState     int
}

func createEnemy(pos pos) {
	var e enemy
	e.pos = pos
	drawables = append(drawables, &e)
}

func (e *enemy) todoEnemy() {
	e.updateAnimation()
}

func (e *enemy) updateAnimation() {
	e.texture = enemyAnimations[0][(animationCycle+e.offsetForAnimation)%6]
}
