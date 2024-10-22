package main

/*
DISCRIPTION

	This is a basic RPG game made with Ebiten.

LOGS

	2024.10.19 ver 0.0.1
	- This is the beginning of the game. The main game loop is set up,
	  but no game logic has been implemented yet.
	- Added button ui functions and basic game variables such as gamestate.
	- Started the menu with options and play buttons

	2024.10.20 ver 0.0.2
	- Moved screenWidth annd screenHeight out of gamestate
	- Added basic testmap

	2024.10.21 ver 0.0.3
	- Added a function create maps
	- Added support on all screen sizes
	- Started adding better mapgeneration


TODO
	important:
		- add a test map  / done
		- add a test character
		- decide whether to use rendered map or not / done
		- better map creation
	unimportant:
		- split up code into individual files
		- move readme to readme
*/

/*
for me

ideas: different stats on different tile types
*/
import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func gameinit() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

}

// Color variable for ui
var (
	uigray        = color.RGBA{128, 128, 128, 255}
	uidarkgray    = color.RGBA{100, 100, 100, 255}
	uilightgray   = color.RGBA{158, 158, 158, 255}
	uitransparent = color.RGBA{0, 0, 0, 0}
)

// Color variable for map rendering
var (
	mlightgreen      = color.RGBA{144, 238, 144, 255}
	mbrown           = color.RGBA{139, 69, 19, 255}
	mdarkgray        = color.RGBA{169, 169, 169, 255}
	mdarkgreen       = color.RGBA{34, 139, 34, 255}
	currenttilecolor = color.RGBA{0, 0, 0, 0}
)

// Screen sizes
var (
	width, height    = ebiten.Monitor().Size()
	screenWidth      = float32(width)
	screenHeight     = float32(height)
	screendivisor    = float32(120)
	intscreendivisor = screenHeight / 9
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

var testmap = gamemap{
	data: [9][16]int{
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		{2, 3, 0, 0, 3, 3, 3, 3, 1, 1, 0, 2, 0, 1, 0, 2},
		{2, 3, 2, 3, 1, 3, 1, 2, 0, 0, 3, 3, 1, 2, 1, 2},
		{2, 1, 3, 3, 2, 0, 3, 0, 3, 1, 3, 3, 2, 3, 2, 2},
		{2, 1, 3, 3, 0, 2, 1, 1, 1, 2, 0, 3, 0, 0, 0, 2},
		{2, 3, 1, 0, 1, 1, 0, 1, 3, 3, 2, 1, 1, 3, 3, 2},
		{2, 2, 2, 3, 1, 0, 2, 1, 3, 3, 2, 1, 2, 3, 2, 2},
		{2, 1, 1, 3, 3, 1, 3, 2, 0, 0, 2, 2, 0, 0, 1, 2},
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
	},
}

var testmap2 = createMap()

type gamemap struct {
	// map data (2D array)
	//
	//0 plains, 1 hills, 2 mountains, 3 forests
	data [9][16]int
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
	currentmapid int

	// character information

	// contains where the character is on the array
	//
	//this translates to gamemap[y][x]
	currentcharposx, currentcharposy int
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

		playbtn := createButton("Play", 25, 25, 150, 50, uitransparent, uilightgray, uigray)
		playbtn.DrawButton(screen)
		if playbtn.pressed {
			globalGameState.stateid = 3
		}

		optionsbtn := createButton("Options", 25, 85, 150, 50, uitransparent, uilightgray, uigray)
		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			globalGameState.stateid = 1
		}
	case 1:
		options_exitbtn := createButton("Back to menu", 25, 25, 150, 50, uitransparent, uilightgray, uigray)
		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			globalGameState.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
	case 3:
		switch globalGameState.currentmapid {
		case 0:
			ebitenutil.DebugPrint(screen, "Test map")
			// this is 32 because we are rendering that part of the map
			// one tile is 240x240 px so this would fill the whole screen
			for i := 0; i < 9; i++ {
				for j := 0; j < 16; j++ {
					switch testmap2.data[i][j] {
					case 0:
						currenttilecolor = mlightgreen
					case 1:
						currenttilecolor = mbrown
					case 2:
						currenttilecolor = mdarkgray
					case 3:
						currenttilecolor = mdarkgreen
					}
					vector.DrawFilledRect(
						screen,
						float32(j*int(intscreendivisor)),
						float32(i*int(intscreendivisor)),
						screendivisor,
						screendivisor,
						currenttilecolor,
						false,
					)
				}
			}
		}
	}

}

// Layout method of the Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	gameinit()
	fmt.Println(createMap())
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
