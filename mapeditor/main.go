package main

import (
	"fmt"
	"image/color"
	"log"
	"rpg/mapio"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	windowWidth  = 1200
	windowHeight = 800
	// Match main game's screendivisor/intscreendivisor (30) so node coordinates align 1:1.
	// This removes the apparent "scaled & shifted" effect that came from using 32 here.
	tileSize = 30
)

type MapEditor struct {
	camera  Camera
	mapData *mapio.MapData
	ui      UI
	tools   ToolSystem
	assets  AssetManager
	// dialogue editing state
	editingDialogue bool
	editDialogueIdx int
	editBuffer      string
}

func (e *MapEditor) Update() error {
	e.camera.Update()
	e.ui.Update()

	// NPC editing & inline dialogue editing
	if e.ui.selectedTool == ToolNPC {
		idx := e.tools.GetSelectedNPC()
		if idx >= 0 && idx < len(e.mapData.NPCs) {
			n := &e.mapData.NPCs[idx]
			// mouse click to select line
			mouseX, mouseY := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				panelX, panelY := 120, 10
				lineStartY := panelY + 55
				lineHeight := 14
				maxLinesDisplay := 8
				if mouseX >= panelX && mouseX <= panelX+260 && mouseY >= lineStartY && mouseY <= lineStartY+maxLinesDisplay*lineHeight {
					clickedIdx := (mouseY - lineStartY) / lineHeight
					if clickedIdx >= 0 && clickedIdx < len(n.Dialogues) && clickedIdx < maxLinesDisplay {
						e.editingDialogue = true
						e.editDialogueIdx = clickedIdx
						e.editBuffer = n.Dialogues[clickedIdx]
					}
				}
			}
			// add new line (Ctrl + '+')
			if ebiten.IsKeyPressed(ebiten.KeyControl) && (inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd)) {
				n.Dialogues = append(n.Dialogues, "New line")
				if !e.editingDialogue {
					e.editDialogueIdx = len(n.Dialogues) - 1
					e.editBuffer = n.Dialogues[e.editDialogueIdx]
					e.editingDialogue = true
				}
			}
			// remove last line (Ctrl + '-')
			if ebiten.IsKeyPressed(ebiten.KeyControl) && (inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract)) {
				if len(n.Dialogues) > 0 {
					n.Dialogues = n.Dialogues[:len(n.Dialogues)-1]
					if e.editDialogueIdx >= len(n.Dialogues) {
						e.editDialogueIdx = len(n.Dialogues) - 1
						if e.editDialogueIdx < 0 {
							e.editingDialogue = false
						}
					}
				}
			}
			// cycle sprite path placeholders (Ctrl + '[' or ']')
			if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyLeftBracket) {
				n.SpritePath = "import/Characters/hamster.png"
			}
			if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyRightBracket) {
				n.SpritePath = "import/Characters/Orc.png"
			}

			// cycle through scanned sprite list with Ctrl+Up / Ctrl+Down
			if ebiten.IsKeyPressed(ebiten.KeyControl) && (inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown)) {
				paths := e.assets.GetNPCSpritePaths()
				if len(paths) > 0 {
					// find current index
					idxCur := -1
					for i, pth := range paths {
						if pth == n.SpritePath {
							idxCur = i
							break
						}
					}
					if idxCur == -1 {
						idxCur = 0
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
						idxCur = (idxCur + 1) % len(paths)
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
						idxCur = (idxCur - 1 + len(paths)) % len(paths)
					}
					n.SpritePath = paths[idxCur]
				}
			}

			// handle text editing
			if e.editingDialogue && e.editDialogueIdx >= 0 && e.editDialogueIdx < len(n.Dialogues) {
				// commit/cancel
				if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) {
					n.Dialogues[e.editDialogueIdx] = strings.TrimSpace(e.editBuffer)
					e.editingDialogue = false
				}
				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					e.editingDialogue = false
				}
				// backspace
				if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
					if len(e.editBuffer) > 0 {
						_, size := utf8.DecodeLastRuneInString(e.editBuffer)
						if size > 0 {
							e.editBuffer = e.editBuffer[:len(e.editBuffer)-size]
						}
					}
				}
				// append typed chars
				for _, r := range ebiten.InputChars() {
					if r >= 32 && r != 127 {
						e.editBuffer += string(r)
					}
				}
			}
		}
	}

	// Handle file operations
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			mapio.SaveMapToFile(e.mapData, "../map.txt")
			fmt.Println("Map saved!")
			e.ui.ShowStatus("Map saved!")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyO) {
			var err error
			e.mapData, err = mapio.LoadMapFromFile("../map.txt")
			if err != nil {
				fmt.Printf("Error loading map: %v\n", err)
				e.ui.ShowStatus("Error loading map!")
			} else {
				fmt.Println("Map loaded!")
				e.ui.ShowStatus("Map loaded!")
			}
		}
		// Undo with Ctrl+Z
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
			if e.tools.Undo(e.mapData) {
				fmt.Println("Undone!")
				e.ui.ShowStatus("Undone!")
			} else {
				e.ui.ShowStatus("Nothing to undo")
			}
		}
		// Redo with Ctrl+Y
		if inpututil.IsKeyJustPressed(ebiten.KeyY) {
			if e.tools.Redo(e.mapData) {
				fmt.Println("Redone!")
				e.ui.ShowStatus("Redone!")
			} else {
				e.ui.ShowStatus("Nothing to redo")
			}
		}
	}

	// Update tools with current UI selection
	e.updateTools()

	return nil
}

