package main

import (
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Global variable for player
var char character = createCharacter("character.png", "character")

// Contains all information about the character
type character struct {
	title   string
	pos     pos
	picture *ebiten.Image
	hp      int

	dashing      bool
	dashStart    time.Time
	speed        float32
	dashDuration float32
	lastDash     time.Time
	dashCooldown time.Duration
}

// Returns a character with the given title and path to the picture
func createCharacter(path, title string) character {
	var c character
	c.title = title
	c.hp = 1000
	c.speed = 250
	c.dashDuration = 200
	c.dashCooldown = 1

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
	c.picture = ebiten.NewImageFromImage(imgData)

	return c
}

// DrawCharacter draws the character centered on the screen
func (c character) DrawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	originalWidth, originalHeight := c.picture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth) * float64(globalGameState.camera.zoom)
	scaleY := float64(screendivisor) / float64(originalHeight) * float64(globalGameState.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Calculate the centered position based on screen dimensions and zoom level
	centerX := (float64(screenWidth) / 2) - (float64(originalWidth) * scaleX / 2)
	centerY := (float64(screenHeight) / 2) - (float64(originalHeight) * scaleY / 2)
	op.GeoM.Translate(centerX, centerY)

	screen.DrawImage(c.picture, op)
}

func (c *character) Die() {
	char.pos.float_y = screenHeight / 2
	char.pos.float_x = screenWidth / 2

	char.hp = 1000
}

// Assuming 'enemy' is of type 'character' or has a position that you can access
func (c *character) Hurt(enemyPos pos) {
	c.hp -= 10
	if c.hp <= 0 {
		c.Die()
	}

	// Calculate the direction to move away from the enemy
	moveAmount := float32(30) // Amount to move away
	directionX := c.pos.float_x - enemyPos.float_x
	directionY := c.pos.float_y - enemyPos.float_y

	// Normalize the direction vector
	length := float32(math.Sqrt(float64(directionX*directionX+directionY*directionY))) * 2.5
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	// Move the character away from the enemy
	c.pos.float_x += directionX * moveAmount
	c.pos.float_y += directionY * moveAmount
}

func (c *character) Dash() {
	// Check if enough time has passed since the last dash to allow dashing again (cooldown)
	if time.Since(c.lastDash) < time.Duration(c.dashCooldown)*time.Second {
		return
	}

	// Start dash if it's not already active
	if !c.dashing {
		c.dashing = true
		c.dashStart = time.Now()
		c.speed = 800 // Increase speed for dash
	}
}
