package main

import (
	"math/rand"
)

const (
	multipler_hillchance     = 1.1
	multipler_forestchance   = 1.1
	multipler_mountainchance = 1.005
	multipler_general_top    = 60
	multipler_general_bottom = 1.3
)

// Returns either true or false
func calcChance(chance float64) bool {
	flip := float64(rand.Intn(100))
	return flip < chance
}

func createMap(_height int) gamemap {
	var (
		hillchance     float32 = 1.0
		forestchance   float32 = 1.0
		mountainchance float32 = 1.0

		m gamemap
	)

	m.height = _height
	m.width = int(float32(_height) * 1.77777777778)

	// 0 = not decided, 1 = mountains, 2 = plains, 3 = hills, 4 = forests
	// Edges are mountains
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {

			if i != 0 && i != m.height-1 && j != 0 && j != m.width-1 {

				//on the top
				switch m.data[i][j-1] {
				case 3:
					hillchance *= multipler_general_top
				case 4:
					forestchance *= multipler_general_top
				}

				//on the bottom
				switch m.data[i-1][j] {
				case 3:
					hillchance *= multipler_general_bottom
				case 4:
					forestchance *= multipler_general_bottom
				}

			}

			if i == 0 || i == m.height-1 || j == 0 || j == m.width-1 {
				m.data[i][j] = 1
			} else if calcChance(float64(forestchance / 2)) {
				m.data[i][j] = 4 // Forest

				hillchance *= multipler_hillchance
				mountainchance *= multipler_mountainchance
				forestchance = 1

			} else if calcChance(float64(hillchance) / 6) {
				m.data[i][j] = 3 // Hill

				forestchance *= multipler_forestchance
				mountainchance *= multipler_mountainchance
				hillchance = 1

			} else if calcChance(float64(mountainchance)) {
				m.data[i][j] = 1 // Mountain

				forestchance *= multipler_forestchance
				hillchance *= multipler_hillchance
				mountainchance = 1
			} else {
				m.data[i][j] = 2 // Plains
			}
		}
	}

	return m
}
