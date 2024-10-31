package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

// Standard positioning used for everything
type pos struct {
	int_x, int_y     int
	float_x, float_y float32
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

func debug() {
}

// ptid calculates and returns the tile coordinates based on the given position.
func ptid(pos pos) (int, int) {
	// Calculate the tile coordinates based on the character's position
	x := int(pos.float_x / screendivisor)
	y := int(pos.float_y / screendivisor)
	return x, y
}

// Pos variables for cursor
var (
	curspos pos
)

func (cursor *pos) updatemouse() {
	// Get the mouse position
	intmx, intmy := ebiten.CursorPosition()
	curspos.float_x, curspos.float_y = float32(intmx), float32(intmy)
	curspos.int_x, curspos.int_y = intmx, intmy
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

// DrawButton draws the button and checks for interaction
func (b *button) DrawButton(screen *ebiten.Image) {
	// Check if the mouse is inside the button's area
	if curspos.float_x >= b.pos.float_x && curspos.float_x <= b.pos.float_x+b.width && curspos.float_y >= b.pos.float_y && curspos.float_y <= b.pos.float_y+b.height {
		b.hovered = true
		// Check if the left mouse button is pressed
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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