// Helper functions to replicate main game's offsetsx/offsetsy exactly
func (e *MapEditor) offsetsx(tobeoffset float32) float32 {
	// Use our window width but same coordinate transformation logic
	return ((tobeoffset-float32(e.camera.X))*float32(e.camera.Zoom) + float32(windowWidth)/2)
}

func (e *MapEditor) offsetsy(tobeoffset float32) float32 {
	// Use our window height but same coordinate transformation logic
	return ((tobeoffset-float32(e.camera.Y))*float32(e.camera.Zoom) + float32(windowHeight)/2)
}

func (e *MapEditor) updateTools() {
	// Sync tool system with UI selection
	e.tools.SetTool(e.ui.GetSelectedTool())

	mouseX, mouseY := ebiten.CursorPosition()

	// Don't interact if mouse is over UI area
	if mouseX < 120 {
		return
	}

	// Convert screen coordinates to world coordinates
	worldX, worldY := e.camera.ScreenToWorld(mouseX, mouseY)

	// Get input states
	leftClick := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	rightClick := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	leftHeld := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	currentTool := e.ui.GetSelectedTool()

	// If NPC tool and cursor is over the NPC overlay panel, suppress tool interactions
	if currentTool == ToolNPC {
		idx := e.tools.GetSelectedNPC()
		if idx >= 0 && idx < len(e.mapData.NPCs) {
			// Mirror overlay geometry from Draw()
			panelX, panelY := 120, 10
			panelW := 260
			// Compute dynamic panel height similar to Draw()
			lines := len(e.mapData.NPCs[idx].Dialogues)
			maxLinesDisplay := 8
			if lines > maxLinesDisplay {
				lines = maxLinesDisplay
			}
			panelH := 110 + lines*14
			if panelH < 140 {
				panelH = 140
			}
			if mouseX >= panelX && mouseX <= panelX+panelW && mouseY >= panelY && mouseY <= panelY+panelH {
				// Inside overlay: let Update() handle line editing clicks; skip placement/deletion
				return
			}
		}
	}

	switch currentTool {
	case ToolPaint, ToolBucket:
		// Handle painting/bucket with left mouse button held
		if leftHeld {
			// Convert to tile coordinates
			tileX, tileY := e.camera.GetTileAtScreenPos(mouseX, mouseY)
			selectedType := e.ui.GetSelectedTileType()
			e.tools.PaintTile(e.mapData, tileX, tileY, selectedType)
		}
	case ToolNode:
		// Handle node tool
		e.tools.HandleNodeTool(e.mapData, worldX, worldY, leftClick, rightClick)
	case ToolPath:
		// Handle path tool
		e.tools.HandlePathTool(e.mapData, worldX, worldY, leftClick, rightClick)
	case ToolNPC:
		// Handle NPC placement/removal
		e.tools.HandleNPCTool(e.mapData, worldX, worldY, leftClick, rightClick)
	}
}

