package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type DamageIndicator struct {
	pos         pos
	amount      float32
	lifetime    float64
	maxLifetime float64
	vy          float32
	col         color.RGBA
	crit        bool
}

var damageIndicators []*DamageIndicator

func AddDamageIndicator(p pos, amount float32, crit bool) {
	di := &DamageIndicator{
		pos:         p,
		amount:      amount,
		lifetime:    0,
		maxLifetime: 0.6, // shorter display
		vy:          50, // faster upward
		col:         color.RGBA{255, 0, 0, 255},
		crit:        crit,
	}
	damageIndicators = append(damageIndicators, di)
}

func updateDamageIndicators(dt float64) {
	write := 0
	for _, di := range damageIndicators {
		di.lifetime += dt
		di.pos.float_y -= di.vy * float32(dt)
		if di.lifetime < di.maxLifetime {
			damageIndicators[write] = di
			write++
		}
	}
	damageIndicators = damageIndicators[:write]
}

func drawDamageIndicators() {
	face := basicfont.Face7x13
	for _, di := range damageIndicators {
		progress := di.lifetime / di.maxLifetime
		if progress < 0 {
			progress = 0
		}
		if progress > 1 {
			progress = 1
		}
		// Pulsing scale (sin wave) with crit amplification
		baseScale := 1.0 + 0.4*math.Sin(progress*math.Pi)
		scale := baseScale
		if di.crit {
			scale *= 1.4
		}
		// Flash stronger early on
		flash := 1.0
		if di.lifetime < 0.25 {
			if int(di.lifetime*50)%2 == 0 { // toggle quickly
				flash = 1.5
			}
		}
		// Fade toward end
		alpha := 1.0 - progress
		if alpha < 0 {
			alpha = 0
		}
		baseR, baseG, baseB := 255.0, 30.0, 30.0
		if di.crit { // yellow for crit
			baseR, baseG, baseB = 255, 220, 40
		}
		r := uint8(clampInt(int(baseR*flash), 0, 255))
		g := uint8(clampInt(int(baseG*flash), 0, 255))
		b := uint8(clampInt(int(baseB*flash), 0, 255))
		a := uint8(255 * alpha)
		prefix := ""
		if di.crit {
			prefix = "★"
		}
		txt := fmt.Sprintf("%s%.0f", prefix, di.amount)
		screenX := float64(offsetsx(di.pos.float_x))
		screenY := float64(offsetsy(di.pos.float_y))
		// Prepare text image
		w := len(txt)*7 + 4
		h := 13 + 4
		tmp := ebiten.NewImage(w, h)
		outline := color.RGBA{0, 0, 0, a}
		for ox := -1; ox <= 1; ox++ {
			for oy := -1; oy <= 1; oy++ {
				if ox == 0 && oy == 0 {
					continue
				}
				text.Draw(tmp, txt, face, 2+ox, 2+oy+10, outline)
			}
		}
		text.Draw(tmp, txt, face, 2, 2+10, color.RGBA{r, g, b, a})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(screenX-float64(w)*scale/2, screenY-float64(h)*scale/2)
		screenGlobal.DrawImage(tmp, op)
	}
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
