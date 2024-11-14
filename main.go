package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Screen sizes
var (
	width, height = ebiten.Monitor().Size()
	screenWidth   = float32(width)
	screenHeight  = float32(height)

	screendivisor    float32
	intscreendivisor int
)

func gameinit() {

	load()

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	globalGameState.currentmap = createMap(36)
	// 1080/36--30
	screendivisor = screenHeight / float32(globalGameState.currentmap.height)
	intscreendivisor = int(screenHeight) / globalGameState.currentmap.height

	char.pos.float_y = screenHeight / 2
	char.pos.float_x = screenWidth / 2

	globalGameState.camera.zoom = 1

}

type Game struct {
}

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

		updateCamera()

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
					(float32(j*intscreendivisor-intscreendivisor/2)+camera.pos.float_x)*camera.zoom+screenWidth/2,
					(float32(i*intscreendivisor-intscreendivisor/2)+camera.pos.float_y)*camera.zoom+screenHeight/2,
					screendivisor*camera.zoom,
					screendivisor*camera.zoom,
					currenttilecolor,
					false,
				)

				// posX := (float32(j*intscreendivisor-intscreendivisor/2)+camera.pos.float_x)*camera.zoom + screenWidth/2
				// posY := (float32(i*intscreendivisor-intscreendivisor/2)+camera.pos.float_y)*camera.zoom + screenHeight/2

				// var a string = strconv.Itoa(j) + " " + strconv.Itoa(i)

				// ebitenutil.DebugPrintAt(screen, "K", int(posX), int(posY))
			}

		}

		char.DrawCharacter(screen)
		checkMovementAndInput()
		updateAnimationCharacter()

		for i := 0; i < len(enemies); i++ {
			enemies[i].Draw(screen)
		}

		// var op *ebiten.DrawImageOptions
		// screen.DrawImage(cutCam(screen, createPos(30, 30)), op)

		drawUi(screen)
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
