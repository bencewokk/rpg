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
	readMapData()
	parseTextureAndSprites()

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	screendivisor = 30
	intscreendivisor = 30

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
				if globalGameState.currentmap.texture[i][j] != nil {

					drawTile(screen, globalGameState.currentmap.texture[i][j], i, j)

				}

			}

		}

		for i := 0; i < len(globalGameState.currentmap.sprites); i++ {
			//fmt.Println(len(currentmap.sprites))
			// sort.Slice(globalGameState.currentmap.sprites, func(i, j int) bool {
			// 	return globalGameState.currentmap.sprites[i].pos.float_y < globalGameState.currentmap.sprites[j].pos.float_y
			// })
			drawSprite(screen, globalGameState.currentmap.sprites[i].texture, globalGameState.currentmap.sprites[i].pos)

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
