package main

import (
	"fmt"
	"image/png"
	"log"
	"math"
	"os"

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
	dashing bool
}

// Returns a character with the given title and path to the picture
func createCharacter(path, title string) character {
	var c character
	c.title = title
	c.hp = 100

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

// DrawCharacter draws the character
func (c character) DrawCharacter(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	originalWidth, originalHeight := c.picture.Size()
	scaleX := float64(screendivisor) / float64(originalWidth)
	scaleY := float64(screendivisor) / float64(originalHeight)
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(float64(c.pos.float_x), float64(c.pos.float_y))

	screen.DrawImage(c.picture, op)
}

// Assuming 'enemy' is of type 'character' or has a position that you can access
func (c *character) Hurt(enemyPos pos) {
	c.hp -= 10
	if c.hp <= 0 {
		// Character is dead
		fmt.Println("Character is dead")
	}

	fmt.Println(lastTwoWays[0])

	// Calculate the direction to move away from the enemy
	moveAmount := float32(10) // Amount to move away
	directionX := c.pos.float_x - enemyPos.float_x
	directionY := c.pos.float_y - enemyPos.float_y

	// Normalize the direction vector
	length := float32(math.Sqrt(float64(directionX*directionX + directionY*directionY)))
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	// Move the character away from the enemy
	c.pos.float_x += directionX * moveAmount
	c.pos.float_y += directionY * moveAmount

	if lastTwoWays[0] != lastTwoWays[1] {
		// You can handle any additional logic here if necessary
	}
}

func (c *character) Dash(screen *ebiten.Image) {
	//var dashSpeed float64 = 10.0       // Speed during dash
	var dashDuration float64 = 0.2     // Time in seconds for the dash
	var dashCooldown float64 = 1.0     // Time in seconds between dashes
	var isDashing bool = false         // Whether the player is currently dashing
	var dashTimeLeft float64 = 0.0     // Time remaining in the current dash
	var dashCooldownTime float64 = 3.0 // Time left for cooldown

	if isDashing {
		// Apply dash movement
		fmt.Println("test")
		dashTimeLeft -= globalGameState.deltatime

		// End dash if time is up
		if dashTimeLeft <= 0 {
			isDashing = false
			dashCooldownTime = dashCooldown
		}
	} else if dashCooldownTime > 0 {
		dashCooldownTime -= globalGameState.deltatime
	} else if ebiten.IsKeyPressed(ebiten.KeyB) {
		fmt.Println("TEST3")
		isDashing = true
		dashTimeLeft = dashDuration
	}

}
