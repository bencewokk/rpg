package main

import (
	"time"
)

type gamemap struct {
	// map data (2D array)
	//
	// 0 = not decided, 1 = mountains, 2 = plains, 3 = hills, 4 = forests
	data [36][64]int

	// height of the map
	//
	//used for rendering and generating the map
	height int
	width  int
}

// read more in gamestate
type camera struct {
	pos pos

	//used in rendering and collision checking
	zoom float32
}

// contains all global information about the game
var globalGameState gamestate

// contains all global information about the game
//
// contains maps
type gamestate struct {
	// 0 menu / 1 menu and options / 3 in game
	stateid int

	// maps are stored in arrays (see in type map)
	//
	// this is the current map that is  being used//while rendered map array size is constant to 144 (16*9) currentmapid is not
	currentmap gamemap

	// counts the time since start of game
	//
	// get updated every frame
	deltatime float64

	// date of last update
	lastUpdateTime time.Time

	// contains the camere positions
	//
	// this is used in the rendering, it offsets the drawing positions
	camera camera
}

func updateCamera() {
	globalGameState.camera.pos.float_x = char.pos.float_x * -1
	globalGameState.camera.pos.float_y = char.pos.float_y * -1
}
