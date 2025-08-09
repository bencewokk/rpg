package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	windowWidth  = 1200
	windowHeight = 800
	tileSize     = 32
)

type MapEditor struct {
	camera  Camera
	mapData MapData
	ui      UI
	tools   ToolSystem
	assets  AssetManager
}

func (e *MapEditor) Update() error {
	e.camera.Update()
	e.ui.Update()
	
	// Handle file operations
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			e.mapData.SaveToFile("../map.txt")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyO) {
			e.mapData.LoadFromFile("../map.txt")
		}
	}
	
	// Update tools with current UI selection
	e.updateTools()
	
	return nil
}

func (e *MapEditor) updateTools() {
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Don't paint if mouse is over UI area
	if mouseX < 120 {
		return
	}
	
	// Convert screen coordinates to tile coordinates
	tileX, tileY := e.camera.GetTileAtScreenPos(mouseX, mouseY)
	
	// Handle painting with left mouse button
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		selectedType := e.ui.GetSelectedTileType()
		e.tools.PaintTile(&e.mapData, tileX, tileY, selectedType)
	}
}

func (e *MapEditor) Draw(screen *ebiten.Image) {
	// Clear screen with light background
	screen.Fill(lightGray)

	// Draw the map
	e.drawMap(screen)

	// Draw UI elements
	e.ui.Draw(screen)

	// Draw debug info
	fps := ebiten.CurrentFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", fps))
}

func (e *MapEditor) drawMap(screen *ebiten.Image) {
	// Calculate visible tiles based on camera position and zoom
	// Use same approach as main game
	startX := int((e.camera.X - float64(windowWidth)/(2*e.camera.Zoom)) / tileSize)
	startY := int((e.camera.Y - float64(windowHeight)/(2*e.camera.Zoom)) / tileSize)
	endX := int((e.camera.X + float64(windowWidth)/(2*e.camera.Zoom)) / tileSize) + 2
	endY := int((e.camera.Y + float64(windowHeight)/(2*e.camera.Zoom)) / tileSize) + 2

	// Clamp to map bounds
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	if endX > e.mapData.Width {
		endX = e.mapData.Width
	}
	if endY > e.mapData.Height {
		endY = e.mapData.Height
	}

	// Draw tiles using main game's coordinate system
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			tileType := e.mapData.GetTile(x, y)
			texture := e.assets.GetTileTexture(tileType)

			if texture != nil {
				op := &ebiten.DrawImageOptions{}
				
				// Apply zoom scaling
				op.GeoM.Scale(e.camera.Zoom, e.camera.Zoom)
				
				// Calculate screen position using same approach as main game
				worldX := float64(x * tileSize)
				worldY := float64(y * tileSize)
				screenX := (worldX-e.camera.X)*e.camera.Zoom + float64(windowWidth)/2
				screenY := (worldY-e.camera.Y)*e.camera.Zoom + float64(windowHeight)/2
				
				op.GeoM.Translate(screenX, screenY)
				screen.DrawImage(texture, op)
			}
		}
	}
	
	// Draw grid if enabled
	if e.ui.ShouldShowGrid() {
		e.drawGrid(screen, startX, startY, endX, endY)
	}
}

func (e *MapEditor) drawGrid(screen *ebiten.Image, startX, startY, endX, endY int) {
	gridColor := color.RGBA{100, 100, 100, 128}
	
	// Draw vertical lines using main game coordinate system
	for x := startX; x <= endX; x++ {
		worldX := float64(x * tileSize)
		screenX := float32((worldX-e.camera.X)*e.camera.Zoom + float64(windowWidth)/2)
		
		if screenX >= 0 && screenX < windowWidth {
			worldStartY := float64(startY * tileSize)
			worldEndY := float64(endY * tileSize)
			startScreenY := float32((worldStartY-e.camera.Y)*e.camera.Zoom + float64(windowHeight)/2)
			endScreenY := float32((worldEndY-e.camera.Y)*e.camera.Zoom + float64(windowHeight)/2)
			vector.StrokeLine(screen, screenX, startScreenY, screenX, endScreenY, 1, gridColor, false)
		}
	}
	
	// Draw horizontal lines using main game coordinate system
	for y := startY; y <= endY; y++ {
		worldY := float64(y * tileSize)
		screenY := float32((worldY-e.camera.Y)*e.camera.Zoom + float64(windowHeight)/2)
		
		if screenY >= 0 && screenY < windowHeight {
			worldStartX := float64(startX * tileSize)
			worldEndX := float64(endX * tileSize)
			startScreenX := float32((worldStartX-e.camera.X)*e.camera.Zoom + float64(windowWidth)/2)
			endScreenX := float32((worldEndX-e.camera.X)*e.camera.Zoom + float64(windowWidth)/2)
			vector.StrokeLine(screen, startScreenX, screenY, endScreenX, screenY, 1, gridColor, false)
		}
	}
}

func (e *MapEditor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	editor := &MapEditor{}

	// Initialize components
	editor.camera = NewCamera()
	editor.mapData = NewMapData(150, 100) // Width x Height from your RPG
	editor.ui = NewUI()
	editor.tools = NewToolSystem()
	editor.assets = NewAssetManager()

	// Try to load existing map from parent directory
	if err := editor.mapData.LoadFromFile("../map.txt"); err != nil {
		fmt.Printf("Could not load existing map: %v\n", err)
		fmt.Println("Starting with empty map...")
	}

	// Try to load assets from parent directory
	if err := editor.assets.LoadAssets("../import"); err != nil {
		fmt.Printf("Could not load assets: %v\n", err)
		fmt.Println("Using fallback colors...")
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("RPG Map Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
