package main

import (
	"math"
	"math/rand"
	"rpg/mapio"
	"time"
)

// runtimeSpawner augments mapio.EnemySpawner with timing & tracking
type runtimeSpawner struct {
	data       mapio.EnemySpawner
	timer      float64
	alive      map[*enemy]struct{}
	nextJitter float64 // randomized offset to desync spawns
}

var spawners []*runtimeSpawner

func initSpawners(m *mapio.MapData) {
	spawners = []*runtimeSpawner{}
	for _, sp := range m.Spawners {
		rs := &runtimeSpawner{data: sp, alive: make(map[*enemy]struct{})}
		rs.nextJitter = rand.Float64() * float64(sp.IntervalSeconds)
		spawners = append(spawners, rs)
	}
}

// call once per frame while in-game (state 3)
func updateSpawners(dt float64) {
	// Clean references to dead enemies
	for _, rs := range spawners {
		for e := range rs.alive {
			if e.dead || e.hp <= 0 { // removed
				delete(rs.alive, e)
			}
		}
	}
	for idx, rs := range spawners {
		rs.timer += dt
		interval := float64(rs.data.IntervalSeconds)
		if interval <= 0 {
			interval = 1
		}
		if len(rs.alive) < rs.data.MaxAlive && rs.timer >= interval+rs.nextJitter {
			spawnEnemyFromSpawner(idx, rs)
			rs.timer = 0
			rs.nextJitter = 0 // after first spawn
		}
	}
}

func spawnEnemyFromSpawner(index int, rs *runtimeSpawner) {
	// Random point within circle (uniform)
	u := rand.Float64()
	r := math.Sqrt(u) * float64(rs.data.Radius)
	theta := rand.Float64() * math.Pi * 2
	x := rs.data.Pos.X + float32(r*math.Cos(theta))
	y := rs.data.Pos.Y + float32(r*math.Sin(theta))
	epos := createPos(x, y)
	// Snap to nearest path point if path network present and within radius to keep enemies on routes
	if len(game.currentmap.paths) > 0 {
		closest, dist := findClosestPointOnPaths(epos)
		if dist < rs.data.Radius*1.1 { // allow small slack
			// Only snap if the closest path point is not wildly outside spawn circle
			epos = closest
		}
	}
	createEnemy(epos)
	e := game.currentmap.enemies[len(game.currentmap.enemies)-1]
	e.homePos = createPos(rs.data.Pos.X, rs.data.Pos.Y)
	e.leashRadius = rs.data.Radius
	e.spawnerIndex = index
	rs.alive[e] = struct{}{}
}

func removeEnemyFromSpawner(e *enemy) {
	if e.spawnerIndex < 0 || e.spawnerIndex >= len(spawners) {
		return
	}
	rs := spawners[e.spawnerIndex]
	delete(rs.alive, e)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
