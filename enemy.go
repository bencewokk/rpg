package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

var enemyAnimations [1][6]*ebiten.Image
var globalAC int

func loadEnemy() {
	enemyAnimations[0][0] = loadPNG("import/Characters/Enemy/enemyidle1.png")
	enemyAnimations[0][1] = loadPNG("import/Characters/Enemy/enemyidle2.png")
	enemyAnimations[0][2] = loadPNG("import/Characters/Enemy/enemyidle3.png")
	enemyAnimations[0][3] = loadPNG("import/Characters/Enemy/enemyidle4.png")
	enemyAnimations[0][4] = loadPNG("import/Characters/Enemy/enemyidle1.png")
	enemyAnimations[0][5] = loadPNG("import/Characters/Enemy/enemyidle2.png")
	enemyAnimations[0][5] = loadPNG("import/Characters/Enemy/enemyidle2.png")
}

// Contains all information about the enemies
type enemy struct {
	id          int
	title       string
	pos         pos
	curtiletype int
	hp          int

	animationCycle int
}

var (
	enemies []enemy
)

// Returns a new enemy with the given title and path to the picture
func createEnemy(title string, id int) enemy {
	var e enemy
	e.title = title
	e.id = id

	e.animationCycle = int(rand.Int32N(5))

	return e
}

func drawEnemy(screen *ebiten.Image, e enemy) {
	op := &ebiten.DrawImageOptions{}
	t := enemyAnimations[0][(e.animationCycle+globalAC)%6]
	// Set up scaling
	originalWidth, originalHeight := t.Size()
	scaleX := float64(screendivisor) / float64(originalWidth) * float64(game.camera.zoom)
	scaleY := float64(screendivisor) / float64(originalHeight) * float64(game.camera.zoom)
	op.GeoM.Scale(scaleX, scaleY)

	// Positioning with respect to camera
	op.GeoM.Translate(
		float64(offsetsx(e.pos.float_y)),
		float64(offsetsy(e.pos.float_x)),
	)

	// Draw the selected portion of the image onto the screen
	screen.DrawImage(t, op)
}

var globalAnimationTimer float64

func updateAnimationEnemies() {
	globalAnimationTimer += game.deltatime
	if globalAnimationTimer >= 0.45 {
		globalAC++
		globalAnimationTimer = 0.0
	}

}

func (e *enemy) Die() {
	e.pos.float_y = screenHeight / 2
	e.pos.float_x = screenWidth / 2

	e.hp = 100
}

// Assuming 'enemy' is of type 'character' or has a position that you can access
func (e *enemy) Hurt(enemyPos pos) {
	e.hp -= 10
	if e.hp <= 0 {
		e.Die()
	}

	// Calculate the direction to move away from the enemy
	moveAmount := float32(30) // Amount to move away
	directionX := e.pos.float_x - enemyPos.float_x
	directionY := e.pos.float_y - enemyPos.float_y

	// Normalize the direction vector
	length := float32(math.Sqrt(float64(directionX*directionX+directionY*directionY))) * 2.5
	if length > 0 {
		directionX /= length
		directionY /= length
	}

	// Move the character away from the enemy
	e.pos.float_x += directionX * moveAmount
	e.pos.float_y += directionY * moveAmount
}

func init() {
	enemies = append(enemies, createEnemy("Enemy 1", 0))
	enemies[0].pos = createPos(400, 400)
	enemies = append(enemies, createEnemy("Enemy 2", 1))
	enemies[1].pos = createPos(300, 300)
	enemies = append(enemies, createEnemy("Enemy 3", 2))
	enemies[2].pos = createPos(700, 700)
}
