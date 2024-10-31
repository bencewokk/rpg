package main

//
import (
	"github.com/hajimehoshi/ebiten/v2"
)

// 0 up, 1 down, 2 right, 3 left
func checkNextTile(way int) bool {

	//top left corner
	topleftx, toplefty := char.pos.float_x, char.pos.float_y
	topleftpos := createPos(topleftx, toplefty)

	//top right corner
	toprightx, toprighty := char.pos.float_x+screendivisor, char.pos.float_y
	toprightpos := createPos(toprightx, toprighty)

	// //bottom left corner
	bottomleftx, bottomlefty := char.pos.float_x, char.pos.float_y+screendivisor
	bottomleftpos := createPos(bottomleftx, bottomlefty)

	// //bottom right corner
	bottomrightx, bottomrighty := char.pos.float_x+screendivisor, char.pos.float_y+screendivisor
	bottomrightpos := createPos(bottomrightx, bottomrighty)

	var x, y int

	switch way {
	//up
	case 0:
		topleftpos.float_y -= speed
		toprightpos.float_y -= speed

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
		bottomrightpos.float_y += speed
		bottomleftpos.float_y += speed

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
		bottomrightpos.float_x += speed
		toprightpos.float_x += speed

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
		topleftpos.float_x -= speed
		bottomleftpos.float_x -= speed

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

func checkmovement() {
	ptid(char.pos)

	if ebiten.IsKeyPressed(ebiten.KeyD) && checkNextTile(2) {
		char.pos.float_x += speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) && checkNextTile(3) {
		char.pos.float_x -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) && checkNextTile(0) {
		char.pos.float_y -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) && checkNextTile(1) {
		char.pos.float_y += speed
	}

}
