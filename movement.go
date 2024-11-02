package main

//
import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	speed       float32 = 250
	lastTwoWays [2]int
)

func checkMovement() {
	// Log current character position (for debugging)
	ptid(char.pos)

	// Handle movement based on key presses and check next tile for collisions
	if ebiten.IsKeyPressed(ebiten.KeyD) && checkNextTile(2) { // Move right
		char.pos.float_x += speed * float32(globalGameState.deltatime)
		updateLastTwoWays(2)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) && checkNextTile(3) { // Move left
		char.pos.float_x -= speed * float32(globalGameState.deltatime)
		updateLastTwoWays(3)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) && checkNextTile(0) { // Move up
		char.pos.float_y -= speed * float32(globalGameState.deltatime)
		updateLastTwoWays(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) && checkNextTile(1) { // Move down
		char.pos.float_y += speed * float32(globalGameState.deltatime)
		updateLastTwoWays(1)
	}

	// Check for collisions with enemies
	for i := range enemies {
		if checkCollision(char.pos, enemies[i].pos) {
			fmt.Println("test2")
			char.Hurt(enemies[i].pos) // Call the Hurt method on collision
		}
	}

}

// updateLastTwoWays updates the last two directions the character moved
func updateLastTwoWays(direction int) {
	lastTwoWays[1] = lastTwoWays[0]
	lastTwoWays[0] = direction
}
