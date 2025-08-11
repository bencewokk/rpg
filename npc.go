package main

import (
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// npc represents a simple non-hostile talkable entity.
type npc struct {
	pos        pos
	texture    *ebiten.Image
	id         int
	dialogue   []string
	line       int
	talking    bool
	talkRadius float32
	name       string
	// simple internal animation frames (fallback, separate from global AnimationManager)
	frames        []*ebiten.Image
	frameIndex    int
	frameElapsed  float64
	frameDuration float64
}

var (
	activeNPC      *npc
	npcSprite      *ebiten.Image // legacy default
	npcSpriteCache = map[string]*ebiten.Image{}
)

// createNPC adds an NPC to the current map with the provided dialogue lines.
func createNPC(p pos, lines []string) { createNPCWithSprite(p, lines, "import/Characters/hamster.png") }

// createNPCWithSprite allows specifying a sprite path; caches images; falls back on default if load fails.
func createNPCWithSprite(p pos, lines []string, spritePath string) {
	if spritePath == "" || spritePath == "-" {
		spritePath = "import/Characters/hamster.png"
	}
	img, ok := npcSpriteCache[spritePath]
	if !ok {
		defer func() { // recover from potential load issues silently
			if r := recover(); r != nil {
				img = npcSprite
			}
		}()
		img = loadPNG(spritePath)
		if img != nil {
			npcSpriteCache[spritePath] = img
		}
	}
	if img == nil {
		if npcSprite == nil {
			npcSprite = loadPNG("import/Characters/hamster.png")
		}
		img = npcSprite
	}
	sheetW, sheetH := img.Size()
	var frames []*ebiten.Image
	// Detect layout types:
	// 1. Horizontal strip (sheetW > sheetH)
	// 2. Vertical strip (sheetH > sheetW)
	// 3. Grid (sheetW == sheetH containing NxN frames)
	if sheetH > 0 && sheetW > sheetH { // horizontal strip
		frameSize := sheetH
		frameCount := sheetW / frameSize
		for i := 0; i < frameCount; i++ {
			r := image.Rect(i*frameSize, 0, (i+1)*frameSize, frameSize)
			if sub, ok := img.SubImage(r).(*ebiten.Image); ok {
				frames = append(frames, sub)
			}
		}
	} else if sheetW > 0 && sheetH > sheetW { // vertical strip
		frameSize := sheetW
		frameCount := sheetH / frameSize
		for i := 0; i < frameCount; i++ {
			r := image.Rect(0, i*frameSize, frameSize, (i+1)*frameSize)
			if sub, ok := img.SubImage(r).(*ebiten.Image); ok {
				frames = append(frames, sub)
			}
		}
	} else if sheetW == sheetH && sheetW > 0 { // potential grid (e.g., 2x2, 3x3, 4x4)
		// Find a reasonable grid factor (prefer 4, then 3, then 2, else no grid)
		candidates := []int{8, 6, 5, 4, 3, 2}
		var grid int
		for _, c := range candidates {
			if sheetW%c == 0 {
				grid = c
				break
			}
		}
		if grid > 1 { // slice grid
			cell := sheetW / grid
			// Avoid absurdly tiny cells (<8px) – treat as single image instead.
			if cell >= 8 {
				for y := 0; y < grid; y++ {
					for x := 0; x < grid; x++ {
						r := image.Rect(x*cell, y*cell, (x+1)*cell, (y+1)*cell)
						if sub, ok := img.SubImage(r).(*ebiten.Image); ok {
							frames = append(frames, sub)
						}
					}
				}
			}
		}
	}
	baseImg := img
	if len(frames) > 0 {
		baseImg = frames[0]
	} else {
		// Fallback: crop leftmost square so we don't render the entire sheet of frames.
		minDim := sheetW
		if sheetH < minDim {
			minDim = sheetH
		}
		if minDim > 0 {
			if sub, ok := img.SubImage(image.Rect(0, 0, minDim, minDim)).(*ebiten.Image); ok {
				baseImg = sub
			}
		}
	}
	n := &npc{pos: p, texture: baseImg, dialogue: lines, talkRadius: 110, name: "NPC", frames: frames, frameDuration: 0.15}
	game.currentmap.npcs = append(game.currentmap.npcs, n)
	drawables = append(drawables, n)
}

// draw implements drawable.
func (n *npc) draw(screen *ebiten.Image) {
	if n.texture == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	// Match earlier entity scale approach (tile sized ~18px base)
	scale := float64(screendivisor) / 18 * float64(game.camera.zoom)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(offsetsx(n.pos.float_x))-float64(screendivisor),
		float64(offsetsy(n.pos.float_y))-float64(screendivisor),
	)
	screen.DrawImage(n.texture, op)

	// Small indicator when in range and not already talking
	if activeNPC == nil {
		c := nearestCharacter(n.pos)
		if c != nil && Distance(c.pos, n.pos) < n.talkRadius {
			ebitenutil.DebugPrintAt(screen, "[E] Talk", int(offsetsx(n.pos.float_x))-20, int(offsetsy(n.pos.float_y))-40)
		}
	}
}

