package main

import (
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Returns either true or false
func calcChance(chance float64) bool {
	flip := float64(rand.Intn(100))
	return flip < chance
}

func loadPNG(path string) *ebiten.Image {
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
	return ebiten.NewImageFromImage(imgData)
}

// Color variable for ui
var (
	uigray        = color.RGBA{128, 128, 128, 255}
	uidarkgray    = color.RGBA{100, 100, 100, 255}
	uilightgray   = color.RGBA{158, 158, 158, 255}
	uilightgray2  = color.RGBA{190, 190, 190, 255}
	uitransparent = color.RGBA{0, 0, 0, 0}
	uilightred    = color.RGBA{190, 75, 71, 255}
	uidarkred     = color.RGBA{129, 0, 0, 255}
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
	float_x, float_y float32
}

func Distance(a, b pos) float32 {
	dx := float64(b.float_x - a.float_x)
	dy := float64(b.float_y - a.float_y)
	return float32(math.Sqrt(dx*dx + dy*dy))
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

	pressedColor  color.RGBA
	hoveredColor  color.RGBA
	inactiveColor color.RGBA

	knobpos  pos
	dragging bool
}

// ptid calculates and returns the tile coordinates based on the given position.
func ptid(pos pos) (int, int) {
	x := int(pos.float_x / screendivisor)
	y := int(pos.float_y / screendivisor)
	return x, y
}

// pos variables for cursor
var (
	curspos pos
)

// create a new button
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

// create a new slider
func createSlider(title string, width, height float32, minval, maxval int, pressedColor, hoveredColor, inactiveColor color.RGBA, pos pos) slider {
	kb := createPos(pos.float_x+8, pos.float_y+4)
	return slider{
		title:         title,
		pos:           pos,
		width:         width,
		height:        height,
		minval:        minval,
		maxval:        maxval,
		pressedColor:  pressedColor,
		hoveredColor:  hoveredColor,
		inactiveColor: inactiveColor,
		knobpos:       kb,
	}
}

func onearg_createPos(u float32) pos {
	return pos{
		float_x: u,
		float_y: u,
	}
}

func createPos(x, y float32) pos {
	return pos{
		float_x: x,
		float_y: y,
	}
}

func inSlide(s *slider) bool {

	if curspos.float_x >= s.pos.float_x+10 &&
		curspos.float_x <= s.pos.float_x+s.width-20 {
		return true
	}

	return false
}

// drawSlider draws a slider and checks for interaction
func (s *slider) DrawSlider(screen *ebiten.Image) {
	// Check if the cursor is hovering over the knob
	if curspos.float_x >= s.knobpos.float_x &&
		curspos.float_x <= s.knobpos.float_x+s.width/50 &&
		curspos.float_y >= s.knobpos.float_y &&
		curspos.float_y <= s.knobpos.float_y+(s.height-7) {

		s.hovered = true

		// Start dragging if the left mouse button is pressed and cursor is over the knob
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.pressed = true
			s.dragging = true
		} else {
			s.pressed = false
			// Stop dragging when the mouse button is released
			s.dragging = false
		}
	} else {
		s.hovered = false
		// If the cursor is no longer on the knob, release the drag state
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.pressed = false
			s.dragging = false
		}
	}

	// Choose color based on button state
	var drawColor color.Color
	if s.pressed {
		drawColor = s.pressedColor
	} else if s.hovered {
		drawColor = s.hoveredColor
	} else {
		drawColor = s.inactiveColor
	}

	//TODO check for outside

	// Move the knob if it is being dragged
	if s.dragging {
		s.knobpos.float_x = curspos.float_x - s.width/100
	}

	// Draw the slider track and knob
	vector.DrawFilledRect(screen, s.pos.float_x, s.pos.float_y, s.width, s.height, uidarkgray, false)
	vector.DrawFilledRect(screen, s.pos.float_x+5, s.pos.float_y+5, s.width-10, s.height-10, uilightgray2, false)
	vector.DrawFilledRect(screen, s.knobpos.float_x, s.knobpos.float_y, s.width/50, s.height-7, drawColor, false)

}

// drawButton draws the button and checks for interaction
func (b *button) DrawButton(screen *ebiten.Image) {
	// Check if the mouse is inside the button's area
	if curspos.float_x >= b.pos.float_x &&
		curspos.float_x <= b.pos.float_x+b.width &&
		curspos.float_y >= b.pos.float_y &&
		curspos.float_y <= b.pos.float_y+b.height {
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
	ebitenutil.DebugPrintAt(screen, b.title, int(b.pos.float_x+10), int(b.pos.float_y+10))
}
