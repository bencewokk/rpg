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

	2024.10.23 ver 0.0.4
	- Created way better map creation
		- Changable map size
		- Better randomization

	2024.10.24 ver 0.0.5
	- Added a character struct
	- Addes standardized position struct

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
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func gameinit() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	fmt.Println()

	globalGameState.currentmap = createMap(36)

	screendivisor = screenHeight / float32(globalGameState.currentmap.height)
	intscreendivisor = int(screenHeight) / globalGameState.currentmap.height

}

// Pos variables for cursor
var (
	curspos pos
)

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
	mbrown           = color.RGBA{225, 216, 161, 255}
	mdarkgray        = color.RGBA{169, 169, 169, 255}
	mdarkgreen       = color.RGBA{75, 156, 0, 255}
	currenttilecolor = color.RGBA{0, 0, 0, 0}
)

// Screen sizes
var (
	width, height = ebiten.Monitor().Size()
	screenWidth   = float32(width)
	screenHeight  = float32(height)

	screendivisor    float32
	intscreendivisor int
)

// Standard positioning used for everything
type pos struct {
	int_x, int_y     int
	float_x, float_y float32
}

// Contains all information about the character
type character struct {
	title   string
	pos     pos
	picture *ebiten.Image
}

// Button struct definition
type button struct {
	title         string
	pos           pos
	width, height float32
	pressed       bool
	hovered       bool
	pressedColor  color.RGBA
	hoveredColor  color.RGBA
	inactiveColor color.RGBA
}

type slider struct {
	title          string
	pos            pos
	width, height  float32
	pressed        bool
	hovered        bool
	maxval, minval int
}

// var testmap = gamemap{
// 	data: [9][16]int{
// 		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
// 		{2, 3, 0, 0, 3, 3, 3, 3, 1, 1, 0, 2, 0, 1, 0, 2},
// 		{2, 3, 2, 3, 1, 3, 1, 2, 0, 0, 3, 3, 1, 2, 1, 2},
// 		{2, 1, 3, 3, 2, 0, 3, 0, 3, 1, 3, 3, 2, 3, 2, 2},
// 		{2, 1, 3, 3, 0, 2, 1, 1, 1, 2, 0, 3, 0, 0, 0, 2},
// 		{2, 3, 1, 0, 1, 1, 0, 1, 3, 3, 2, 1, 1, 3, 3, 2},
// 		{2, 2, 2, 3, 1, 0, 2, 1, 3, 3, 2, 1, 2, 3, 2, 2},
// 		{2, 1, 1, 3, 3, 1, 3, 2, 0, 0, 2, 2, 0, 0, 1, 2},
// 		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
// 	},
// }

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

	// contains where the character is on the array
	//
	// this translates to gamemap[y][x]
	//currentcharpos pos
}

// Create a new button
func createButton(title string, width, height float32, pressedColor, hoveredColor, inactiveColor color.RGBA, pos pos) button {
	return button{
		title:         title,
		pos:           pos,
		width:         width,
		height:        height,
		pressedColor:  pressedColor,
		hoveredColor:  hoveredColor,
		inactiveColor: inactiveColor,
	}
}

// Create a new slider
func createSlider(title string, width, height float32, minval, maxval int, pos pos) slider {
	return slider{
		title:  title,
		pos:    pos,
		width:  width,
		height: height,
		minval: minval,
		maxval: maxval,
	}
}

func onearg_createPos(u float32) pos {
	return pos{
		int_x:   int(u),
		int_y:   int(u),
		float_x: u,
		float_y: u,
	}
}

func createPos(x, y float32) pos {
	return pos{
		int_x:   int(x),
		int_y:   int(y),
		float_x: x,
		float_y: y,
	}
}

// DrawSlider draws a slider and check for interaction
func (s *slider) DrawSlider(screen *ebiten.Image) {

}

func createCharacter(path, title string) character {
	var c character
	// Open the image file
	file, err := os.Open(title)
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

	op.GeoM.Scale(100, 100)

	c.picture.DrawImage(screen, op)
}

// DrawButton draws the button and checks for interaction
func (b *button) DrawButton(screen *ebiten.Image) {
	// Check if the mouse is inside the button's area
	if curspos.float_x >= b.pos.float_x && curspos.float_x <= b.pos.float_x+b.width && curspos.float_y >= b.pos.float_y && curspos.float_y <= b.pos.float_y+b.height {
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
	vector.DrawFilledRect(screen, b.pos.float_x, b.pos.float_y, float32(b.width), float32(b.height), drawColor, false)

	// Draw the button title as text
	ebitenutil.DebugPrintAt(screen, b.title, b.pos.int_x+10, b.pos.int_y+10)
}

type Game struct{}

// Update method of the Game
func (g *Game) Update() error {
	return nil
}

// Draw method of the Game
func (g *Game) Draw(screen *ebiten.Image) {
	// Get the mouse position
	intmx, intmy := ebiten.CursorPosition()
	curspos.float_x, curspos.float_y = float32(intmx), float32(intmy)
	curspos.int_x, curspos.int_y = intmx, intmy

	//char := createCharacter("C:/Users/bence/Desktop/vsc/golang/rpg/character.png", "character")

	//char.DrawCharacter(screen)

	switch globalGameState.stateid {
	case 0:

		playbtn := createButton("Play", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
		playbtn.DrawButton(screen)
		if playbtn.pressed {
			globalGameState.stateid = 3
		}

		optionsbtn := createButton("Options", 150, 50, uitransparent, uilightgray, uigray, createPos(25, 85))
		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			globalGameState.stateid = 1
		}
	case 1:
		options_exitbtn := createButton("Back to menu", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			globalGameState.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
	case 3:
		ebitenutil.DebugPrint(screen, "Test map")

		// this is 32 because we are rendering that part of the map
		// one tile is 240x240 px so this would fill the whole screen
		// 0 = not decided, 1 = mountains, 2 = plains, 3 = hills, 4 = forests
		for i := 0; i < globalGameState.currentmap.height; i++ {
			for j := 0; j < globalGameState.currentmap.width; j++ {
				switch globalGameState.currentmap.data[i][j] {
				case 2:
					currenttilecolor = mlightgreen
				case 3:
					currenttilecolor = mbrown
				case 1:
					currenttilecolor = mdarkgray
				case 4:
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
