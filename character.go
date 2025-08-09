package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARSPEED   = 260 // was 200
	DASHSPEED   = 800 // was 700
	BOOSTSPEED  = 340 // was 300
	ATTACKSPEED = 160 // character movement speed while attacking

	// Damage randomness
	MIN_DAMAGE      = 4.0
	MAX_DAMAGE      = 10.0
	CRIT_CHANCE     = 0.15
	CRIT_MULTIPLIER = 2.0
)

type character struct {
	pos     pos
	texture *ebiten.Image
	id      int

	speed              float32
	offsetForAnimation int
	animationState     int

	running     bool
	facingNorth int // between 0 and 1

	untilNewDash   float64
	untilEndOfDash float64
	dashing        bool

	untilEndOfBoost float64

	hp float32
	// UI smoothed values
	uiHp float32

	attacking                bool
	sinceAttack              float64
	attackCooldown           float64
	offsetForAnimationAttack int

	// queued fast-paced combat
	queuedAttack bool

	// New unified animation player
	animPlayer AnimationPlayer
	// cached state to decide which animation to play
	currentAnimName string
}

func createCharacter() {
	var c character

	c.hp = 100
	c.uiHp = c.hp
	c.pos = createPos(screenWidth/2, screenHeight/2)
	c.speed = CHARSPEED

	c.offsetForAnimation = rand.IntN(5)

	drawables = append(drawables, &c)
	game.currentmap.players = append(game.currentmap.players, &c)
}

func nearestCharacter(pos pos) *character {
	var closest int
	var closestDistance float32

	for i := 0; i < len(game.currentmap.players); i++ {
		if closestDistance > Distance(pos, game.currentmap.players[i].pos) {
			closestDistance = Distance(pos, game.currentmap.players[i].pos)
			closest = i
		}
	}

	return game.currentmap.players[closest]
}

func charactersInRange(pos pos, distance float32) []*character {
	var cs []*character
	for i := 0; i < len(game.currentmap.players); i++ {
		if Distance(pos, game.currentmap.players[i].pos) > distance {
			cs = append(cs, game.currentmap.players[i])
		}
	}
	return cs
}

func (c *character) updateCamera() {
	game.camera.pos = c.pos
}

func (c *character) updateAnimation() {
	// Decide desired animation key
	var desired string
	facing := "front"
	if c.facingNorth == 1 {
		facing = "back"
	}
	if c.attacking {
		desired = "attack_" + facing
	} else if c.running {
		desired = "run_" + facing
	} else {
		desired = "idle_" + facing
	}

	if desired != c.currentAnimName {
		anim := animationManager.Get("character", desired)
		if anim != nil {
			c.animPlayer.SetAnimation(anim, true)
			c.currentAnimName = desired
		}
	}

	img := c.animPlayer.Update(game.deltatime)
	if img != nil {
		c.texture = img
	}

	if c.attacking {
		c.speed = ATTACKSPEED
		if c.sinceAttack < 0 { // end of current swing
			if c.queuedAttack {
				c.queuedAttack = false
				c.attackCooldown = 0 // chain immediately
				c.attack()           // start next attack in combo
				return
			}
			c.attacking = false
			c.speed = CHARSPEED
			c.attackCooldown = 0.25 // was 0.5
		}
	}
	c.running = false // reset flagged each movement update
}

// applyKnockback sets or stacks a velocity-based knockback on an enemy.
func applyKnockback(e *enemy, from *character, strength float32) {
	dx := e.pos.float_x - from.pos.float_x
	dy := e.pos.float_y - from.pos.float_y
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if dist > 0.0001 {
		dx /= dist
		dy /= dist
	} else {
		dx, dy = 0, -1 // fallback upward
	}
	// New target velocity to add
	addVX := dx * strength
	addVY := dy * strength
	// Stack but clamp to max magnitude
	e.knockbackVX += addVX
	e.knockbackVY += addVY
	mag := float32(math.Sqrt(float64(e.knockbackVX*e.knockbackVX + e.knockbackVY*e.knockbackVY)))
	if mag > KNOCKBACK_MAX_STACK {
		scale := KNOCKBACK_MAX_STACK / mag
		e.knockbackVX *= scale
		e.knockbackVY *= scale
	}
	e.knockbackTime = KNOCKBACK_DURATION
}

func (c *character) attack() {
	c.attacking = true
	c.sinceAttack = 0.32           // was 0.52
	c.offsetForAnimationAttack = 0 // legacy field, retained for now

	es := enemiesInRange(c.pos, 80)

	for i := 0; i < len(es); i++ {
		e := es[i]
		// Roll base damage in range
		varDmg := MIN_DAMAGE + rand.Float32()*(MAX_DAMAGE-MIN_DAMAGE)
		crit := false
		// Critical hit roll
		if rand.Float32() < CRIT_CHANCE {
			varDmg *= CRIT_MULTIPLIER
			crit = true
		}
		e.hp -= float32(varDmg)
		e.hit = true
		e.sinceHit = 0.2
		applyKnockback(e, c, KNOCKBACK_BASE_STRENGTH)
		AddDamageIndicator(e.pos, float32(varDmg), crit)
	}

}

func (c *character) checkHp() {
	if c.hp < 1 {
		removeAtID(c.id, drawables)
	}
}

func (c *character) todoCharacter() {
	c.checkHp()
	c.updateCamera()
	c.checkMovement()
	c.updateAnimation()
}
