package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type bar struct {
	pos      pos
	width    float32
	height   float32
	color    color.Color
	maxValue float32
}

// Creates and returns a new bar with the specified properties
func newBar(x, y, width, height float32, color color.Color, maxValue float32) bar {
	return bar{
		pos:      createPos(x, y),
		width:    width,
		height:   height,
		color:    color,
		maxValue: maxValue,
	}
}

var secondaryBar float32 = 600

func fillBar() {
	for secondaryBar < 600 {
		secondaryBar += 1.5
		time.Sleep(5 * time.Millisecond)

	}

	secondaryBar = 600
	char.barFilling = false
}

func drawUi(s *ebiten.Image) {

	if !char.dashing && !char.barFilling {
		char.barFilling = true
		go fillBar()
	}

	if char.dashing {
		secondaryBar = 0
	}
	vector.DrawFilledRect(s, 60, screenHeight-30, 600, 15, uidarkred, false)
	vector.DrawFilledRect(s, 60, screenHeight-30, secondaryBar, 15, uilightred, false)
}
