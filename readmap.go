package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readMapData() {
	filename := "map.txt"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isReadingSprites := false

	y := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Look for the sprite section header
		if line == "---SPRITES---" {
			isReadingSprites = true
			continue
		}

		// Process sprite data
		if isReadingSprites {
			// 	// Split sprite data by commas
			// 	values := strings.Split(line, ",")
			// 	if len(values) != 3 {
			// 		fmt.Println("Invalid sprite data:", line)
			// 		continue
			// 	}

			// 	// Parse the values for the sprite
			// 	typeOf, err := strconv.Atoi(strings.TrimSpace(values[0]))
			// 	if err != nil {
			// 		fmt.Println("Error parsing sprite type:", err)
			// 		continue
			// 	}

			// 	floatX, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 32)
			// 	if err != nil {
			// 		fmt.Println("Error parsing sprite X position:", err)
			// 		continue
			// 	}

			// 	floatY, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 32)
			// 	if err != nil {
			// 		fmt.Println("Error parsing sprite Y position:", err)
			// 		continue
			// 	}

			// 	// Create the sprite and add it to the map
			// 	sprite := createSprite(createPos(float32(floatX), float32(floatY)), typeOf)
			// 	game.currentmap.sprites = append(game.currentmap.sprites, sprite)
		} else {
			// Process map data
			values := strings.Split(line, ",")
			for x, value := range values {
				value = strings.TrimSpace(value)
				if value == "" {
					continue
				}

				intValue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println("Error parsing map value:", err)
					return
				}

				game.currentmap.data[y][x] = intValue
			}
			y++
		}
	}

	game.currentmap.height = 100
	game.currentmap.width = 150

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	fmt.Println("File read successfully!")
}
