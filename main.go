package main

/*
DISCRIPTION

	This is a basic RPG game made with Ebiten.

LOGS

	2024.10.19 ver 0.0.1
	-This is the beginning of the game. The main game loop is set up, but no game logic has been implemented yet.
	-Added button ui functions and basic game variables such as gamestate.
	-Started the menu with options and play buttons

TODO
	important:
		- add a test map
		- add a test character
		- decide whether to use rendered map or not
	unimportant:
*/

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func gameinit() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")
	width, height := ebiten.Monitor().Size()
	globalGameState.screenWidth = float32(width)
	globalGameState.screenHeight = float32(height)
}

// Color variables
var (
	gray        = color.RGBA{128, 128, 128, 255}
	darkergray  = color.RGBA{100, 100, 100, 255}
	lightergray = color.RGBA{158, 158, 158, 255}
	transparent = color.RGBA{0, 0, 0, 0}
)

// Button struct definition
type button struct {
	title         string
	x, y          float32
	width, height float32
	pressed       bool
	hovered       bool
	pressedColor  color.RGBA
	hoveredColor  color.RGBA
	inactiveColor color.RGBA
}

// contains all global information about the game
var globalGameState gamestate

// contains all global information about the game
type gamestate struct {
	// 0 menu / 1 menu and options / 3 in game
	stateid int

	// screen width (x) and height (y)
	//
	// defaults are the user's screen dimensions, declared in gameinit
	//
	// constant
	screenWidth, screenHeight float32

	// maps are stored in arrays (see in type map)
	//
	// while rendered map array size is constant to 144 (16*9) currentmapid is not
	currentmapid int
}

// Create a new button
func createButton(title string, x, y, width, height float32, pressedColor, hoveredColor, inactiveColor color.RGBA) button {
	return button{
		title:         title,
		x:             x,
		y:             y,
		width:         width,
		height:        height,
		pressedColor:  pressedColor,
		hoveredColor:  hoveredColor,
		inactiveColor: inactiveColor,
	}
}

// DrawButton draws the button and checks for interaction
func (b *button) DrawButton(screen *ebiten.Image) {
	// Get the mouse position
	intmx, intmy := ebiten.CursorPosition()
	mx, my := float32(intmx), float32(intmy)

	// Check if the mouse is inside the button's area
	if mx >= b.x && mx <= b.x+b.width && my >= b.y && my <= b.y+b.height {
		b.hovered = true
		// Check if the left mouse button is pressed
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.pressed = true
		} else {
			b.pressed = false
		}
	} else {
		b.hovered = false
		b.pressed = false
	}

	// Choose color based on button state
	var drawColor color.Color
	if b.pressed {
		drawColor = b.pressedColor
	} else if b.hovered {
		drawColor = b.hoveredColor
	} else {
		drawColor = b.inactiveColor
	}

	// Draw the button rectangle
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), drawColor, false)

	// Draw the button title as text
	ebitenutil.DebugPrintAt(screen, b.title, int(b.x)+10, int(b.y)+10)
}

type Game struct{}

// Update method of the Game
func (g *Game) Update() error {
	return nil
}

// Draw method of the Game
func (g *Game) Draw(screen *ebiten.Image) {

	switch globalGameState.stateid {
	case 0:

		playbtn := createButton("Play", 25, 25, 150, 50, transparent, lightergray, gray)
		playbtn.DrawButton(screen)
		if playbtn.pressed {
			globalGameState.stateid = 3
		}

		optionsbtn := createButton("Options", 25, 85, 150, 50, transparent, lightergray, gray)
		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			globalGameState.stateid = 1
		}
	case 1:
		options_exitbtn := createButton("Back to menu", 25, 25, 150, 50, transparent, lightergray, gray)
		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			globalGameState.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, globalGameState.screenWidth-250, globalGameState.screenHeight-50, darkergray, false)
	case 3:

	}

}

// Layout method of the Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	gameinit()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