func (e *MapEditor) Draw(screen *ebiten.Image) {
	// Clear screen with light background
	screen.Fill(lightGray)

	// Draw the map
	e.drawMap(screen)

	// Draw nodes and paths
	e.drawNodes(screen)

	// Draw UI elements
	e.ui.Draw(screen)

	// NPC overlay with inline editing
	if e.ui.selectedTool == ToolNPC {
		idx := e.tools.GetSelectedNPC()
		if idx >= 0 && idx < len(e.mapData.NPCs) {
			n := e.mapData.NPCs[idx]
			panelX, panelY := 120, 10
			maxLinesDisplay := 8
			linesShown := len(n.Dialogues)
			if linesShown > maxLinesDisplay {
				linesShown = maxLinesDisplay
			}
			panelH := float32(110 + linesShown*14)
			if panelH < 140 {
				panelH = 140
			}
			vector.DrawFilledRect(screen, float32(panelX), float32(panelY), 260, panelH, color.RGBA{50, 50, 60, 200}, false)
			vector.StrokeRect(screen, float32(panelX), float32(panelY), 260, panelH, 2, color.RGBA{0, 0, 0, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("NPC %d", idx), panelX+10, panelY+10)
			ebitenutil.DebugPrintAt(screen, "Name:"+n.Name, panelX+10, panelY+25)
			ebitenutil.DebugPrintAt(screen, "Sprite:"+n.SpritePath, panelX+10, panelY+40)
			lineBaseY := panelY + 55
			for i := 0; i < linesShown; i++ {
				text := n.Dialogues[i]
				if e.editingDialogue && e.editDialogueIdx == i {
					text = e.editBuffer + "|"
				}
				// highlight row
				if e.editDialogueIdx == i && e.editingDialogue {
					vector.DrawFilledRect(screen, float32(panelX+5), float32(lineBaseY+i*14-2), 250, 14, color.RGBA{90, 90, 110, 220}, false)
				} else if e.editDialogueIdx == i {
					vector.DrawFilledRect(screen, float32(panelX+5), float32(lineBaseY+i*14-2), 250, 14, color.RGBA{70, 70, 90, 180}, false)
				}
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", i+1, text), panelX+10, lineBaseY+i*14)
			}
			if len(n.Dialogues) > maxLinesDisplay {
				ebitenutil.DebugPrintAt(screen, "(scroll TBD)", panelX+10, lineBaseY+linesShown*14)
			}
			// instructions footer
			instrY := int(panelY) + int(panelH) - 44
			paths := e.assets.GetNPCSpritePaths()
			if len(paths) > 0 {
				// show index
				curIdx := -1
				for i, p := range paths {
					if p == n.SpritePath {
						curIdx = i
						break
					}
				}
				if curIdx >= 0 {
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Sprite %d/%d", curIdx+1, len(paths)), panelX+10, instrY)
				}
				instrY += 14
				ebitenutil.DebugPrintAt(screen, "Ctrl+Up/Down cycle sprites", panelX+10, instrY)
				instrY += 14
			}
			ebitenutil.DebugPrintAt(screen, "+ add  - remove", panelX+10, instrY)
			instrY += 14
			ebitenutil.DebugPrintAt(screen, "Click line then type, Enter=save Esc=cancel", panelX+10, instrY)
		}
	}

	// Draw debug info
	fps := ebiten.CurrentFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", fps))
}

func (e *MapEditor) drawMap(screen *ebiten.Image) {
	// Calculate visible tiles based on camera position and zoom
	// Use same approach as main game
	startX := int((e.camera.X - float64(windowWidth)/(2*e.camera.Zoom)) / tileSize)
	startY := int((e.camera.Y - float64(windowHeight)/(2*e.camera.Zoom)) / tileSize)
	endX := int((e.camera.X+float64(windowWidth)/(2*e.camera.Zoom))/tileSize) + 2
	endY := int((e.camera.Y+float64(windowHeight)/(2*e.camera.Zoom))/tileSize) + 2

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

				// Center tiles like in the main game (tile spans center +/- tileSize/2)
				worldX := float64(x*tileSize - tileSize/2)
				worldY := float64(y*tileSize - tileSize/2)
				screenX := (worldX-e.camera.X)*e.camera.Zoom + float64(windowWidth)/2
				screenY := (worldY-e.camera.Y)*e.camera.Zoom + float64(windowHeight)/2

				op.GeoM.Translate(screenX, screenY)
				screen.DrawImage(texture, op)
			}
		}
	}

	// Draw NPC markers (simple circles with initial letter)
	for _, n := range e.mapData.NPCs {
		sx := e.offsetsx(n.Pos.X)
		sy := e.offsetsy(n.Pos.Y)
		vector.DrawFilledCircle(screen, sx, sy, 6, color.RGBA{200, 180, 60, 255}, false)
		vector.StrokeCircle(screen, sx, sy, 6, 2, color.RGBA{0, 0, 0, 255}, false)
		label := "N"
		if len(n.Name) > 0 {
			label = string([]rune(n.Name)[0])
		}
		ebitenutil.DebugPrintAt(screen, label, int(sx)-3, int(sy)-20)
	}
}

