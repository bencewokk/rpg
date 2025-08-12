package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ENEMYNORMALSPEED = 70  // was 50
	ENEMYALLERTSPEED = 130 // was 100
	// Knockback tuning
	KNOCKBACK_BASE_STRENGTH = 520  // initial velocity magnitude applied on hit
	KNOCKBACK_MAX_STACK     = 780  // cap on stacked knockback velocity
	KNOCKBACK_DURATION      = 0.18 // time (s) during which knockback overrides AI movement
)

type enemy struct {
	pos     pos
	texture *ebiten.Image
	id      int

	speed    float32
	velocity float32

	offsetForAnimation int
	animationState     int

	aiState int

	inPatrol   bool
	target     pos
	route      []pos
	routeIndex int
	chasing    bool

	sleeping   bool
	sinceSleep float64

	hp       float32
	hit      bool
	sinceHit float64

	// knockback state
	knockbackVX   float32
	knockbackVY   float32
	knockbackTime float64

	// unified animation player
	animPlayer      AnimationPlayer
	currentAnimName string

	dead bool // set true once removed; prevents further interaction

	// spawning / leash data (optional)
	homePos      pos
	leashRadius  float32
	spawnerIndex int // index into runtime spawner list, -1 if not from spawner
}

func enemiesInRange(pos pos, distance float32) []*enemy {
	var es []*enemy
	for i := 0; i < len(game.currentmap.enemies); i++ {
		enemy := game.currentmap.enemies[i]
		if enemy == nil || enemy.dead || enemy.hp <= 0 {
			continue // skip dead / removed enemies
		}
		if Distance(pos, enemy.pos) < distance {
			es = append(es, enemy)
		}
	}
	return es
}

func createEnemy(pos pos) {
	var e enemy
	e.pos = pos
	e.speed = ENEMYNORMALSPEED
	e.hp = 60
	e.offsetForAnimation = rand.Intn(5)
	e.spawnerIndex = -1
	game.currentmap.enemies = append(game.currentmap.enemies, &e)
	drawables = append(drawables, &e)
}

func (e *enemy) todoEnemy() {
	if e.dead { // don't update animation/state if dead
		return
	}
	e.updateAnimation()
	e.updateState()
	e.checkHp()

}

func (e *enemy) updateAnimation() {
	if e.dead {
		return
	}
	// Decide desired animation
	desired := "idle"
	if e.animationState == 1 { // moving
		desired = "run"
	}
	if desired != e.currentAnimName {
		anim := animationManager.Get("enemy", desired)
		if anim != nil {
			e.animPlayer.SetAnimation(anim, true)
			e.currentAnimName = desired
		}
	}
	img := e.animPlayer.Update(game.deltatime)
	if img != nil {
		e.texture = img
	}
	e.animationState = 0
}

func (e *enemy) moveTowards(target pos) {
	e.animationState = 1
	dx := target.float_x - e.pos.float_x
	dy := target.float_y - e.pos.float_y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length > 0 {
		dx /= length
		dy /= length

		e.pos.float_x += dx * e.speed * float32(game.deltatime)
		e.pos.float_y += dy * e.speed * float32(game.deltatime)
	}

}

func (e *enemy) patrol() {

	e.sinceSleep += game.deltatime

	if !e.inPatrol {
		e.inPatrol = true
		e.sinceSleep = 0
		e.routeIndex = 0
		// Constrain patrol within leash if this enemy came from a spawner with leashRadius
		startNode := findClosestNode(e.pos)
		goalPos := randomPointWithinRange(startNode, 6)
		if e.leashRadius > 0 {
			allowed := nodesWithinCircle(e.homePos, e.leashRadius)
			startID := findClosestAllowedNodeID(e.pos, allowed)
			goalID := findClosestAllowedNodeID(goalPos, allowed)
			if startID >= 0 && goalID >= 0 {
				e.route = findShortestPathPositionsConstrained(startID, goalID, allowed)
			}
		}
		if len(e.route) == 0 { // fallback to original global pathing if constraints fail
			e.route = findShortestPathPositions(findClosestNode(e.pos).id, findClosestNode(goalPos).id)
		}
		e.target = e.route[e.routeIndex]

	} else {
		e.moveTowards(e.target)
	}

	if Distance(e.target, e.pos) < 30 {
		if e.routeIndex == len(e.route)-1 {
			e.inPatrol = false
			e.sinceSleep = 0
			e.routeIndex = 0
		} else {
			e.routeIndex++
			e.target = e.route[e.routeIndex]
		}
	}
}

