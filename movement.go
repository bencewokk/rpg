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
	ptid(char.pos)
	for {
		time.Sleep(10 * time.Millisecond)
		if !checkNextTile() {
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				char.pos.float_x += 0.2
			}

			if ebiten.IsKeyPressed(ebiten.KeyA) {
				char.pos.float_x -= 0.2
			}

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				char.pos.float_y -= 0.2
			}

			if ebiten.IsKeyPressed(ebiten.KeyS) {
				char.pos.float_y += 0.2
			}
		}
	}
}
