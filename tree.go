package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type tree struct {
	typeOf  int
	pos     pos
	texture *ebiten.Image
	treeId  int
}

func createSprite(pos pos, typeOf int) {
	var t tree

	switch typeOf {
	case 0: // tree
		t.typeOf = 0
		t.treeId = rand.Intn(len(trees))
		t.texture = trees[t.treeId]
		t.pos = pos

		drawables = append(drawables, &t)
	}
}
