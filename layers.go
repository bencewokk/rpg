package main

import (
	"fmt"
	"sort"
)

type drawable struct {
	character *character
	sprite    *sprite
}

var drawables []drawable

func addDrawable[T character | sprite](item T, index int) {
	var d drawable
	drawables = append(drawables, d)

	switch v := any(item).(type) {
	case character:
		drawables[index].character = &v
	case sprite:
		drawables[index].sprite = &v
	default:
		fmt.Println("Unsupported type", v)
	}
}

func addAllToDrawables() {
	i := 0
	for ; i < len(globalGameState.currentmap.sprites); i++ {
		addDrawable(globalGameState.currentmap.sprites[i], i)
	}

	addDrawable(char, i)
}

func sortDrawablesByY() {
	sort.Slice(drawables, func(i, j int) bool {
		var y1, y2 float32
		var offset float32 = 46

		if drawables[i].character != nil {
			y1 = drawables[i].character.pos.float_y

		} else if drawables[i].sprite != nil {
			y1 = drawables[i].sprite.pos.float_y + offset
		}

		if drawables[j].character != nil {
			y2 = drawables[j].character.pos.float_y
		} else if drawables[j].sprite != nil {
			y2 = drawables[j].sprite.pos.float_y + offset
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
