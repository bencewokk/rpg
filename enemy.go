package main

import (
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Contains all information about the enemies
type enemy struct {
	id          int
	title       string
	pos         pos
	picture     *ebiten.Image
	curtiletype int
	hp          int
}

var (
	enemies []enemy
)

// Returns a new enemy with the given title and path to the picture
func createEnemy(path, title string, id int) enemy {
	var e enemy
	e.title = title
	e.id = id

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
	camera := globalGameState.camera.pos

	originalWidth, originalHeight := e.picture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth)
	scaleY := float64(screendivisor) / float64(originalHeight)
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(
		float64(e.pos.float_x+screenWidth/2+camera.float_x)-float64(intscreendivisor)/2,
		float64(e.pos.float_y+screenHeight/2+camera.float_y)-float64(intscreendivisor)/2,
	)

	screen.DrawImage(e.picture, op)
}

func (e *enemy) Die() {
	e.pos.float_y = screenHeight / 2
	e.pos.float_x = screenWidth / 2

	e.hp = 100
}

// Assuming 'enemy' is of type 'character' or has a position that you can access
func (e *enemy) Hurt(enemyPos pos) {
	e.hp -= 10
	if e.hp <= 0 {
		e.Die()
	}

	// Calculate the direction to move away from the enemy
	moveAmount := float32(30) // Amount to move away
	directionX := e.pos.float_x - enemyPos.float_x
	directionY := e.pos.float_y - enemyPos.float_y

	// Normalize the direction vector
	length := float32(math.Sqrt(float64(directionX*directionX+directionY*directionY))) * 2.5
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	// Move the character away from the enemy
	e.pos.float_x += directionX * moveAmount
	e.pos.float_y += directionY * moveAmount
}

func init() {
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 1", 0))
	enemies[0].pos = createPos(60, 60)
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 2", 1))
	enemies[1].pos = createPos(120, 90)
	enemies = append(enemies, createEnemy("enemy.png", "Enemy 3", 2))
	enemies[2].pos = createPos(60, 270)
}
