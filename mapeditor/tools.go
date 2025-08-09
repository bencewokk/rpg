package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ToolSystem struct {
	painting               bool
	lastPaintX, lastPaintY int
}

func NewToolSystem() ToolSystem {
	return ToolSystem{
		painting: false,
	}
}

func (t *ToolSystem) Update(mapData *MapData, camera *Camera) {
	mouseX, mouseY := ebiten.CursorPosition()

	// Don't paint if mouse is over UI area
	if mouseX < 120 {
		return
	}

	// Convert screen coordinates to tile coordinates
	tileX, tileY := camera.GetTileAtScreenPos(mouseX, mouseY)

	// Handle painting with left mouse button
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Only paint if we're not already painting this tile (to avoid rapid repainting)
		if !t.painting || t.lastPaintX != tileX || t.lastPaintY != tileY {
			// Get the selected tile type from UI (we'll need to pass this in)
			// For now, we'll get it from the global UI reference
			selectedType := 2 // Default to plains, this will be fixed
			mapData.SetTile(tileX, tileY, selectedType)

			t.painting = true
			t.lastPaintX = tileX
			t.lastPaintY = tileY
		}
	} else {
		t.painting = false
	}
}

// PaintTile paints a tile at the given coordinates with the specified type
func (t *ToolSystem) PaintTile(mapData *MapData, tileX, tileY, tileType int) {
	mapData.SetTile(tileX, tileY, tileType)
}