func (e *MapEditor) drawNodes(screen *ebiten.Image) {
	// Draw paths first (so they appear behind nodes)
	for _, path := range e.mapData.Paths {
		nodeA := e.mapData.FindNodeByID(path.NodeAID)
		nodeB := e.mapData.FindNodeByID(path.NodeBID)

		if nodeA != nil && nodeB != nil {
			// Use exact same coordinate transformation as main game
			screenX1 := e.offsetsx(nodeA.Pos.X)
			screenY1 := e.offsetsy(nodeA.Pos.Y)
			screenX2 := e.offsetsx(nodeB.Pos.X)
			screenY2 := e.offsetsy(nodeB.Pos.Y)

			// Draw path line
			vector.StrokeLine(screen, screenX1, screenY1, screenX2, screenY2, 2, color.RGBA{100, 100, 255, 255}, false)
		}
	}

	// Draw path being created
	if e.tools.IsCreatingPath() {
		startNodeID := e.tools.GetPathStartNodeID()
		startNode := e.mapData.FindNodeByID(startNodeID)
		if startNode != nil {
			mouseX, mouseY := ebiten.CursorPosition()
			screenX1 := e.offsetsx(startNode.Pos.X)
			screenY1 := e.offsetsy(startNode.Pos.Y)

			// Draw preview line to mouse cursor
			vector.StrokeLine(screen, screenX1, screenY1, float32(mouseX), float32(mouseY), 2, color.RGBA{255, 255, 100, 128}, false)
		}
	}

	// Draw nodes
	for _, node := range e.mapData.Nodes {
		// Use exact same coordinate transformation as main game
		screenX := e.offsetsx(node.Pos.X)
		screenY := e.offsetsy(node.Pos.Y)

		// Choose node color based on selection
		nodeColor := color.RGBA{255, 100, 100, 255} // Default red
		if node.ID == e.tools.GetSelectedNodeID() {
			nodeColor = color.RGBA{100, 255, 100, 255} // Green if selected
		}
		if node.ID == e.tools.GetPathStartNodeID() {
			nodeColor = color.RGBA{255, 255, 100, 255} // Yellow if path start
		}

		// Draw node circle
		radius := float32(8.0 * e.camera.Zoom)
		if radius < 4 {
			radius = 4
		}
		vector.DrawFilledCircle(screen, screenX, screenY, radius, nodeColor, false)
		vector.StrokeCircle(screen, screenX, screenY, radius, 2, color.RGBA{0, 0, 0, 255}, false)

		// Draw node ID
		if e.camera.Zoom >= 0.5 { // Only show IDs when zoomed in enough
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", node.ID), int(screenX)-5, int(screenY)-20)
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
	editor.mapData = mapio.NewMapData(150, 100) // Width x Height from your RPG
	editor.ui = NewUI()
	editor.tools = NewToolSystem()
	editor.assets = NewAssetManager()

	// Try to load existing map from the main project
	if mapData, err := mapio.LoadMapFromFile("../map.txt"); err != nil {
		fmt.Printf("Could not load existing map: %v\n", err)
		fmt.Println("Starting with empty map...")
	} else {
		editor.mapData = mapData
		fmt.Printf("Loaded map: %dx%d with %d nodes, %d paths, %d sprites\n",
			mapData.Width, mapData.Height, len(mapData.Nodes), len(mapData.Paths), len(mapData.Sprites))
	}

	// Try to load assets from import directory
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
