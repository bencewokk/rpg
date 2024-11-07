package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func gameinit() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	globalGameState.currentmap = createMap(36)

	screendivisor = screenHeight / float32(globalGameState.currentmap.height)
	intscreendivisor = int(screenHeight) / globalGameState.currentmap.height

	char.pos.float_y = 90
	char.pos.float_x = 90
}

// Screen sizes
var (
	width, height = ebiten.Monitor().Size()
	screenWidth   = float32(width)
	screenHeight  = float32(height)

	screendivisor    float32
	intscreendivisor int
)

type Game struct{}

// Update method of the Game
func (g *Game) Update() error {
	return nil
}

var (
	// Case 0
	playbtn    = createButton("Play", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
	optionsbtn = createButton("Options", 150, 50, uitransparent, uilightgray, uigray, createPos(25, 85))

	// Case 1
	options_exitbtn = createButton("Back to menu", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
	testslider      = createSlider("testslider", 500, 20, 5, 10, uigray, uilightgray, uigray, createPos(230, 80))
)

// Draw method of the Game
func (g *Game) Draw(screen *ebiten.Image) {

	camera := globalGameState.camera

	now := time.Now()
	globalGameState.deltatime = now.Sub(globalGameState.lastUpdateTime).Seconds()
	globalGameState.lastUpdateTime = now

	curspos.updatemouse()

	switch globalGameState.stateid {
	case 0:

		playbtn.DrawButton(screen)
		if playbtn.pressed {
			globalGameState.stateid = 3
		}

		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			globalGameState.stateid = 1
		}

	case 1:

		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			globalGameState.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
		testslider.DrawSlider(screen)

	case 3:

		//TODO redo this comment and make this into a function
		for i := 0; i < globalGameState.currentmap.height; i++ {
			for j := 0; j < globalGameState.currentmap.width; j++ {
				switch globalGameState.currentmap.data[i][j] {
				case 2:
					currenttilecolor = mlightgreen
				case 3:
					currenttilecolor = mbrown
				case 1:
					currenttilecolor = mdarkgray
				case 4:
					currenttilecolor = mdarkgreen
				}
				vector.DrawFilledRect(
					screen,
					(float32(j*intscreendivisor)+camera.pos.float_x)*2,
					(float32(i*intscreendivisor)+camera.pos.float_y)*2,
					screendivisor*2,
					screendivisor*2,
					currenttilecolor,
					false,
				)
			}
		}

		char.DrawCharacter(screen)
		checkMovement()

		for i := 0; i < len(enemies); i++ {
			enemies[i].Draw(screen)
		}

		// var op *ebiten.DrawImageOptions
		// screen.DrawImage(cutCam(screen, createPos(30, 30)), op)

	}

	fps := ebiten.CurrentFPS()
	fpsText := fmt.Sprintf("FPS: %.2f", fps)
	ebitenutil.DebugPrint(screen, fpsText)

}

// Layout method of the Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	gameinit()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
