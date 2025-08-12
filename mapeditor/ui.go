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
	toolButtons      [6]Button // added Spawner
	selectedTool     ToolType
	statusMessage    string
	statusTimer      int
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
		selectedTool:     ToolPaint, // Start with paint tool
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

	// Create tool buttons (add NPC)
	toolNames := []string{"Paint", "Bucket", "Node", "Path", "NPC", "Spawner"}
	toolY := startY + 4*(buttonHeight+10) + 20 // Below tile buttons

	for i := 0; i < len(toolNames); i++ {
		ui.toolButtons[i] = Button{
			X:     startX,
			Y:     toolY + i*(buttonHeight/2+5),
			W:     buttonWidth,
			H:     buttonHeight / 2,
			Text:  toolNames[i],
			Color: lightGray,
		}
	}

	return ui
}

func (ui *UI) Update() {
	mouseX, mouseY := ebiten.CursorPosition()

	// Update status timer
	if ui.statusTimer > 0 {
		ui.statusTimer--
		if ui.statusTimer <= 0 {
			ui.statusMessage = ""
		}
	}

	// Update tile buttons
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

	// Update tool buttons
	for i := 0; i < len(ui.toolButtons); i++ { // iterate tools
		btn := &ui.toolButtons[i]

		// Check if mouse is over button
		btn.Hovered = mouseX >= btn.X && mouseX < btn.X+btn.W &&
			mouseY >= btn.Y && mouseY < btn.Y+btn.H

		// Check for button click
		if btn.Hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			ui.selectedTool = ToolType(i)
			btn.Pressed = true
		} else {
			btn.Pressed = false
		}
	}

	// Handle keyboard shortcuts only when Ctrl held so normal typing works in editors
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		// Tile selections
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
		// Tool selections
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			ui.selectedTool = ToolPaint
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			ui.selectedTool = ToolBucket
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			ui.selectedTool = ToolNode
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			ui.selectedTool = ToolPath
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			ui.selectedTool = ToolNPC
		}
		if !ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
			ui.selectedTool = ToolSpawner
		}
		// Toggle grid
		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			ui.showGrid = !ui.showGrid
		}
	}
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Draw tool panel background
	vector.DrawFilledRect(screen, 10, 10, 100, 500, mediumGray, false)
	vector.StrokeRect(screen, 10, 10, 100, 500, 1, darkGray, false)

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

	// Draw tool buttons
	for i, btn := range ui.toolButtons {
		buttonColor := btn.Color
		if ToolType(i) == ui.selectedTool {
			buttonColor = lightBlue // Highlight selected tool
		}
		if btn.Hovered {
			buttonColor = color.RGBA{
				uint8(min(int(buttonColor.R)+20, 255)),
				uint8(min(int(buttonColor.G)+20, 255)),
				uint8(min(int(buttonColor.B)+20, 255)),
				255,
			}
		}

		// Draw tool button background
		vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), buttonColor, false)

		// Draw tool button border
		borderColor := darkGray
		if ToolType(i) == ui.selectedTool {
			borderColor = color.RGBA{0, 0, 0, 255} // Black border for selected
		}
		vector.StrokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 2, borderColor, false)

		// Draw tool button text
		textX := btn.X + 5
		textY := btn.Y + 8
		ebitenutil.DebugPrintAt(screen, btn.Text, textX, textY)
	}

	// Draw instructions
	instructionsY := 420
	ebitenutil.DebugPrintAt(screen, "Controls:", 20, instructionsY)
	ebitenutil.DebugPrintAt(screen, "0-3: Select tile", 20, instructionsY+15)
	ebitenutil.DebugPrintAt(screen, "P: Paint tool", 20, instructionsY+30)
	ebitenutil.DebugPrintAt(screen, "B: Bucket tool", 20, instructionsY+45)
	ebitenutil.DebugPrintAt(screen, "N: Node tool", 20, instructionsY+60)
	ebitenutil.DebugPrintAt(screen, "M: Path tool", 20, instructionsY+75)
	ebitenutil.DebugPrintAt(screen, "C: NPC tool", 20, instructionsY+90)
	ebitenutil.DebugPrintAt(screen, "S: Spawner tool", 20, instructionsY+105)
	ebitenutil.DebugPrintAt(screen, "LClick: Use tool", 20, instructionsY+120)
	ebitenutil.DebugPrintAt(screen, "RClick: Delete/Cancel", 20, instructionsY+135)
	ebitenutil.DebugPrintAt(screen, "MClick: Pan camera", 20, instructionsY+150)
	ebitenutil.DebugPrintAt(screen, "Wheel: Zoom", 20, instructionsY+165)
	ebitenutil.DebugPrintAt(screen, "G: Toggle grid", 20, instructionsY+180)
	ebitenutil.DebugPrintAt(screen, "Ctrl+Shift+S: Save", 20, instructionsY+195)
	ebitenutil.DebugPrintAt(screen, "Ctrl+Z: Undo", 20, instructionsY+210)
	ebitenutil.DebugPrintAt(screen, "Ctrl+Y: Redo", 20, instructionsY+225)

	// Draw status message if active
	if ui.statusMessage != "" {
		ebitenutil.DebugPrintAt(screen, ui.statusMessage, 120, 20)
	}

	// NPC editing helper panel (draw to right of main panel) - read-only summary for now
	if ui.selectedTool == ToolNPC {
		vector.DrawFilledRect(screen, 10, 520, 100, 140, mediumGray, false)
		ebitenutil.DebugPrintAt(screen, "NPC Tool", 20, 525)
		if ui.selectedTool == ToolSpawner {
			vector.DrawFilledRect(screen, 10, 520, 100, 100, mediumGray, false)
			ebitenutil.DebugPrintAt(screen, "Spawner", 20, 525)
			ebitenutil.DebugPrintAt(screen, "L: place", 20, 540)
			ebitenutil.DebugPrintAt(screen, "R: delete", 20, 555)
		}
		ebitenutil.DebugPrintAt(screen, "L: place/select", 20, 540)
		ebitenutil.DebugPrintAt(screen, "R: del", 20, 555)
		ebitenutil.DebugPrintAt(screen, "+: add line", 20, 570)
		ebitenutil.DebugPrintAt(screen, "-: pop line", 20, 585)
		// We can't access mapData directly here; selection summary is drawn in main Draw if needed.
	}
	if ui.selectedTool == ToolSpawner {
		vector.DrawFilledRect(screen, 10, 520, 100, 120, mediumGray, false)
		ebitenutil.DebugPrintAt(screen, "Spawner", 20, 525)
		ebitenutil.DebugPrintAt(screen, "Q/E rad -/+", 12, 540)
		ebitenutil.DebugPrintAt(screen, "Z/X cnt -/+", 12, 555)
		ebitenutil.DebugPrintAt(screen, "F/R int -/+", 12, 570)
		ebitenutil.DebugPrintAt(screen, "Ctrl bigger", 12, 585)
		ebitenutil.DebugPrintAt(screen, "L place", 20, 600)
		ebitenutil.DebugPrintAt(screen, "R del", 20, 615)
	}
}

func (ui *UI) GetSelectedTileType() int {
	return ui.selectedTileType
}

func (ui *UI) GetSelectedTool() ToolType {
	return ui.selectedTool
}

func (ui *UI) ShouldShowGrid() bool {
	return ui.showGrid
}

func (ui *UI) ShowStatus(message string) {
	ui.statusMessage = message
	ui.statusTimer = 120 // Show for 2 seconds at 60fps
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
