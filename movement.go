package main

//
import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func checkNextTile() bool {

	return false
}

func checkmovement() {

	for {
		time.Sleep(2500000)
		if !checkNextTile() {
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				char.pos.float_x += 0.3
			}

			if ebiten.IsKeyPressed(ebiten.KeyA) {
				char.pos.float_x -= 0.3
			}

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				char.pos.float_y -= 0.3
			}

			if ebiten.IsKeyPressed(ebiten.KeyS) {
				char.pos.float_y += 0.3
			}
		}
	}
}
