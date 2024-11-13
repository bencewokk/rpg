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

var a float32

func drawUi(s *ebiten.Image) {

	if a/0.29354207436 < 300 {

		a = float32(time.Since(char.dashStart).Milliseconds())
	}
	vector.DrawFilledRect(s, 60, screenHeight-30, 300, 20, uidarkgray, false)
	vector.DrawFilledRect(s, 60, screenHeight-30, float32(a), 20, mlightgreen, false)
}
