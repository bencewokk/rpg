package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MapData holds the tile data for the map
type MapData struct {
	Width  int
	Height int
	Tiles  [][]int // 2D array of tile types
}

// NewMapData creates a new map with the specified dimensions
func NewMapData(width, height int) MapData {
	tiles := make([][]int, height)
	for i := range tiles {
		tiles[i] = make([]int, width)
	}
	
	return MapData{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

// GetTile returns the tile type at the specified coordinates
func (m *MapData) GetTile(x, y int) int {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return 0 // Return void for out-of-bounds
	}
	return m.Tiles[y][x]
}

// SetTile sets the tile type at the specified coordinates
func (m *MapData) SetTile(x, y, tileType int) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}
	m.Tiles[y][x] = tileType
}

// LoadFromFile loads map data from a file (compatible with your RPG format)
func (m *MapData) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	y := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for section headers and skip them for now
		if strings.HasPrefix(line, "---") {
			break // Stop at first section header (sprites, nodes, etc.)
		}
		
		// Parse tile data
		values := strings.Split(line, ",")
		for x, value := range values {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}
			
			intValue, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			
			if x < m.Width && y < m.Height {
				m.Tiles[y][x] = intValue
			}
		}
		y++
		
		if y >= m.Height {
			break
		}
	}
	
	fmt.Printf("Loaded map: %dx%d\n", m.Width, m.Height)
	return nil
}

// SaveToFile saves the map data to a file (compatible with your RPG format)
func (m *MapData) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write map data
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if x > 0 {
				file.WriteString(", ")
			}
			file.WriteString(strconv.Itoa(m.Tiles[y][x]))
		}
		file.WriteString("\n")
	}
	
	// TODO: Later we'll add sprite, node, and path data here
	
	fmt.Printf("Saved map to %s\n", filename)
	return nil
}