func (e *enemy) chase() {
	nearestP := nearestCharacter(e.pos)
	e.moveTowards(nearestP.pos)
}

func (e *enemy) checkCollision(posA, posB pos) bool {
	if Distance(posA, posB) < 50 {
		return true
	}
	return false
}

func (e *enemy) hurt(c *character) {
	e.hp -= 0.3
	c.hp -= 0.1
}

func (e *enemy) checkHp() {
	if e.dead {
		return
	}
	if e.hp <= 0 {
		// Mark dead so we don't process logic any further
		e.dead = true
		// Inform spawner system if applicable
		removeEnemyFromSpawner(e)

		// Remove from enemies slice
		for i, en := range game.currentmap.enemies {
			if en == e {
				game.currentmap.enemies = append(game.currentmap.enemies[:i], game.currentmap.enemies[i+1:]...)
				break
			}
		}
		// Remove from drawables slice
		for i, d := range drawables {
			if de, ok := d.(*enemy); ok && de == e {
				drawables = append(drawables[:i], drawables[i+1:]...)
				break
			}
		}
	}
}

func (e *enemy) updateState() {
	if e.dead {
		return
	}

	e.sinceHit -= game.deltatime

	if e.sinceHit < 0 {

		e.hit = false
	}

	// Apply knockback if active (takes precedence over normal AI movement)
	if e.knockbackTime > 0 {
		dt := float32(game.deltatime)
		e.knockbackTime -= game.deltatime
		e.pos.float_x += e.knockbackVX * dt
		e.pos.float_y += e.knockbackVY * dt
		// Exponential damping for a smooth ease-out feel
		damping := float32(math.Exp(-8 * game.deltatime))
		e.knockbackVX *= damping
		e.knockbackVY *= damping
		// Flag as moving so run animation can play (optional)
		e.animationState = 1
		return
	}

	player := nearestCharacter(e.pos)
	// Absolute hard leash: if beyond 4x leash radius, force return
	if e.spawnerIndex >= 0 && e.leashRadius > 0 {
		dx := e.pos.float_x - e.homePos.float_x
		dy := e.pos.float_y - e.homePos.float_y
		distHome := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if distHome > e.leashRadius*4 {
			// run back towards home at alert speed, stop chase
			e.chasing = false
			e.animationState = 1
			if distHome > 0 {
				dx /= distHome; dy /= distHome
			}
			retSpeed := float32(ENEMYALLERTSPEED)
			e.pos.float_x -= dx * retSpeed * float32(game.deltatime)
			e.pos.float_y -= dy * retSpeed * float32(game.deltatime)
			return
		}
	}

	if Distance(e.pos, player.pos) > 100 && !e.chasing {
		e.speed = ENEMYNORMALSPEED
		nearestP, distanceToNearest := findClosestPointOnPaths(e.pos)
		switch e.aiState {
		case 0: // move towards nearest path
			if distanceToNearest > 40 {
				e.moveTowards(nearestP)
			} else {
				e.aiState = 1
			}
		case 1:
			if e.sinceSleep > 1 {
				e.patrol()
			} else {
				e.sinceSleep += game.deltatime
			}
		}
	} else {

		for i := 0; i < len(game.currentmap.players); i++ {
			c := game.currentmap.players[i]
			if checkCollision(e.pos, c.pos) {
				e.hurt(c)
			}
		}

		e.chasing = true
		e.speed = ENEMYALLERTSPEED
		e.inPatrol = false
		// Allow chase up to 4x leash if applicable; otherwise normal
		if e.spawnerIndex >= 0 && e.leashRadius > 0 {
			// Only chase if within 4x leash from home; else will be handled by early return above
			e.chase()
		} else {
			e.chase()
		}

		if Distance(e.pos, player.pos) > 400 {
			e.chasing = false
		}
	}

	// Leash handling: if spawned from spawner and too far from home, force return
	if e.spawnerIndex >= 0 && e.leashRadius > 0 {
		dx := e.pos.float_x - e.homePos.float_x
		dy := e.pos.float_y - e.homePos.float_y
		dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if dist > e.leashRadius*1.1 { // allow small overflow then pull back
			// override movement: move towards home rapidly
			// simple direct return
			e.animationState = 1
			if dist > 0 {
				dx /= dist
				dy /= dist
			}
			retSpeed := float32(ENEMYALLERTSPEED)
			e.pos.float_x -= dx * retSpeed * float32(game.deltatime)
			e.pos.float_y -= dy * retSpeed * float32(game.deltatime)
			// stop chasing once outside leash
			e.chasing = false
		}
	}
}
