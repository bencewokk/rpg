package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

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

var dryTransitionsTextures = map[string]*ebiten.Image{
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

func parseTextureAndSprites() {

	var (
		Grass_S1 *ebiten.Image = loadPNG("import/tiles/Grass_S1.png")
		Grass_S2 *ebiten.Image = loadPNG("import/tiles/Grass_S2.png")
		Grass_S3 *ebiten.Image = loadPNG("import/tiles/Grass_S3.png")
		Grass_S6 *ebiten.Image = loadPNG("import/tiles/Grass_S6.png")
		Grass_S8 *ebiten.Image = loadPNG("import/tiles/Grass_S8.png")

		Grass_S4 *ebiten.Image = loadPNG("import/tiles/Grass_S4.png") // normal
		Grass_S5 *ebiten.Image = loadPNG("import/tiles/Grass_S5.png") // normal
		Grass_S7 *ebiten.Image = loadPNG("import/tiles/Grass_S7.png") // normal
	)

	var grassTextures = []*ebiten.Image{
		Grass_S1, Grass_S2, Grass_S3, Grass_S6, Grass_S8, Grass_S4, Grass_S5, Grass_S7,
	}

	for i := 0; i < game.currentmap.height; i++ {
		for j := 0; j < game.currentmap.width; j++ {

			var textureID string = ""

			if game.currentmap.data[i][j] == 3 {
				// Check for out-of-bounds for each neighboring tile
				if game.currentmap.data[i-1][j] == 3 { // upper
					textureID += "D"
				} else {
					textureID += "G"
				}

				if game.currentmap.data[i][j-1] == 3 { // left
					textureID += "D"
				} else {
					textureID += "G"
				}

				if game.currentmap.data[i][j+1] == 3 { // right
					textureID += "D"
				} else {
					textureID += "G"
				}

				if game.currentmap.data[i+1][j] == 3 { // lower
					textureID += "D"
				} else {
					textureID += "G"
				}

				// Assign the appropriate texture based on the ID
				if texture, exists := dryTransitionsTextures[textureID]; exists {
					game.currentmap.texture[i][j] = texture
				}
			} else if game.currentmap.data[i][j] == 2 {
				if calcChance(10) {
					game.currentmap.texture[i][j] = grassTextures[rand.Int31n(5)]
				} else {
					game.currentmap.texture[i][j] = grassTextures[rand.Int31n(3)+5]
				}
			}
		}
	}
}

var trees = []*ebiten.Image{tree_S1, tree_S2}

var (
	tree_S1 *ebiten.Image = loadPNG("import/prop/tree1.png")
	tree_S2 *ebiten.Image = loadPNG("import/prop/tree2.png")
)

var characterAnimations [4][]*ebiten.Image

func loadChar() {

	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle.png"))
	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S2.png"))
	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S3.png"))
	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S4.png"))
	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S5.png"))
	characterAnimations[0] = append(characterAnimations[0], loadPNG("import/Characters/Character/Front_C_Idle_S6.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S2.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S3.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S4.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S5.png"))
	characterAnimations[1] = append(characterAnimations[1], loadPNG("import/Characters/Character/Back_C_Idle_S6.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S2.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S3.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S4.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S5.png"))
	characterAnimations[2] = append(characterAnimations[2], loadPNG("import/Characters/Character/Front_C_Running_S6.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S2.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S3.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S4.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S5.png"))
	characterAnimations[3] = append(characterAnimations[3], loadPNG("import/Characters/Character/Back_C_Running_S6.png"))
}

var enemyAnimations [2][6]*ebiten.Image

func loadEnemy() {
	enemyAnimations[0][0] = loadPNG("import/Characters/Enemy/enemyidle1.png")
	enemyAnimations[0][1] = loadPNG("import/Characters/Enemy/enemyidle2.png")
	enemyAnimations[0][2] = loadPNG("import/Characters/Enemy/enemyidle3.png")
	enemyAnimations[0][3] = loadPNG("import/Characters/Enemy/enemyidle4.png")
	enemyAnimations[0][4] = loadPNG("import/Characters/Enemy/enemyidle1.png")
	enemyAnimations[0][5] = loadPNG("import/Characters/Enemy/enemyidle2.png")
	enemyAnimations[0][5] = loadPNG("import/Characters/Enemy/enemyidle2.png")

	enemyAnimations[1][0] = loadPNG("import/Characters/Enemy/enemyrun1.png")
	enemyAnimations[1][1] = loadPNG("import/Characters/Enemy/enemyrun2.png")
	enemyAnimations[1][2] = loadPNG("import/Characters/Enemy/enemyrun3.png")
	enemyAnimations[1][3] = loadPNG("import/Characters/Enemy/enemyrun4.png")
	enemyAnimations[1][4] = loadPNG("import/Characters/Enemy/enemyrun1.png")
	enemyAnimations[1][5] = loadPNG("import/Characters/Enemy/enemyrun2.png")
	enemyAnimations[1][5] = loadPNG("import/Characters/Enemy/enemyrun2.png")
}

var (
	animationTimer float64
	animationCycle int
)

func updateAnimationCycle() {
	animationTimer += game.deltatime

	if animationTimer > 0.13 {
		animationCycle++
		animationTimer = 0
	}
}
