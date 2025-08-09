package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type damage struct {
	pos pos

	wanterFor  float64
	sinceDrawn float64
	damage     int
}

func (c *character) drawUi() {
	// Smooth HP animation (lerp)
	target := c.hp
	diff := target - c.uiHp
	c.uiHp += diff * float32(math.Min(1, game.deltatime*8)) // respond quickly but smooth

	// Panel background
	panelW := float32(160)
	panelH := float32(70)
	panelX := float32(20)
	panelY := float32(20)
	// Shadow
	vector.DrawFilledRect(screenGlobal, panelX+3, panelY+3, panelW, panelH, color.RGBA{0, 0, 0, 80}, false)
	// Main panel
	vector.DrawFilledRect(screenGlobal, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 200}, false)
	// Top highlight strip
	vector.DrawFilledRect(screenGlobal, panelX, panelY, panelW, 4, color.RGBA{90, 90, 120, 255}, false)

	// HP bar container
	barX := panelX + 12
	barY := panelY + 14
	barW := panelW - 24
	barH := float32(14)
	// Background
	vector.DrawFilledRect(screenGlobal, barX, barY, barW, barH, color.RGBA{50, 50, 60, 255}, false)
	// Gradient fill based on smoothed uiHp (assuming max 100)
	pct := clampFloat(c.uiHp/100, 0, 1)
	fillW := barW * pct
	// Draw segmented gradient (simple 4 segments to fake gradient)
	segments := 4
	for i := 0; i < segments; i++ {
		segPct0 := float32(i) / float32(segments)
		segPct1 := float32(i+1) / float32(segments)
		if segPct1 > pct { // partial segment
			segPct1 = pct
		}
		if segPct1 <= segPct0 {
			break
		}
		// Color interpolate red -> orange -> yellow -> lime
		prog := (segPct0 + segPct1) * 0.5
		col := gradient4(prog)
		segX := barX + barW*segPct0
		segW := barW * (segPct1 - segPct0)
		vector.DrawFilledRect(screenGlobal, segX, barY, segW, barH, col, false)
	}
	// HP loss overlay (recent damage flash)
	if diff < -0.1 { // took damage
		vector.DrawFilledRect(screenGlobal, barX+fillW, barY, barW-fillW, barH, color.RGBA{200, 30, 30, 80}, false)
	}

	// Outline
	drawRectStroke(screenGlobal, barX, barY, barW, barH, color.RGBA{15, 15, 20, 255})

	// Dash cooldown circular bar (mini) on right of panel
	cx := panelX + panelW - 26
	cy := panelY + 26
	radius := float32(18)
	// Background circle
	drawCircleFilled(screenGlobal, cx, cy, radius, color.RGBA{40, 40, 50, 220})
	// Cooldown fill (radial sweep)
	if c.untilNewDash > 0 {
		pctCd := clampFloat(float32(c.untilNewDash/1.5), 0, 1)
		drawArcFilled(screenGlobal, cx, cy, radius-3, pctCd, color.RGBA{120, 180, 255, 230})
	} else {
		// Ready state glow
		drawCircleFilled(screenGlobal, cx, cy, radius-3, color.RGBA{90, 255, 140, 230})
	}
}

// Helper for gradient coloring 0..1 across 4 key colors
func gradient4(t float32) color.RGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	// key colors
	keys := []color.RGBA{
		{230, 40, 30, 255},  // red
		{255, 120, 30, 255}, // orange
		{255, 220, 40, 255}, // yellow
		{90, 255, 90, 255},  // lime
	}
	seg := float32(len(keys) - 1)
	f := t * seg
	i := int(f)
	if i >= len(keys)-1 {
		return keys[len(keys)-1]
	}
	local := f - float32(i)
	a := keys[i]
	b := keys[i+1]
	return color.RGBA{
		uint8(float32(a.R) + (float32(b.R)-float32(a.R))*local),
		uint8(float32(a.G) + (float32(b.G)-float32(a.G))*local),
		uint8(float32(a.B) + (float32(b.B)-float32(a.B))*local),
		255,
	}
}

// drawRectStroke draws a 1px outline
func drawRectStroke(dst *ebiten.Image, x, y, w, h float32, c color.Color) {
	vector.DrawFilledRect(dst, x, y, w, 1, c, false)
	vector.DrawFilledRect(dst, x, y+h-1, w, 1, c, false)
	vector.DrawFilledRect(dst, x, y, 1, h, c, false)
	vector.DrawFilledRect(dst, x+w-1, y, 1, h, c, false)
}

// drawCircleFilled approximates a filled circle
func drawCircleFilled(dst *ebiten.Image, cx, cy, r float32, col color.Color) {
	steps := int(20 + r*0.8)
	for i := 0; i < steps; i++ {
		a0 := float64(i) * 2 * math.Pi / float64(steps)
		a1 := float64(i+1) * 2 * math.Pi / float64(steps)
		x0 := cx + float32(math.Cos(a0))*r
		y0 := cy + float32(math.Sin(a0))*r
		x1 := cx + float32(math.Cos(a1))*r
		y1 := cy + float32(math.Sin(a1))*r
		// triangle fan using thin quads (approx)
		vector.DrawFilledRect(dst, x0, y0, 1+(x1-x0), 1+(y1-y0), col, false)
	}
}

// drawArcFilled draws a simple pie arc from 0..pct (0..1)
func drawArcFilled(dst *ebiten.Image, cx, cy, r, pct float32, col color.Color) {
	steps := int(60 * pct)
	if steps < 1 {
		return
	}
	for i := 0; i < steps; i++ {
		a := float64(i) * 2 * math.Pi / 60.0
		x := cx + float32(math.Cos(a))*r
		y := cy + float32(math.Sin(a))*r
		vector.DrawFilledRect(dst, x, y, 2, 2, col, false)
	}
}

func clampFloat(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
