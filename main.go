package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Screen sizes
var (
	width, height = ebiten.Monitor().Size()
	screenWidth   = float32(width)
	screenHeight  = float32(height)

	screendivisor    float32
	intscreendivisor int
)

// Audio
var (
	audioCtx           *audio.Context
	menuMusicPlayer    *audio.Player
	prevMenuMusicState string
)

func gameinit() {
	rand.Seed(time.Now().UnixNano())

	readMapData()
	parseTextureAndSprites()

	// Initialize new animation system
	animationManager = NewAnimationManager()
	if err := animationManager.LoadManifest("import/animations.json"); err != nil {
		fmt.Println("Animation manifest load failed:", err)
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("rpg")

	createCharacter()
	createEnemy(createPos(500, 500))
	createEnemy(createPos(700, 500))
	createEnemy(createPos(500, 400))
	createEnemy(createPos(400, 900))

	// Sample NPC (after enemies so player can interact)
	createNPC(createPos(600, 600), []string{
		"Hey there adventurer!",
		"Nice day to wander the plains, isn't it?",
		"Come back later, I might have a quest for you.",
	})
	createEnemy(createPos(500, 500))
	createEnemy(createPos(700, 500))
	createEnemy(createPos(500, 400))
	createEnemy(createPos(400, 900))
	createEnemy(createPos(500, 500))
	createEnemy(createPos(700, 500))
	createEnemy(createPos(500, 400))
	createEnemy(createPos(400, 900))

	screendivisor = 30
	intscreendivisor = 30

	game.camera.zoom = 1

	// Initialize audio context & load menu music
	audioCtx = audio.NewContext(44100)
	musicPath := "music/rpg main theme.wav"
	if info, statErr := os.Stat(musicPath); statErr != nil {
		log.Println("[AUDIO] music file not found:", musicPath, statErr)
	} else {
		log.Println("[AUDIO] music file found:", musicPath, "size=", info.Size())
		f, err := os.Open(musicPath)
		if err != nil {
			log.Println("[AUDIO] open error:", err)
		} else {
			data, rerr := io.ReadAll(f)
			f.Close()
			if rerr != nil {
				log.Println("[AUDIO] read error:", rerr)
			} else {
				log.Println("[AUDIO] read into memory bytes=", len(data))
				r := bytes.NewReader(data)
				stream, err := wav.DecodeWithSampleRate(44100, r)
				if err != nil {
					log.Println("[AUDIO] decode error:", err)
				} else {
					loop := audio.NewInfiniteLoop(stream, stream.Length())
					p, err := audioCtx.NewPlayer(loop)
					if err != nil {
						log.Println("[AUDIO] player create error:", err)
					} else {
						p.SetVolume(0.6)
						menuMusicPlayer = p
						menuMusicPlayer.Play()
						prevMenuMusicState = "playing"
						log.Println("[AUDIO] menu music started (looping, vol=0.6)")
					}
				}
			}
		}
	}

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
	// non-hostile talkable NPCs
	npcs []*npc

	paths []path
	nodes []node

	players []*character
	enemies []*enemy
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
	go checkZoom()

	// Update cursor position
	cx, cy := ebiten.CursorPosition()
	curspos.float_x = float32(cx)
	curspos.float_y = float32(cy)

	// Update buttons (hover/click detection) before using their pressed flags
	playbtn.UpdateButton()
	optionsbtn.UpdateButton()
	options_exitbtn.UpdateButton()
	exitbtn.UpdateButton()

	if optionsbtn.pressed {
		game.stateid = 1
	}
	if options_exitbtn.pressed {
		game.stateid = 0
	}
	if playbtn.pressed { game.stateid = 3 }
	if exitbtn.pressed {
		fmt.Println("exited with code 0")
		os.Exit(0)
	}

	// ESC: from game -> menu; from menu/options -> exit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch game.stateid {
		case 3:
			game.stateid = 0
		case 0, 1:
			fmt.Println("exited with code 0")
			os.Exit(0)
		}
	}

	// Reset one-shot pressed state so it doesn't trigger repeatedly
	playbtn.pressed = false
	optionsbtn.pressed = false
	options_exitbtn.pressed = false
	exitbtn.pressed = false

	// (Animation cycle handled per AnimationPlayer now)

	// Damage indicators
	updateDamageIndicators(game.deltatime)

	// NPC interaction handling (only ingame)
	if game.stateid == 3 {
		updateNPCInteractions()
		updateNPCAnimations(game.deltatime)
	}

	// Menu music state management
	if menuMusicPlayer != nil {
		desiredPlaying := (game.stateid == 0 || game.stateid == 1)
		currentlyPlaying := menuMusicPlayer.IsPlaying()
		if desiredPlaying && !currentlyPlaying {
			menuMusicPlayer.Play()
		} else if !desiredPlaying && currentlyPlaying {
			menuMusicPlayer.Pause()
		}
		// Debug state change log
		stateNow := "paused"
		if menuMusicPlayer.IsPlaying() {
			stateNow = "playing"
		}
		if stateNow != prevMenuMusicState {
			log.Printf("[AUDIO] menu music state -> %s (stateid=%d)\n", stateNow, game.stateid)
			prevMenuMusicState = stateNow
		}
	}

	return nil
}

var (
	// Case 0
	playbtn    = createButton("Play", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
	optionsbtn = createButton("Options", 150, 50, uitransparent, uilightgray, uigray, createPos(25, 85))
	exitbtn    = createButton("Exit", 150, 50, uitransparent, uilightgray, uigray, createPos(25, 145))

	// Case 1
	options_exitbtn = createButton("Back to menu", 150, 50, uitransparent, uilightgray, uigray, onearg_createPos(25))
	testslider      = createSlider("testslider", 500, 20, 5, 10, uigray, uilightgray, uigray, createPos(230, 80))
)

var screenGlobal *ebiten.Image

// Draw method of the Game
func (g *Game) Draw(screen *ebiten.Image) {

	screenGlobal = screen

	now := time.Now()
	game.deltatime = now.Sub(game.lastUpdateTime).Seconds()
	game.lastUpdateTime = now

	// ESC handling moved to Update for state-aware behavior

	switch game.stateid {
	case 0:

		playbtn.DrawButton(screen)

		optionsbtn.DrawButton(screen)
		exitbtn.DrawButton(screen)

	case 1:

		options_exitbtn.DrawButton(screen)

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
			drawables[i].giveId(i)
			drawables[i].draw(screen)
		}

		// for i := 0; i < len(game.currentmap.paths); i++ {
		// 	drawPath(screen, game.currentmap.paths[i])
		// }

		// for i := 0; i < len(game.currentmap.nodes); i++ {
		// 	n := game.currentmap.nodes[i]
		// 	ebitenutil.DebugPrintAt(screen, strconv.Itoa(n.id), int(offsetsx(n.pos.float_x)), int(offsetsy(n.pos.float_y)))
		// }

		p := 0
		game.currentmap.players[p].drawUi()

		// Draw floating damage after entities so it's on top
		drawDamageIndicators()
		// Draw conversation if active
		drawConversationUI(screen)

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
