package main

import (
	"math"
	"rpg/mapio"

	"github.com/hajimehoshi/ebiten/v2"
)

type ToolType int

const (
	ToolPaint ToolType = iota
	ToolBucket
	ToolNode
	ToolPath
)

// Action represents a single undoable action
type Action struct {
	ActionType string       // "paint" or "bucket"
	Changes    []TileChange // List of tile changes
}

// TileChange represents a single tile modification
type TileChange struct {
	X, Y     int
	OldValue int
	NewValue int
}

type ToolSystem struct {
	currentTool            ToolType
	painting               bool
	lastPaintX, lastPaintY int

	// Node editing
	selectedNodeID  int
	creatingPath    bool
	pathStartNodeID int

	// Undo/Redo system
	history      []Action
	historyIndex int
	maxHistory   int
}

func NewToolSystem() ToolSystem {
	return ToolSystem{
		currentTool:     ToolPaint,
		painting:        false,
		selectedNodeID:  -1,
		creatingPath:    false,
		pathStartNodeID: -1,
		history:         make([]Action, 0),
		historyIndex:    -1,
		maxHistory:      100, // Keep last 100 actions
	}
}

func (t *ToolSystem) Update(mapData *mapio.MapData, camera *Camera) {
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
func (t *ToolSystem) PaintTile(mapData *mapio.MapData, tileX, tileY, tileType int) {
	if t.currentTool == ToolPaint {
		// Record single tile change for undo
		oldValue := mapData.GetTile(tileX, tileY)
		if oldValue != tileType {
			changes := []TileChange{{X: tileX, Y: tileY, OldValue: oldValue, NewValue: tileType}}
			t.addToHistory(Action{ActionType: "paint", Changes: changes})
			mapData.SetTile(tileX, tileY, tileType)
		}
	} else if t.currentTool == ToolBucket {
		// Bucket fill - we'll collect all changes first
		changes := t.getBucketFillChanges(mapData, tileX, tileY, tileType)
		if len(changes) > 0 {
			t.addToHistory(Action{ActionType: "bucket", Changes: changes})
			t.applyChanges(mapData, changes)
		}
	}
}

// SetTool changes the current tool
func (t *ToolSystem) SetTool(tool ToolType) {
	t.currentTool = tool
}

// GetCurrentTool returns the current tool
func (t *ToolSystem) GetCurrentTool() ToolType {
	return t.currentTool
}

// GetToolName returns the name of the current tool
func (t *ToolSystem) GetToolName() string {
	switch t.currentTool {
	case ToolPaint:
		return "Paint"
	case ToolBucket:
		return "Bucket"
	case ToolNode:
		return "Node"
	case ToolPath:
		return "Path"
	default:
		return "Unknown"
	}
}

// bucketFill implements flood fill algorithm - now just for calculating changes
func (t *ToolSystem) getBucketFillChanges(mapData *mapio.MapData, startX, startY, newTileType int) []TileChange {
	// Get the original tile type at the starting position
	originalTileType := mapData.GetTile(startX, startY)

	// If the new tile type is the same as the original, do nothing
	if originalTileType == newTileType {
		return []TileChange{}
	}

	var changes []TileChange
	visited := make(map[int]map[int]bool)

	// Initialize visited map
	for y := 0; y < mapData.Height; y++ {
		visited[y] = make(map[int]bool)
	}

	// Use a stack-based flood fill to avoid recursion stack overflow
	type point struct {
		x, y int
	}

	stack := []point{{startX, startY}}

	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		x, y := current.x, current.y

		// Check bounds
		if x < 0 || x >= mapData.Width || y < 0 || y >= mapData.Height {
			continue
		}

		// Check if already visited
		if visited[y][x] {
			continue
		}

		// Check if this tile matches the original type
		if mapData.GetTile(x, y) != originalTileType {
			continue
		}

		// Mark as visited
		visited[y][x] = true

		// Record this change
		changes = append(changes, TileChange{
			X: x, Y: y,
			OldValue: originalTileType,
			NewValue: newTileType,
		})

		// Add adjacent tiles to stack (4-directional)
		stack = append(stack,
			point{x + 1, y}, // Right
			point{x - 1, y}, // Left
			point{x, y + 1}, // Down
			point{x, y - 1}, // Up
		)
	}

	return changes
}

// addToHistory adds an action to the undo history
func (t *ToolSystem) addToHistory(action Action) {
	// Remove any actions after current position (for redo invalidation)
	if t.historyIndex < len(t.history)-1 {
		t.history = t.history[:t.historyIndex+1]
	}

	// Add new action
	t.history = append(t.history, action)
	t.historyIndex++

	// Limit history size
	if len(t.history) > t.maxHistory {
		t.history = t.history[1:]
		t.historyIndex--
	}
}

