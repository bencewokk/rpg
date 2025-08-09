package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type UI struct {
	selectedTileType int
	showGrid         bool
	tileButtons      [4]Button
}

type Button struct {
	X, Y, W, H int
	Text       string
	Color      color.RGBA
	Pressed    bool
	Hovered    bool
}

func NewUI() UI {
	ui := UI{
		selectedTileType: 2, // Start with grass/plains
		showGrid:         true,
	}
	
	// Create tile type buttons
	buttonWidth := 80
	buttonHeight := 60
	startX := 20
	startY := 20
	
	tileNames := []string{"Void", "Mountain", "Plains", "Dry"}
	tileColors := []color.RGBA{
		voidColor,
		darkGray,
		lightGreen,
		lightBrown,
	}
	
	for i := 0; i < 4; i++ {
		ui.tileButtons[i] = Button{
			X:     startX,
			Y:     startY + i*(buttonHeight+10),
			W:     buttonWidth,
			H:     buttonHeight,
			Text:  tileNames[i],
			Color: tileColors[i],
		}
	}
	
	return ui
}

func (ui *UI) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Update buttons
	for i := 0; i < 4; i++ {
		btn := &ui.tileButtons[i]
		
		// Check if mouse is over button
		btn.Hovered = mouseX >= btn.X && mouseX < btn.X+btn.W &&
			mouseY >= btn.Y && mouseY < btn.Y+btn.H
		
		// Check for button click
		if btn.Hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			ui.selectedTileType = i
			btn.Pressed = true
		} else {
			btn.Pressed = false
		}
	}
	
	// Handle keyboard shortcuts
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		ui.selectedTileType = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		ui.selectedTileType = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		ui.selectedTileType = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		ui.selectedTileType = 3
	}
	
	// Toggle grid
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		ui.showGrid = !ui.showGrid
	}
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Draw tool panel background
	vector.DrawFilledRect(screen, 10, 10, 100, 300, mediumGray, false)
	vector.StrokeRect(screen, 10, 10, 100, 300, 1, darkGray, false)
	
	// Draw tile type buttons
	for i, btn := range ui.tileButtons {
		buttonColor := btn.Color
		if i == ui.selectedTileType {
			// Highlight selected button with a lighter version
			buttonColor = color.RGBA{
				uint8(min(int(btn.Color.R)+40, 255)),
				uint8(min(int(btn.Color.G)+40, 255)),
				uint8(min(int(btn.Color.B)+40, 255)),
				255,
			}
		}
		if btn.Hovered {
			// Make hovered buttons slightly brighter
			buttonColor = color.RGBA{
				uint8(min(int(buttonColor.R)+20, 255)),
				uint8(min(int(buttonColor.G)+20, 255)),
				uint8(min(int(buttonColor.B)+20, 255)),
				255,
			}
		}
		
		// Draw button background
		vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), buttonColor, false)
		
		// Draw button border
		borderColor := darkGray
		if i == ui.selectedTileType {
			borderColor = color.RGBA{0, 0, 0, 255} // Black border for selected
		}
		vector.StrokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 2, borderColor, false)
		
		// Draw button text
		textX := btn.X + 5
		textY := btn.Y + 15
		ebitenutil.DebugPrintAt(screen, btn.Text, textX, textY)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%d)", i), textX, textY+15)
	}
	
	// Draw instructions
	instructionsY := 320
	ebitenutil.DebugPrintAt(screen, "Controls:", 20, instructionsY)
	ebitenutil.DebugPrintAt(screen, "0-3: Select tile", 20, instructionsY+15)
	ebitenutil.DebugPrintAt(screen, "LClick: Paint tile", 20, instructionsY+30)
	ebitenutil.DebugPrintAt(screen, "MClick: Pan camera", 20, instructionsY+45)
	ebitenutil.DebugPrintAt(screen, "Wheel: Zoom", 20, instructionsY+60)
	ebitenutil.DebugPrintAt(screen, "G: Toggle grid", 20, instructionsY+75)
	ebitenutil.DebugPrintAt(screen, "Ctrl+S: Save", 20, instructionsY+90)
	ebitenutil.DebugPrintAt(screen, "Ctrl+O: Load", 20, instructionsY+105)
}

func (ui *UI) GetSelectedTileType() int {
	return ui.selectedTileType
}

func (ui *UI) ShouldShowGrid() bool {
	return ui.showGrid
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
