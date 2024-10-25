package main

import (
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamemap struct {
	// map data (2D array)
	//
	// 0 = not decided, 1 = mountains, 2 = plains, 3 = hills, 4 = forests
	data [72][128]int

	// height of the map
	//
	//used for rendering and generating the map
	height int
	width  int
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
	//  this is the current map that is  being used//while rendered map array size is constant to 144 (16*9) currentmapid is not
	currentmap gamemap
}

// Global variable for player
var char character = createCharacter("character.png", "character")

// Contains all information about the character
type character struct {
	title   string
	pos     pos
	picture *ebiten.Image
}

// Returns a character with the given title and path to the picture
func createCharacter(path, title string) character {
	var c character
	c.title = title

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
	c.picture.Bounds()
	// Draw the character's image onto the screen
	op.GeoM.Translate(float64(c.pos.float_x), float64(c.pos.float_y))

	// Original image size
	originalWidth, originalHeight := c.picture.Size()

	// Calculate scaling factors
	scaleX := float64(screendivisor) / float64(originalWidth) * 1.1
	scaleY := float64(screendivisor) / float64(originalHeight) * 1.1

	// Apply scaling factors
	op.GeoM.Scale(scaleX, scaleY)

	// Draw the character's image onto the screen
	screen.DrawImage(c.picture, op)
}
