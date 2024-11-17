package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamemap struct {
	// map data (2D array)
	//
	// 0 = not decided, 1 = mountains, 2 = plains, 3 = dry
	data    [100][150]int
	texture [100][150]*ebiten.Image

	// height of the map
	//
	//used for rendering and generating the map
	height int
	width  int

	sprites []sprite
}

type drawable struct {
	typeOf int
	pos    pos
}

// read more in gamestate
type camera struct {
	pos pos

	//used in rendering and collision checking
	zoom float32
}

type sprite struct {
	pos     pos
	texture *ebiten.Image
}

func offsetsx(tobeoffset float32) float32 {
	return ((tobeoffset+globalGameState.camera.pos.float_x)*globalGameState.camera.zoom + screenWidth/2)
}
func offsetsy(tobeoffset float32) float32 {
	return ((tobeoffset+globalGameState.camera.pos.float_y)*globalGameState.camera.zoom + screenHeight/2)

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

	// contains the camera positions
	//
	// this is used in the rendering, it offsets the drawing positions
	camera camera
}

func updateCamera() {
	globalGameState.camera.pos.float_x = char.pos.float_x * -1
	globalGameState.camera.pos.float_y = char.pos.float_y * -1
}
