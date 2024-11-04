package main

//
import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	lastTwoWays [2]int
)

func checkMovement() {
	// Log current character position (for debugging)
	ptid(char.pos)

	// Handle movement based on key presses and check next tile for collisions
	if ebiten.IsKeyPressed(ebiten.KeyD) && checkNextTile(2) { // Move right
		char.pos.float_x += char.speed * float32(globalGameState.deltatime)
		updateLastTwoWays(2)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) && checkNextTile(3) { // Move left
		char.pos.float_x -= char.speed * float32(globalGameState.deltatime)
		updateLastTwoWays(3)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) && checkNextTile(0) { // Move up
		char.pos.float_y -= char.speed * float32(globalGameState.deltatime)
		updateLastTwoWays(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) && checkNextTile(1) { // Move down
		char.pos.float_y += char.speed * float32(globalGameState.deltatime)
		updateLastTwoWays(1)
	}

	// Handle dash timing and cooldown
	if char.dashing {
		elapsed := time.Since(char.dashStart)
		if elapsed < time.Duration(char.dashDuration)*time.Millisecond {
			// Dash is ongoing
			fmt.Println("Time since dash started:", elapsed)
		} else {
			// Reset dashing state and initiate cooldown
			char.dashing = false
			char.speed = 250           // Reset speed after dash ends
			char.lastDash = time.Now() // Record end time for cooldown tracking
			fmt.Println("Dash ended, cooldown started")
		}
	}

	// Check if dash key is pressed and dash is not already active
	if ebiten.IsKeyPressed(ebiten.KeyShift) && !char.dashing {
		char.Dash()
	}

	// Check for collisions with enemies
	for i := range enemies {
		if checkCollision(char.pos, enemies[i].pos) {
			char.Hurt(enemies[i].pos)
		}
	}

}

// updateLastTwoWays updates the last two directions the character moved
func updateLastTwoWays(direction int) {
	lastTwoWays[1] = lastTwoWays[0]
	lastTwoWays[0] = direction
}
