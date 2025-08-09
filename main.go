package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// Audio
var (
	audioCtx           *audio.Context
	menuMusicPlayer    *audio.Player
	prevMenuMusicState string
	baseMusicVolume    = 0.6
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
	stateid   int
	prevState int

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
	resumeBtn.UpdateButton()
	pauseMenuBtn.UpdateButton()
	pauseExitBtn.UpdateButton()

	if optionsbtn.pressed {
		game.stateid = 1
	}
	if options_exitbtn.pressed {
		game.stateid = 0
	}
	if playbtn.pressed {
		game.stateid = 3
	}
	if exitbtn.pressed {
		fmt.Println("exited with code 0")
		os.Exit(0)
	}
	// Pause overlay buttons (state 2)
	if resumeBtn.pressed && game.stateid == 2 {
		game.stateid = 3
	}
	if pauseMenuBtn.pressed && game.stateid == 2 {
		game.stateid = 0
	}
	if pauseExitBtn.pressed && game.stateid == 2 {
		fmt.Println("exited with code 0")
		os.Exit(0)
	}

	// ESC behavior with pause state
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch game.stateid {
		case 3: // in game -> pause
			game.stateid = 2
		case 2: // pause -> back to game
			game.stateid = 3
		case 0, 1: // menus -> exit
			fmt.Println("exited with code 0")
			os.Exit(0)
		}
	}
	// Also allow 'P' to toggle pause in-game
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if game.stateid == 3 {
			game.stateid = 2
		} else if game.stateid == 2 {
			game.stateid = 3
		}
	}

	// Reset one-shot pressed state so it doesn't trigger repeatedly
	playbtn.pressed = false
	optionsbtn.pressed = false
	options_exitbtn.pressed = false
	exitbtn.pressed = false
	resumeBtn.pressed = false
	pauseMenuBtn.pressed = false
	pauseExitBtn.pressed = false

	// (Animation cycle handled per AnimationPlayer now)

	// Damage indicators
	updateDamageIndicators(game.deltatime)

	// Menu particle/visual effects
	if game.stateid == 0 || game.stateid == 1 {
		updateMenuEffects(game.deltatime)
	}

	// NPC interaction handling (only ingame, not paused)
	if game.stateid == 3 {
		updateNPCInteractions()
		updateNPCAnimations(game.deltatime)
	}

	// Music volume management (keeps music always playing; lowers on pause)
	updateMusicVolume(game.stateid, game.deltatime)
	game.prevState = game.stateid

	// Always keep music playing (started in init). Volume handled above.

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

	// Pause (state 2)
	resumeBtn    = createButton("Resume", 150, 45, uitransparent, uilightgray, uigray, createPos(50, 60))
	pauseMenuBtn = createButton("Menu", 150, 45, uitransparent, uilightgray, uigray, createPos(50, 115))
	pauseExitBtn = createButton("Exit", 150, 45, uitransparent, uilightgray, uigray, createPos(50, 170))
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
		drawFancyMenu(screen, 0)

	case 1:
		drawFancyMenu(screen, 1)

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

	case 2: // paused overlay: draw game scene behind then overlay
		// First draw game world (reuse logic from case 3 without changing state)
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
		p := 0
		game.currentmap.players[p].drawUi()
		// damage + conversations on top
		drawDamageIndicators()
		drawConversationUI(screen)
		// Washed overlay (desaturated feel via tinted semi-transparent layer)
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{40, 40, 40, 170}, false)
		// Pause panel
		panelW := float32(260)
		panelH := float32(200)
		panelX := (screenWidth - panelW) / 2
		panelY := (screenHeight - panelH) / 2
		vector.DrawFilledRect(screen, panelX+4, panelY+4, panelW, panelH, color.RGBA{0, 0, 0, 120}, false) // shadow
		vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{70, 80, 75, 230}, false)
		// Title
		fmtStr := "PAUSED"
		fbx := int(panelX + (panelW-float32(len(fmtStr))*7)/2)
		fby := int(panelY + 10)
		ebitenutil.DebugPrintAt(screen, fmtStr, fbx, fby)
		// Reposition pause buttons relative to panel for neat layout
		resumeBtn.pos.float_x = panelX + panelW/2 - resumeBtn.width/2
		resumeBtn.pos.float_y = panelY + 50
		pauseMenuBtn.pos.float_x = panelX + panelW/2 - pauseMenuBtn.width/2
		pauseMenuBtn.pos.float_y = panelY + 100
		pauseExitBtn.pos.float_x = panelX + panelW/2 - pauseExitBtn.width/2
		pauseExitBtn.pos.float_y = panelY + 150
		resumeBtn.DrawButton(screen)
		pauseMenuBtn.DrawButton(screen)
		pauseExitBtn.DrawButton(screen)

	}

	fps := ebiten.CurrentFPS()
	fpsText := fmt.Sprintf("FPS: %.2f", fps)
	ebitenutil.DebugPrint(screen, fpsText)
}

// Layout method of the Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// managePauseAudio handles quick fade out/in when pausing/unpausing.
func managePauseAudio(prev, cur int, dt float64) {
	if menuMusicPlayer == nil {
		return
	}
	// Only affect if we were/are in menu music playing states (0,1) OR paused state
	if prev == cur {
		return
	}
	// If entering pause from menu (0 or 1) we fast fade
	if (cur == 2) && (prev == 0 || prev == 1) {
		// reduce volume rapidly
		v := menuMusicPlayer.Volume()
		v -= float64(dt) * 2.5
		if v < 0 {
			v = 0
		}
		menuMusicPlayer.SetVolume(v)
		if v == 0 {
			menuMusicPlayer.Pause()
		}
	}
	// If leaving pause back to menu, restore volume quickly and play
	if (prev == 2) && (cur == 0 || cur == 1) {
		if !menuMusicPlayer.IsPlaying() {
			menuMusicPlayer.Play()
		}
		v := menuMusicPlayer.Volume()
		v += float64(dt) * 3
		if v > 0.6 {
			v = 0.6
		}
		menuMusicPlayer.SetVolume(v)
	}
}

// updateMusicVolume keeps music playing; when paused (state 2) it squashes volume lower and restores otherwise.
func updateMusicVolume(state int, dt float64) {
	if menuMusicPlayer == nil {
		return
	}
	// ensure playing
	if !menuMusicPlayer.IsPlaying() {
		menuMusicPlayer.Play()
	}
	target := baseMusicVolume
	if state == 2 { // paused -> squash
		target = baseMusicVolume * 0.25
	}
	cur := menuMusicPlayer.Volume()
	// smooth approach
	speed := 2.5 // volume units per second
	if math.Abs(cur-target) < 0.01 {
		menuMusicPlayer.SetVolume(target)
		return
	}
	if cur < target {
		cur += speed * dt
	} else {
		cur -= speed * dt
	}
	if cur < 0 {
		cur = 0
	}
	if cur > baseMusicVolume {
		cur = baseMusicVolume
	}
	menuMusicPlayer.SetVolume(cur)
}

func main() {
	gameinit()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
