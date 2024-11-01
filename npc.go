package main

import (
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Contains all information about the enemies
type enemy struct {
	title       string
	pos         pos
	picture     *ebiten.Image
	curtiletype int
}

var (
	enemies []enemy
)

// Returns a new enemy with the given title and path to the picture
func createEnemy(path, title string) enemy {
	var e enemy
	e.title = title

	// Open the image file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the image file into an image.Image
	imgData, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the image.Image to an *ebiten.Image
	e.picture = ebiten.NewImageFromImage(imgData)

	return e
}

// Draw the enemies
func (e enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	originalWidth, originalHeight := e.picture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth)
	scaleY := float64(screendivisor) / float64(originalHeight)
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(float64(e.pos.float_x), float64(e.pos.float_y))

	screen.DrawImage(e.picture, op)
}

func init() {
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 1"))
	enemies[0].pos = createPos(60, 60)
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 2"))
	enemies[1].pos = createPos(120, 90)
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 3"))
	enemies[2].pos = createPos(60, 270)
}
