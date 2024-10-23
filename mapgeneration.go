package main

import (
	"fmt"
	"math/rand"
)

const (
	multipler_hillchance     = 2.4
	multipler_forestchance   = 2.1
	multipler_mountainchance = 1.1
	multipler_general        = 3
)

// Takes a float64 between 0 and 1 and returns either true or false
func calcChance(chance float64) bool {
	flip := float64(rand.Intn(100))

	if flip < chance {
		return true
	}

	return false
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

			if i != 0 && i != m.width-1 && j != 0 && j != m.height-1 {
				//on the left
				switch m.data[i][j-1] {
				case 1:
					//mountainchance *= multipler_general
				case 3:
					hillchance *= multipler_general * 2
					fmt.Println("i was here")
					fmt.Println(m.data[i][j-1])
				case 4:
					forestchance *= multipler_general * 2
					fmt.Println("i was here 2")
				}

				//on the top
				switch m.data[i-1][j] {
				case 1:
					//mountainchance *= multipler_general
				case 3:
					hillchance *= multipler_general / 3
				case 4:
					forestchance *= multipler_general / 3
				}

			}

			if i == 0 || i == m.height-1 || j == 0 || j == m.width-1 {
				m.data[i][j] = 1
			} else if calcChance(float64(forestchance)) {
				m.data[i][j] = 4 // Forest

				hillchance *= multipler_hillchance
				mountainchance *= multipler_mountainchance
				forestchance = 2

			} else if calcChance(float64(hillchance)) {
				m.data[i][j] = 3 // Hill

				forestchance *= multipler_forestchance
				mountainchance *= multipler_mountainchance
				hillchance = 2

			} else if calcChance(float64(mountainchance)) {
				m.data[i][j] = 1 // Mountain

				forestchance *= multipler_forestchance
				hillchance *= multipler_hillchance
				mountainchance = 2
			} else {
				m.data[i][j] = 2 // Plains
			}

			// // Update the chances for each iteration
			// forestchance *= multipler_forestchance
			// hillchance *= multipler_hillchance
			// mountainchance *= multipler_mountainchance

			//  else {
			// 	// Compare random value to each chance
			// 	if r < mountainchance {
			// 		m.data[i][j] = 1 // Mountain
			// 	} else if r < mountainchance+hillchance {
			// 		m.data[i][j] = 3 // Hill
			// 	} else {
			// 		m.data[i][j] = 4 // Forest
			// 	}

		}
	}

	return m
}
