package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type gamemap struct {
	// map data (2D array)
	//
	// 0 = not decided, 1 = mountains, 2 = plains, 3 = dry
	data    [200][150]int
	texture [200][150]*ebiten.Image

	// height of the map
	//
	//used for rendering and generating the map
	height int
	width  int
}

// read more in gamestate
type camera struct {
	pos pos

	//used in rendering and collision checking
	zoom float32
}

// contains all global information about the game
var globalGameState gamestate

// contains all global information about the game
//
// contains maps
type gamestate struct {
	// 0 menu / 1 menu and options / 3 in game
	stateid int

	// maps are stored in arrays (see in type map)
	//
	// this is the current map that is  being used//while rendered map array size is constant to 144 (16*9) currentmapid is not
	currentmap gamemap

	// counts the time since start of game
	//
	// get updated every frame
	deltatime float64

	// date of last update
	lastUpdateTime time.Time

	// contains the camera positions
	//
	// this is used in the rendering, it offsets the drawing positions
	camera camera
}

func updateCamera() {
	globalGameState.camera.pos.float_x = char.pos.float_x * -1
	globalGameState.camera.pos.float_y = char.pos.float_y * -1
}

// Function to load map data from a file and set width/height in currentmap
func readMapData() {
	filename := "map.txt"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	y := 0
	maxWidth := 0

	// First pass to determine the maximum width of the map
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		values := strings.Split(line, ",")

		rowWidth := len(values)
		if rowWidth > maxWidth {
			maxWidth = rowWidth
		}
		y++
	}

	// Set the map dimensions based on the number of rows (y) and maximum width
	globalGameState.currentmap.width = maxWidth
	globalGameState.currentmap.height = y

	// Now read the data into the map
	for i, line := range lines {
		values := strings.Split(line, ",")
		for x, value := range values {
			// Trim leading and trailing whitespace before parsing
			value = strings.TrimSpace(value)

			// Skip empty strings
			if value == "" {
				continue
			}

			intValue, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Error parsing integer:", err)
				return
			}
			// Use the current index (i) and x to fill the map
			globalGameState.currentmap.data[i][x] = intValue
		}

	}

	// Check for errors after scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	fmt.Println("Map data loaded successfully!")
	fmt.Printf("Map dimensions: %d x %d\n", globalGameState.currentmap.width, globalGameState.currentmap.height)
}
