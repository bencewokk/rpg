package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// accumulate time for UI pulses
var pulseTime float64

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
	// Widen panel to accommodate separate area for dash circle
	panelW := float32(230)
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
	// Reserve space on the right for the circular dash widget (radius 18 + margins)
	circleRadius := float32(18)
	circleMargin := float32(14) // space between bar end and circle edge
	reserved := circleRadius*2 + circleMargin
	barW := panelW - 24 - reserved
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

	// Dash cooldown circular widget using proper vector paths
	cx := panelX + panelW - (circleRadius + 8) // center inside reserved area
	cy := panelY + 26
	radius := circleRadius
	// Thicker ring rendering
	ringThickness := float32(8)
	if ringThickness > radius-2 {
		ringThickness = radius - 2
	}
	// Base ring background
	drawFilledCircle(screenGlobal, cx, cy, radius, color.RGBA{25, 25, 32, 200})
	drawFilledCircle(screenGlobal, cx, cy, radius-ringThickness, color.RGBA{30, 30, 40, 255}) // carve inner hole
	if c.untilNewDash > 0 {
		remaining := clampFloat(float32(c.untilNewDash/1.5), 0, 1)
		drawRingArc(screenGlobal, cx, cy, radius-1, ringThickness, remaining, color.RGBA{120, 180, 255, 240})
		// inner core background
		innerCoreR := radius - ringThickness - 2
		if innerCoreR > 4 {
			drawFilledCircle(screenGlobal, cx, cy, innerCoreR, color.RGBA{35, 35, 45, 255})
		}
	} else {
		pulseTime += game.deltatime
		pulse := float32(0.6 + 0.1*math.Sin(pulseTime*4))
		col := color.RGBA{uint8(60 + 30*pulse), uint8(200 + 40*pulse), uint8(110 + 30*pulse), 255}
		readyR := radius - ringThickness + 2
		if readyR < 4 {
			readyR = radius / 2
		}
		drawFilledCircle(screenGlobal, cx, cy, readyR, col)
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
// Proper filled circle using vector path
func drawFilledCircle(dst *ebiten.Image, cx, cy, r float32, col color.Color) {
	// Approximate filled circle with triangle fan
	steps := int(32 + r*0.5)
	if steps < 12 {
		steps = 12
	}
	cx64 := float64(cx)
	cy64 := float64(cy)
	for i := 0; i < steps; i++ {
		a0 := float64(i) * 2 * math.Pi / float64(steps)
		a1 := float64(i+1) * 2 * math.Pi / float64(steps)
		x0 := cx64 + float64(r)*math.Cos(a0)
		y0 := cy64 + float64(r)*math.Sin(a0)
		x1 := cx64 + float64(r)*math.Cos(a1)
		y1 := cy64 + float64(r)*math.Sin(a1)
		// Using tiny rectangles to approximate triangles (cheap but acceptable for small radii)
		vector.DrawFilledRect(dst, float32(x0), float32(y0), 1, 1, col, false)
		vector.DrawFilledRect(dst, float32(x1), float32(y1), 1, 1, col, false)
		vector.DrawFilledRect(dst, cx, cy, 1, 1, col, false)
	}
}

// Filled pie wedge (0..pct of full circle clockwise starting at angle 0 (pointing right))
// drawRingArc draws a ring sector from angle 0 to pct*2PI with given outer radius and thickness
func drawRingArc(dst *ebiten.Image, cx, cy, outerR, thickness, pct float32, col color.Color) {
	if pct <= 0 {
		return
	}
	if pct > 1 {
		pct = 1
	}
	innerR := outerR - thickness
	if innerR < 0 {
		innerR = 0
	}
	steps := int(64 * pct)
	if steps < 8 {
		steps = 8
	}
	total := float64(pct) * 2 * math.Pi
	cx64 := float64(cx)
	cy64 := float64(cy)
	for i := 0; i < steps; i++ {
		a0 := total * float64(i) / float64(steps)
		a1 := total * float64(i+1) / float64(steps)
		// For each segment, approximate quad between inner/outer radii
		ox0 := cx64 + float64(outerR)*math.Cos(a0)
		oy0 := cy64 + float64(outerR)*math.Sin(a0)
		ox1 := cx64 + float64(outerR)*math.Cos(a1)
		oy1 := cy64 + float64(outerR)*math.Sin(a1)
		ix0 := cx64 + float64(innerR)*math.Cos(a0)
		iy0 := cy64 + float64(innerR)*math.Sin(a0)
		ix1 := cx64 + float64(innerR)*math.Cos(a1)
		iy1 := cy64 + float64(innerR)*math.Sin(a1)
		// Rasterize by drawing small rects at the four vertices
		vector.DrawFilledRect(dst, float32(ox0), float32(oy0), 1, 1, col, false)
		vector.DrawFilledRect(dst, float32(ox1), float32(oy1), 1, 1, col, false)
		vector.DrawFilledRect(dst, float32(ix0), float32(iy0), 1, 1, col, false)
		vector.DrawFilledRect(dst, float32(ix1), float32(iy1), 1, 1, col, false)
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
