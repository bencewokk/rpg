package main

import (
	"math/rand"
)

// contains how should each globan chance variable be incremented
const (
	hillchance     = 1.3
	forestchance   = 1.3
	mountainchance = 1.1
)

// contains chances what each tile type has to be created
type tile struct {
	tag int

	plainchance    int
	hillchance     int
	forestchance   int
	mountainchance int
}

// creates a map with a 16:9 aspect ratio
//
// return is a 2d array
//
// 0 not decided, 1 mountains, 2 plains, 3 hills, 4 forests
//
// first global chance is checked wheter we should add a hill or a forest
func createMap() gamemap {
	var m gamemap
	var seed int64
	for i := 0; i < 9; i++ {
		for j := 0; j < 16; j++ {
			rand.Seed(seed)
			a := rand.Int()
			seed++
			// 0 not decided, 1 mountains, 2 plains, 3 hills, 4 forests
			//
			//edges are mountains

			if i == 0 || i == 8 || j == 0 || j == 15 {
				m.data[i][j] = 1
			} else {
				switch a % 3 {
				case 0:
					m.data[i][j] = 0
				case 1:
					m.data[i][j] = 1
				case 2:
					m.data[i][j] = 3
				}
			}

			//
		}
	}

	return m
}
