package main

import (
	"fmt"
	"sort"
)

type drawable struct {
	character *character
	sprite    *sprite
	enemy     *enemy
}

var drawables []drawable

func addDrawable[T character | sprite | enemy](item T, index int) {
	var d drawable
	drawables = append(drawables, d)

	switch v := any(item).(type) {
	case character:
		drawables[index].character = &v
	case sprite:
		drawables[index].sprite = &v
	case enemy:
		drawables[index].enemy = &v
	default:
		fmt.Println("Unsupported type", v)
	}
}

func addAllToDrawables() {
	u := 0
	for i := 0; i < len(game.currentmap.sprites); i++ {
		addDrawable(game.currentmap.sprites[i], u)
		u++
	}
	for i := 0; i < len(enemies); i++ {
		addDrawable(enemies[i], u)
		u++
	}

	addDrawable(char, u)
}

func sortDrawablesByY() {
	sort.Slice(drawables, func(i, j int) bool {
		var y1, y2 float32
		var offset float32 = 46

		if drawables[i].character != nil {
			y1 = drawables[i].character.pos.float_y
		} else if drawables[i].sprite != nil {
			y1 = drawables[i].sprite.pos.float_y + offset
		} else if drawables[i].enemy != nil {
			y1 = drawables[i].enemy.pos.float_y + offset
		}

		if drawables[j].character != nil {
			y2 = drawables[j].character.pos.float_y
		} else if drawables[j].sprite != nil {
			y2 = drawables[j].sprite.pos.float_y + offset
		} else if drawables[i].enemy != nil {
			y2 = drawables[i].enemy.pos.float_y + offset
		}

		return y1 < y2
	})
}

func updatePositions() {
	for i := range drawables {
		if drawables[i].character != nil {
			drawables[i].character.pos.float_y = char.pos.float_y - 50
		}
	}
}
