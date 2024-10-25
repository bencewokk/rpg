package main

import (
	"log"

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

// Draw method of the Game
func (g *Game) Draw(screen *ebiten.Image) {
	curspos.updatemouse()

	switch globalGameState.stateid {
	case 0:

		playbtn := createButton("Play", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
		playbtn.DrawButton(screen)
		if playbtn.pressed {
			globalGameState.stateid = 3
		}

		optionsbtn := createButton("Options", 150, 50, uitransparent, uilightgray, uigray, createPos(25, 85))
		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			globalGameState.stateid = 1
		}
	case 1:
		options_exitbtn := createButton("Back to menu", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			globalGameState.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
	case 3:
		ebitenutil.DebugPrint(screen, "Test map")

		//TODO redo this comment
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
					float32(j*int(intscreendivisor)),
					float32(i*int(intscreendivisor)),
					screendivisor,
					screendivisor,
					currenttilecolor,
					false,
				)
			}
		}

		go checkmovement()
		char.DrawCharacter(screen)
	}
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
