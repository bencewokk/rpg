package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

var drawables []drawable

type drawable interface {
	draw(sceen *ebiten.Image)
	pos()
}

func drawTile(screen, t *ebiten.Image, i, j int) {
	op := &ebiten.DrawImageOptions{}

	originalWidth, originalHeight := t.Size()
	scaleX := float64(screendivisor) / float64(originalWidth) * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / float64(originalHeight) * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(
		float64(offsetsx(float32(j*intscreendivisor-intscreendivisor/2))),
		float64(offsetsy(float32(i*intscreendivisor-intscreendivisor/2))))
	screen.DrawImage(t, op)
}
s
func (e *enemy) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	e.todoEnemy()

	// Set up scaling
	originalWidth, originalHeight := e.texture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth) * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / float64(originalHeight) * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Positioning with respect to camera
	op.GeoM.Translate(
		float64(offsetsx(e.pos.float_y)),
		float64(offsetsy(e.pos.float_x)),
	)

	// Draw the selected portion of the image onto the screen
	screen.DrawImage(e.texture, op)
}

type sprite struct {
	typeOf  int
	pos     pos
	texture *ebiten.Image
}

func (s *sprite) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(float64(game.camera.zoom)*1.7, float64(game.camera.zoom)*1.7)

	op.GeoM.Translate(
		float64(offsetsx(s.pos.float_x)),
		float64(offsetsy(s.pos.float_y)))
	screen.DrawImage(s.texture, op)

}

func (c *character) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	c.todoCharacter()

	originalWidth, originalHeight := c.texture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth) * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / float64(originalHeight) * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Calculate the centered position based on screen dimensions and zoom level
	centerX := (float64(screenWidth) / 2) - (float64(originalWidth) * scaleX / 2)
	centerY := (float64(screenHeight) / 2) - (float64(originalHeight) * scaleY / 2)
	op.GeoM.Translate(centerX, centerY)

	screen.DrawImage(c.texture, op)
}

func sortDrawables() {
	sort.Slice(drawables, func(a, b int) bool {
		return drawables[a].pos.float_y < drawables[b].pos.float_y
	})
}
