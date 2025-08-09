package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// updateMenuEffects kept for compatibility (now a no-op)
func updateMenuEffects(dt float64) {}

// drawFancyMenu: simplified flat background & basic title/buttons
func drawFancyMenu(screen *ebiten.Image, state int) {
	// Flat background
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{30, 70, 50, 255}, false)

	// Title
	title := "RPG"
	ebitenutil.DebugPrintAt(screen, title, int(screenWidth/2)-20+1, 41+1)
	ebitenutil.DebugPrintAt(screen, title, int(screenWidth/2)-20, 41)
	subtitle := "Demo"
	ebitenutil.DebugPrintAt(screen, subtitle, int(screenWidth/2)-24, 60)

	switch state {
	case 0:
		playbtn.DrawButton(screen)
		optionsbtn.DrawButton(screen)
		exitbtn.DrawButton(screen)
	case 1:
		options_exitbtn.DrawButton(screen)
		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
		testslider.DrawSlider(screen)
	}
}
