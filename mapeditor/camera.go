package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Camera struct {
	X, Y  float64
	Zoom  float64
	dragging bool
	lastMouseX, lastMouseY int
}

func NewCamera() Camera {
	return Camera{
		X:    0,
		Y:    0,
		Zoom: 1.0,
	}
}

func (c *Camera) Update() {
	// Handle mouse panning
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Start dragging with middle mouse button
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		c.dragging = true
		c.lastMouseX, c.lastMouseY = mouseX, mouseY
	}
	
	// Stop dragging when middle mouse button is released
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		c.dragging = false
	}
	
	// Pan camera while dragging
	if c.dragging {
		deltaX := float64(mouseX - c.lastMouseX)
		deltaY := float64(mouseY - c.lastMouseY)
		
		c.X -= deltaX / c.Zoom
		c.Y -= deltaY / c.Zoom
		
		c.lastMouseX, c.lastMouseY = mouseX, mouseY
	}
	
	// Handle keyboard panning
	panSpeed := 300.0 / c.Zoom
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		c.X -= panSpeed / 60.0 // Assuming 60 FPS
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		c.X += panSpeed / 60.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Y -= panSpeed / 60.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Y += panSpeed / 60.0
	}
	
	// Handle zooming with mouse wheel (simple approach like main game)
	_, yScroll := ebiten.Wheel()
	if yScroll != 0 {
		// Simple zoom adjustment like in the main game
		if yScroll > 0 && c.Zoom < 4.0 {
			c.Zoom += 0.1
		} else if yScroll < 0 && c.Zoom > 0.25 {
			c.Zoom -= 0.1
		}
	}
}

// ScreenToWorld converts screen coordinates to world coordinates (matching main game approach)
func (c *Camera) ScreenToWorld(screenX, screenY int) (float64, float64) {
	// Convert screen coordinates back to world coordinates
	// Based on: screen = (world - camera) * zoom + screenCenter
	// Solve for world: world = (screen - screenCenter) / zoom + camera
	worldX := (float64(screenX)-float64(windowWidth)/2)/c.Zoom + c.X
	worldY := (float64(screenY)-float64(windowHeight)/2)/c.Zoom + c.Y
	return worldX, worldY
}

// WorldToScreen converts world coordinates to screen coordinates (matching main game approach)
func (c *Camera) WorldToScreen(worldX, worldY float64) (int, int) {
	// Based on main game's offsetsx/offsetsy functions:
	// screen = (world - camera) * zoom + screenCenter
	screenX := int((worldX-c.X)*c.Zoom + float64(windowWidth)/2)
	screenY := int((worldY-c.Y)*c.Zoom + float64(windowHeight)/2)
	return screenX, screenY
}

// GetTileAtScreenPos returns the tile coordinates at the given screen position
func (c *Camera) GetTileAtScreenPos(screenX, screenY int) (int, int) {
	worldX, worldY := c.ScreenToWorld(screenX, screenY)
	tileX := int(worldX / tileSize)
	tileY := int(worldY / tileSize)
	return tileX, tileY
}
