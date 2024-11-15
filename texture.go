package main

import "fmt"

func (gamemap *gamemap) parseTexture(pos pos) {
	//for i, j := 0, 0; i < len(gamemap.data); i++ {
	//switch gamemap.data[i][j] {

	i, j := ptid(pos)

	fmt.Println(j, i)

	// if gamemap.data[j+1][i+1] s
	// if gamemap.data[j+1][i]
	// if gamemap.data[j+1][i-1]
	// if gamemap.data[j][i-1]
	// if gamemap.data[j][i+1]
	// if gamemap.data[j-1][i+1]
	// if gamemap.data[j-1][i]
	// if gamemap.data[j-1][i-1]

	//}
	//}
}