// applyChanges applies a list of tile changes to the map
func (t *ToolSystem) applyChanges(mapData *mapio.MapData, changes []TileChange) {
	for _, change := range changes {
		mapData.SetTile(change.X, change.Y, change.NewValue)
	}
}

// Undo reverts the last action
func (t *ToolSystem) Undo(mapData *mapio.MapData) bool {
	if t.historyIndex < 0 {
		return false // Nothing to undo
	}

	action := t.history[t.historyIndex]

	// Revert all changes in reverse order
	for _, change := range action.Changes {
		mapData.SetTile(change.X, change.Y, change.OldValue)
	}

	t.historyIndex--
	return true
}

// Redo applies the next action in history
func (t *ToolSystem) Redo(mapData *mapio.MapData) bool {
	if t.historyIndex >= len(t.history)-1 {
		return false // Nothing to redo
	}

	t.historyIndex++
	action := t.history[t.historyIndex]

	// Apply all changes
	for _, change := range action.Changes {
		mapData.SetTile(change.X, change.Y, change.NewValue)
	}

	return true
}

// CanUndo returns true if there are actions to undo
func (t *ToolSystem) CanUndo() bool {
	return t.historyIndex >= 0
}

// HandleNodeTool handles node creation, selection, and deletion
func (t *ToolSystem) HandleNodeTool(mapData *mapio.MapData, worldX, worldY float64, leftClick, rightClick bool) {
	if leftClick {
		// Check if clicking on existing node
		nodeID := t.findNodeAtPosition(mapData, worldX, worldY, 16.0) // 16 pixel tolerance

		if nodeID >= 0 {
			// Select existing node
			t.selectedNodeID = nodeID
		} else {
			// Create new node
			newID := mapData.GetNextNodeID()
			mapData.AddNode(newID, float32(worldX), float32(worldY))
			t.selectedNodeID = newID
		}
	}

	if rightClick && t.selectedNodeID >= 0 {
		// Delete selected node
		mapData.RemoveNode(t.selectedNodeID)
		t.selectedNodeID = -1
	}
}

// HandlePathTool handles path creation between nodes
func (t *ToolSystem) HandlePathTool(mapData *mapio.MapData, worldX, worldY float64, leftClick, rightClick bool) {
	if leftClick {
		nodeID := t.findNodeAtPosition(mapData, worldX, worldY, 16.0)

		if nodeID >= 0 {
			if !t.creatingPath {
				// Start creating path
				t.pathStartNodeID = nodeID
				t.creatingPath = true
				t.selectedNodeID = nodeID
			} else {
				// Complete path creation
				if nodeID != t.pathStartNodeID {
					// Calculate distance as cost
					nodeA := mapData.FindNodeByID(t.pathStartNodeID)
					nodeB := mapData.FindNodeByID(nodeID)
					if nodeA != nil && nodeB != nil {
						dx := nodeB.Pos.X - nodeA.Pos.X
						dy := nodeB.Pos.Y - nodeA.Pos.Y
						cost := float32(dx*dx + dy*dy)           // Squared distance for cost
						cost = float32(math.Sqrt(float64(cost))) // Actual distance

						mapData.AddPath(t.pathStartNodeID, nodeID, cost)
					}
				}
				t.creatingPath = false
				t.pathStartNodeID = -1
			}
		}
	}

	if rightClick {
		// Cancel path creation
		t.creatingPath = false
		t.pathStartNodeID = -1
	}
}

// findNodeAtPosition finds a node within tolerance distance of the given position
func (t *ToolSystem) findNodeAtPosition(mapData *mapio.MapData, worldX, worldY, tolerance float64) int {
	for _, node := range mapData.Nodes {
		dx := float64(node.Pos.X) - worldX
		dy := float64(node.Pos.Y) - worldY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= tolerance {
			return node.ID
		}
	}
	return -1
}

// GetSelectedNodeID returns the currently selected node ID
func (t *ToolSystem) GetSelectedNodeID() int {
	return t.selectedNodeID
}

// IsCreatingPath returns true if currently in path creation mode
func (t *ToolSystem) IsCreatingPath() bool {
	return t.creatingPath
}

// GetPathStartNodeID returns the node ID where path creation started
func (t *ToolSystem) GetPathStartNodeID() int {
	return t.pathStartNodeID
}
