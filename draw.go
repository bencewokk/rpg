package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

var drawables []drawable

type drawable interface {
	draw(sceen *ebiten.Image)
	Y() float32
}

func (c *character) getId(id int) {
	c.id = id
}
func (e *enemy) getId(id int) {
	e.id = id
}
func (t *tree) getId(id int) {
	t.id = id
}

func removeAtID(id int, d []drawable) []drawable {
	var dr []drawable
	dr = append(d[:id], d[id+1:]...)
	return dr
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

func (e *enemy) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	e.todoEnemy()

	scaleX := float64(screendivisor) / 18 * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / 18 * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Positioning with respect to camera
	op.GeoM.Translate(
		float64(offsetsx(e.pos.float_x)),
		float64(offsetsy(e.pos.float_y)),
	)

	// Draw the selected portion of the image onto the screen
	screen.DrawImage(e.texture, op)
}

func (t *tree) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(float64(game.camera.zoom)*1.7, float64(game.camera.zoom)*1.7)

	op.GeoM.Translate(
		float64(offsetsx(t.pos.float_x)),
		float64(offsetsy(t.pos.float_y)))
	screen.DrawImage(t.texture, op)

}

func (c *character) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	c.todoCharacter()

	originalWidth, originalHeight := c.texture.Size()
	scaleX := float64(screendivisor) / 18 * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / 18 * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Calculate the centered position based on screen dimensions and zoom level
	centerX := (float64(screenWidth) / 2) - (float64(originalWidth) * scaleX / 2)
	centerY := (float64(screenHeight) / 2) - (float64(originalHeight) * scaleY / 2)
	op.GeoM.Translate(centerX, centerY)

	screen.DrawImage(c.texture, op)
}

func (c *character) Y() float32 {
	return c.pos.float_y
}

func (t *tree) Y() float32 {
	switch t.treeId {
	case 0:
		return t.pos.float_y + 95
	case 1:
		return t.pos.float_y + 102
	}
	return t.pos.float_y
}

func (e *enemy) Y() float32 {
	return e.pos.float_y + 25
}

func sortDrawables() {
	sort.Slice(drawables, func(a, b int) bool {
		return drawables[a].Y() < drawables[b].Y()
	})
}