// advance animation based on dt
func (n *npc) update(dt float64) {
	if len(n.frames) == 0 {
		return
	}
	n.frameElapsed += dt
	if n.frameElapsed >= n.frameDuration {
		n.frameElapsed -= n.frameDuration
		n.frameIndex = (n.frameIndex + 1) % len(n.frames)
		n.texture = n.frames[n.frameIndex]
	}
}

// update all NPC animations
func updateNPCAnimations(dt float64) {
	for _, n := range game.currentmap.npcs {
		n.update(dt)
	}
}

// Y depth ordering (similar to enemy offset)
func (n *npc) Y() float32    { return n.pos.float_y }
func (n *npc) giveId(id int) { n.id = id }

// updateNPCInteractions handles starting and advancing conversations.
func updateNPCInteractions() {
	if len(game.currentmap.players) == 0 {
		return
	}
	player := game.currentmap.players[0]

	// If already talking
	if activeNPC != nil {
		// Advance dialogue
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			activeNPC.line++
			if activeNPC.line >= len(activeNPC.dialogue) {
				// End conversation
				activeNPC.talking = false
				activeNPC = nil
			}
		}
		// Cancel with Escape
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyE) {
			activeNPC.talking = false
			activeNPC = nil
		}
		return
	}

	// Not in conversation: check for nearby NPC and E press
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		var closest *npc
		var closestDist float32 = 1e9
		for _, n := range game.currentmap.npcs {
			d := Distance(player.pos, n.pos)
			if d < n.talkRadius && d < closestDist {
				closestDist = d
				closest = n
			}
		}
		if closest != nil {
			activeNPC = closest
			activeNPC.talking = true
			activeNPC.line = 0
		}
	}
}

// drawConversationUI renders the active conversation box.
func drawConversationUI(screen *ebiten.Image) {
	if activeNPC == nil {
		return
	}
	line := ""
	if activeNPC.line < len(activeNPC.dialogue) {
		line = activeNPC.dialogue[activeNPC.line]
	}
	// Box dimensions
	w := float32(screenWidth * 0.5)
	h := float32(90)
	x := (screenWidth - w) / 2
	y := screenHeight - h - 40
	// Background
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{0, 0, 0, 180}, false)
	// Border
	vector.DrawFilledRect(screen, x, y, w, 3, color.RGBA{255, 255, 255, 200}, false)
	vector.DrawFilledRect(screen, x, y+h-3, w, 3, color.RGBA{255, 255, 255, 200}, false)
	vector.DrawFilledRect(screen, x, y, 3, h, color.RGBA{255, 255, 255, 200}, false)
	vector.DrawFilledRect(screen, x+w-3, y, 3, h, color.RGBA{255, 255, 255, 200}, false)

	// Text (wrap simple)
	wrapped := wrapText(line, 60)
	for i, l := range wrapped {
		ebitenutil.DebugPrintAt(screen, l, int(x)+12, int(y)+12+i*14)
	}
	// Prompt
	ebitenutil.DebugPrintAt(screen, "[Space] Next  [E/Esc] Close", int(x)+12, int(y)+int(h)-18)
}

// wrapText naive word wrap.
func wrapText(s string, max int) []string {
	if len(s) <= max {
		return []string{s}
	}
	words := strings.Split(s, " ")
	var lines []string
	cur := ""
	for _, w := range words {
		if len(cur)+len(w)+1 > max {
			lines = append(lines, cur)
			cur = w
		} else {
			if cur == "" {
				cur = w
			} else {
				cur += " " + w
			}
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}
