package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// pos takes in the upper left corners pos
//
// returns an image representing what should be  drawn on display after cut
func cutCam(s *ebiten.Image, pos pos) *ebiten.Image {
	return s.SubImage(image.Rect(int(pos.float_x), int(pos.float_y), int(pos.float_x)+width, int(pos.float_y)+height)).(*ebiten.Image)
}
