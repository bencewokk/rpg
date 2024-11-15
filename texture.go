package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func parseTexture(pos pos) {
	gamemap := globalGameState.currentmap

	var (
		DDDDD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DDDD.png")
		DDDDG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DDDG.png")
		DDDGD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DDGD.png")
		DDDGG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DDGG.png")
		DDGDD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DGDD.png")
		DDGDG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DGDG.png")
		DDGGD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DGGD.png")
		DDGGG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_DGGG.png")
		DGDDD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GDDD.png")
		DGDDG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GDDG.png")
		DGDGD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GDGD.png")
		DGDGG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GDGG.png")
		DGGDD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GGDD.png")
		DGGDG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GGDG.png")
		DGGGD *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GGGD.png")
		DGGGG *ebiten.Image = loadPNG("import/tiles/Dry2Grass_GGGG.png")
	)

	var textures = map[string]*ebiten.Image{
		"DDDD": DDDDD,
		"DDDG": DDDDG,
		"DDGD": DDDGD,
		"DDGG": DDDGG,
		"DGDD": DDGDD,
		"DGDG": DDGDG,
		"DGGD": DDGGD,
		"DGGG": DDGGG,
		"GDDD": DGDDD,
		"GDDG": DGDDG,
		"GDGD": DGDGD,
		"GDGG": DGDGG,
		"GGDD": DGGDD,
		"GGDG": DGGDG,
		"GGGD": DGGGD,
		"GGGG": DGGGG,
	}

	i, j := ptid(pos)

	// for i := 0; i < globalGameState.currentmap.height; i++ {
	// 	for j := 0; j < globalGameState.currentmap.width; j++ {
	var textureID string

	// Ensure we're processing tiles of type 3
	//if gamemap.data[i][j] == 3 {

	gamemap.data[j][i] = 0

	// Check for out-of-bounds for each neighboring tile
	if i > 0 && gamemap.data[i-1][j] == 3 { // upper
		textureID += "D"
	} else {
		textureID += "G"
	}

	if j > 0 && gamemap.data[i][j-1] == 3 { // left
		textureID += "D"
	} else {
		textureID += "G"
	}

	if j < globalGameState.currentmap.width-1 && gamemap.data[i][j+1] == 3 { // right
		textureID += "D"
	} else {
		textureID += "G"
	}

	if i < globalGameState.currentmap.height-1 && gamemap.data[i+1][j] == 3 { // lower
		textureID += "D"
	} else {
		textureID += "G"
	}

	// Assign the appropriate texture based on the ID
	if texture, exists := textures[textureID]; exists {
		gamemap.texture[i][j] = texture
	}
	//}
	// 	}
	// }

}
