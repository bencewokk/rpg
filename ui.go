package main

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type damage struct {
	pos pos

	wanterFor  float64
	sinceDrawn float64
	damage     int
}

func (c *character) drawUi() {
	vector.DrawFilledRect(screenGlobal, 25, 50, 15, 200, uilightred, false)
	vector.DrawFilledRect(screenGlobal, 25, 50, 15, c.hp*2, uidarkred, false)

	vector.DrawFilledRect(screenGlobal, 25, 260, 15, 300, mlightgreen, false)
	if c.untilNewDash > 0 {
		vector.DrawFilledRect(screenGlobal, 25, 260, 15, float32(c.untilNewDash)*300/1.5, mdarkgreen, false)
	}
}
