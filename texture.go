package main

import "fmt"

func (gamemap *gamemap) parseTexture(pos pos) {
	//for i, j := 0, 0; i < len(gamemap.data); i++ {
	//switch gamemap.data[i][j] {

	i, j := ptid(pos)

	fmt.Println(j, i)

	var textureID string = ""
	if gamemap.data[j-1][i+1] == 2 { // upper right
		textureID += "G"
	} // upper right
	// if gamemap.data[j-1][i] // upper
	// gamemap.data[j-1][i-1] = 4 // upper leftÍ

	// 0 = not decided, 1 = mountains, 2 = plains, 3 = dry
	// if gamemap.data[j+1][i+1]  //lower right
	// if gamemap.data[j+1][i] // lower
	// gamemap.data[j+1][i-1] = 4 // lower left
	// if gamemap.data[j][i-1] // left
	// if gamemap.data[j][i+1] // right

	fmt.Println(textureID)
	//}
	//}
}
