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

	readMapData()
	parseTextureAndSprites()

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	createCharacter()
	createEnemy(createPos(500, 500))

	loadChar()
	fmt.Println(enemyAnimations)
	loadEnemy()
	fmt.Println(enemyAnimations)

	screendivisor = 30
	intscreendivisor = 30

	game.camera.zoom = 1

}

type gamemap struct {
	// map data (2D array)
	//
	// 0 = not decided, 1 = mountains, 2 = plains, 3 = dry
	data    [100][150]int
	texture [100][150]*ebiten.Image

	// height of the map
	//
	//used for rendering and generating the map
	height int
	width  int
}

// read more in gamestate
type camera struct {
	pos pos

	//used in rendering and collision checking
	zoom float32
}

func offsetsx(tobeoffset float32) float32 {
	return ((tobeoffset-game.camera.pos.float_x)*game.camera.zoom + screenWidth/2)
}
func offsetsy(tobeoffset float32) float32 {
	return ((tobeoffset-game.camera.pos.float_y)*game.camera.zoom + screenHeight/2)

}

var game Game

type Game struct {
	// 0 menu / 1 menu and options / 3 in game
	stateid int

	// maps are stored in arrays (see in type map)
	//
	// this is the current map that is  being used//while rendered map array size is constant to 144 (16*9) currentmapid is not
	currentmap gamemap

	// counts the time since start of game
	//
	// get updated every frame
	deltatime float64

	// date of last update
	lastUpdateTime time.Time

	// contains the camera positions
	//
	// this is used in the rendering, it offsets the drawing positions
	camera camera
}

// Update method of the Game
func (g *Game) Update() error {

	curspos.updatemouse()
	go checkZoom()

	updateAnimationCycle()

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
	game.deltatime = now.Sub(game.lastUpdateTime).Seconds()
	game.lastUpdateTime = now

	switch game.stateid {
	case 0:

		playbtn.DrawButton(screen)
		if playbtn.pressed {
			game.stateid = 3
		}

		optionsbtn.DrawButton(screen)
		if optionsbtn.pressed {
			game.stateid = 1
		}

	case 1:

		options_exitbtn.DrawButton(screen)
		if options_exitbtn.pressed {
			game.stateid = 0
		}

		vector.DrawFilledRect(screen, 200, 25, screenWidth-250, screenHeight-50, uidarkgray, false)
		testslider.DrawSlider(screen)

	case 3:
		sortDrawables()

		for i := 0; i < game.currentmap.height; i++ {
			for j := 0; j < game.currentmap.width; j++ {
				if game.currentmap.texture[i][j] != nil {
					drawTile(screen, game.currentmap.texture[i][j], i, j)
				}
			}
		}

		for i := 0; i < len(drawables); i++ {
			drawables[i].draw(screen)
		}

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
