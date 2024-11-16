package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	NORMALSPEED = 150
	ONDRYSPEED  = 300
	DASHSPEED   = 500
)

// Global variable for player
var char character = createCharacter("character")

// Load all animations
//
// 0 FRONT IDLE
//
// 1 BACK IDLE
//
// 2 FRONT RUNNING
//
// 3 BACK RUNNING
func load() {
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle.png"))
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S2.png"))
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S3.png"))
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S4.png"))
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S5.png"))
	char.allAnimations[0] = append(char.allAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S6.png"))

	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle.png"))
	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S2.png"))
	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S3.png"))
	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S4.png"))
	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S5.png"))
	char.allAnimations[1] = append(char.allAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S6.png"))

	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running.png"))
	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S2.png"))
	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S3.png"))
	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S4.png"))
	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S5.png"))
	char.allAnimations[2] = append(char.allAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S6.png"))

	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running.png"))
	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S2.png"))
	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S3.png"))
	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S4.png"))
	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S5.png"))
	char.allAnimations[3] = append(char.allAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S6.png"))
}

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
	barFilling   bool

	currentAnimationState_primary   int // use this to index allAnimation []EbitenImage
	currentAnimationState_secondary int // use this to index allAnimation []EbitenImage

	allAnimations [6][]*ebiten.Image
	facingFront   bool
	running       bool

	lastSpeedBoostTime time.Time
}

// Returns a character with the given title and path to the picture
func createCharacter(title string) character {
	var c character
	c.title = title
	c.hp = 1000
	c.speed = NORMALSPEED
	c.dashDuration = 600
	c.dashCooldown = 2

	// Initialize the animation arrays to avoid index out of range
	// Initialize the first animation (idle)
	c.allAnimations[0] = make([]*ebiten.Image, 0)

	c.picture = loadPNG("import/Characters/Character/Front_C_Idle.png")

	return c
}

var animationTimer float64

// Updates animation on global character
func updateAnimationCharacter() {
	animationTimer += globalGameState.deltatime

	if char.dashing {
		animationTimer += 0.01
	}

	if animationTimer >= 0.15 {
		if char.facingFront {
			if char.running {
				char.currentAnimationState_secondary = 2
			}
		} else {
			if char.running {
				char.currentAnimationState_secondary = 3
			}
		}

		char.currentAnimationState_primary = (char.currentAnimationState_primary + 1) % 6
		char.picture = char.allAnimations[char.currentAnimationState_secondary][char.currentAnimationState_primary]
		animationTimer = 0.0
	}

}

// DrawCharacter draws the character centered on the screen
func (c *character) DrawCharacter(screen *ebiten.Image) {

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

	char.running = false
	switch char.facingFront {
	case true:
		char.currentAnimationState_secondary = 0
	case false:
		char.currentAnimationState_secondary = 1
	}

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
		c.speed = DASHSPEED // Increase speed for dash
	}
}
