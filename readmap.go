package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func readMapSprites() {
	filename := "map.txt"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isReadingSprites := false
	globalGameState.currentmap.sprites = nil // Clear existing sprites

	y := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Look for the sprite section header
		if line == "---SPRITES---" {
			isReadingSprites = true
			continue
		}

		// Process sprite data when in the sprite section
		if isReadingSprites {
			// Split sprite data by commas
			values := strings.Split(line, ",")
			if len(values) != 3 {
				fmt.Println("Invalid sprite data:", line)
				continue
			}

			// Parse the values for the sprite
			typeOf, err := strconv.Atoi(strings.TrimSpace(values[0]))
			if err != nil {
				fmt.Println("Error parsing sprite type:", err)
				continue
			}

			floatX, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 32)
			if err != nil {
				fmt.Println("Error parsing sprite X position:", err)
				continue
			}

			floatY, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 32)
			if err != nil {
				fmt.Println("Error parsing sprite Y position:", err)
				continue
			}

			// Create the sprite and add it to the map
			sprite := createSprite(createPos(float32(floatX), float32(floatY)), typeOf)
			globalGameState.currentmap.sprites = append(globalGameState.currentmap.sprites, sprite)
		} else {
			y++
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	fmt.Println("File read successfully!")
}
