package main

//
import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 0 up, 1 down, 2 right, 3 left
func (c *character) checkNextTile(way int) bool {

	//camera := game.camera

	charx := (c.pos.float_x)
	chary := (c.pos.float_y)

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
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(toprightpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		return true
	//down
	case 1:
		bottomrightpos.float_y += 3
		bottomleftpos.float_y += 3

		x, y = ptid(bottomleftpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(bottomrightpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		return true

	//right
	case 2:
		bottomrightpos.float_x += 3
		toprightpos.float_x += 3

		x, y = ptid(bottomrightpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(toprightpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		return true

	//left
	case 3:
		topleftpos.float_x -= 3
		bottomleftpos.float_x -= 3

		x, y = ptid(topleftpos)
		if game.currentmap.data[y][x] == 1 {
			return false
		}

		x, y = ptid(bottomleftpos)
		if game.currentmap.data[y][x] == 1 {
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
		for i := 0; i < 4 && game.camera.zoom > 0.5; i++ {
			time.Sleep(4 * time.Millisecond)
			game.camera.zoom -= 0.02
		}
	} else if my > 0 {
		for i := 0; i < 4 && game.camera.zoom < 2.5; i++ {
			time.Sleep(4 * time.Millisecond)
			game.camera.zoom += 0.02
		}
	}
}

func (c *character) checkMovement() {

	// Handle movement based on key presses and check next tile for collisions
	if ebiten.IsKeyPressed(ebiten.KeyD) && c.checkNextTile(2) { // Move right
		c.pos.float_x += c.speed * float32(game.deltatime)
		c.running = true

	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && c.checkNextTile(3) { // Move left
		c.pos.float_x -= c.speed * float32(game.deltatime)
		c.running = true

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) && c.checkNextTile(0) { // Move up
		c.pos.float_y -= c.speed * float32(game.deltatime)
		c.running = true
		c.facingNorth = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && c.checkNextTile(1) { // Move down
		c.pos.float_y += c.speed * float32(game.deltatime)
		c.running = true
		c.facingNorth = 0
	}

}
