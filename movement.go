package main

//
import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 0 up, 1 down, 2 right, 3 left
func checkNextTile(way int) bool {

	//camera := globalGameState.camera

	charx := (char.pos.float_x)
	chary := (char.pos.float_y)

	//top left corner
	topleftx, toplefty := charx, chary
	topleftpos := createPos(topleftx, toplefty)

	//top right corner
	toprightx, toprighty := (charx + screendivisor), chary
	toprightpos := createPos(toprightx, toprighty)

	// //bottom left corner
	bottomleftx, bottomlefty := charx, (chary + screendivisor)
	bottomleftpos := createPos(bottomleftx, bottomlefty)

	// //bottom right corner
	bottomrightx, bottomrighty := (charx + screendivisor), (chary + screendivisor)
	bottomrightpos := createPos(bottomrightx, bottomrighty)

	var x, y int
	switch way {
	//up
	case 0:
		topleftpos.float_y -= 3
		toprightpos.float_y -= 3

		x, y = ptid(topleftpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(toprightpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		return true
	//down
	case 1:
		bottomrightpos.float_y += 3
		bottomleftpos.float_y += 3

		x, y = ptid(bottomleftpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(bottomrightpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		return true

	//right
	case 2:
		bottomrightpos.float_x += 3
		toprightpos.float_x += 3

		x, y = ptid(bottomrightpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(toprightpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		return true

	//left
	case 3:
		topleftpos.float_x -= 3
		bottomleftpos.float_x -= 3

		x, y = ptid(topleftpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(bottomleftpos)
		if globalGameState.currentmap.data[y][x] == 1 {
			return false
		}

		return true
	}
	return false
}

// 0 up, 1 down, 2 right, 3 left
func checkCollision(first, second pos) bool {

	// Calculate the bounding box for `first`
	firstMinX, firstMinY := first.float_x, first.float_y
	firstMaxX, firstMaxY := first.float_x+screendivisor, first.float_y+screendivisor

	// Calculate the bounding box for `second`
	secondMinX, secondMinY := second.float_x, second.float_y
	secondMaxX, secondMaxY := second.float_x+screendivisor, second.float_y+screendivisor

	// Check for overlap on both axes
	if firstMaxX > secondMinX && firstMinX < secondMaxX &&
		firstMaxY > secondMinY && firstMinY < secondMaxY {
		return true // Collision detected
	}

	return false // No collision
}

func checkZoom() {
	_, my := ebiten.Wheel()

	if my < 0 {
		for i := 0; i < 4 && globalGameState.camera.zoom > 0.5; i++ {
			time.Sleep(4 * time.Millisecond)
			globalGameState.camera.zoom -= 0.02
		}
	} else if my > 0 {
		for i := 0; i < 4 && globalGameState.camera.zoom < 2.5; i++ {
			time.Sleep(4 * time.Millisecond)
			globalGameState.camera.zoom += 0.02
		}
	}
}

var cursor pos

func checkMovementAndInput() {

	go checkZoom()

	intmx, intmy := ebiten.CursorPosition()
	cursor.float_x, cursor.float_y = (float32(intmx)+globalGameState.camera.pos.float_x)*globalGameState.camera.zoom+screenWidth/2,
		(float32(intmy)+globalGameState.camera.pos.float_y)*globalGameState.camera.zoom+screenHeight/2
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		globalGameState.currentmap.parseTexture(cursor)
	}

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
	// 	parseTexture(globalGameState.currentmap, cursor)
	// }

	// Handle movement based on key presses and check next tile for collisions
	if ebiten.IsKeyPressed(ebiten.KeyD) && checkNextTile(2) { // Move right
		char.pos.float_x += char.speed * float32(globalGameState.deltatime)
		char.running = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) && checkNextTile(3) { // Move left
		char.pos.float_x -= char.speed * float32(globalGameState.deltatime)
		char.running = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) && checkNextTile(0) { // Move up
		char.pos.float_y -= char.speed * float32(globalGameState.deltatime)
		char.running = true
		char.facingFront = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) && checkNextTile(1) { // Move down
		char.pos.float_y += char.speed * float32(globalGameState.deltatime)
		char.running = true
		char.facingFront = true
	}

	// Handle dash timing and cooldown
	if char.dashing {
		elapsed := time.Since(char.dashStart)
		if elapsed > time.Duration(char.dashDuration)*time.Millisecond {
			char.dashing = false
			char.speed = 200           // Reset speed after dash ends
			char.lastDash = time.Now() // Record end time for cooldown tracking
		}

		if elapsed < time.Duration(char.dashDuration)*time.Millisecond/2 {
			globalGameState.camera.zoom += 0.003
		} else {
			globalGameState.camera.zoom -= 0.003
		}
	}

	// Check if dash key is pressed and dash is not already active
	if ebiten.IsKeyPressed(ebiten.KeyShift) && !char.dashing {
		char.Dash()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		//fmt.Println("attacking")
	}

	// Check for collisions with enemies
	for i := range enemies {
		if checkCollision(char.pos, enemies[i].pos) {
			char.Hurt(enemies[i].pos)
			enemies[i].Hurt(char.pos)
		}
	}

}
