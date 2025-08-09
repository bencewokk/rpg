package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Position structure for nodes
type Pos struct {
	X, Y float32
}

// Node represents a pathfinding node
type Node struct {
	ID  int
	Pos Pos
}

// Path represents a connection between two nodes
type Path struct {
	NodeAID int
	NodeBID int
	Cost    float32
}

// MapData holds the tile data for the map
type MapData struct {
	Width  int
	Height int
	Tiles  [][]int // 2D array of tile types
	Nodes  []Node  // Pathfinding nodes
	Paths  []Path  // Connections between nodes
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
		Nodes:  []Node{},
		Paths:  []Path{},
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

// AddNode adds a new node at the specified position
func (m *MapData) AddNode(x, y float32) int {
	id := len(m.Nodes)
	node := Node{
		ID:  id,
		Pos: Pos{X: x, Y: y},
	}
	m.Nodes = append(m.Nodes, node)
	return id
}

// RemoveNode removes a node and all paths connected to it
func (m *MapData) RemoveNode(nodeID int) {
	// Remove the node
	for i, node := range m.Nodes {
		if node.ID == nodeID {
			m.Nodes = append(m.Nodes[:i], m.Nodes[i+1:]...)
			break
		}
	}

	// Remove all paths connected to this node
	newPaths := []Path{}
	for _, path := range m.Paths {
		if path.NodeAID != nodeID && path.NodeBID != nodeID {
			newPaths = append(newPaths, path)
		}
	}
	m.Paths = newPaths
}

// AddPath adds a connection between two nodes
func (m *MapData) AddPath(nodeAID, nodeBID int) {
	// Calculate distance between nodes as cost
	var nodeA, nodeB *Node
	for i := range m.Nodes {
		if m.Nodes[i].ID == nodeAID {
			nodeA = &m.Nodes[i]
		}
		if m.Nodes[i].ID == nodeBID {
			nodeB = &m.Nodes[i]
		}
	}

	if nodeA == nil || nodeB == nil {
		return
	}

	// Check if path already exists
	for _, path := range m.Paths {
		if (path.NodeAID == nodeAID && path.NodeBID == nodeBID) ||
			(path.NodeAID == nodeBID && path.NodeBID == nodeAID) {
			return // Path already exists
		}
	}

	dx := nodeB.Pos.X - nodeA.Pos.X
	dy := nodeB.Pos.Y - nodeA.Pos.Y
	cost := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	path := Path{
		NodeAID: nodeAID,
		NodeBID: nodeBID,
		Cost:    cost,
	}
	m.Paths = append(m.Paths, path)
}

// RemovePath removes a connection between two nodes
func (m *MapData) RemovePath(nodeAID, nodeBID int) {
	newPaths := []Path{}
	for _, path := range m.Paths {
		if !((path.NodeAID == nodeAID && path.NodeBID == nodeBID) ||
			(path.NodeAID == nodeBID && path.NodeBID == nodeAID)) {
			newPaths = append(newPaths, path)
		}
	}
	m.Paths = newPaths
}

// FindNodeAt finds a node near the specified position (within radius)
func (m *MapData) FindNodeAt(x, y, radius float32) *Node {
	for i := range m.Nodes {
		node := &m.Nodes[i]
		dx := node.Pos.X - x
		dy := node.Pos.Y - y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if distance <= radius {
			return node
		}
	}
	return nil
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
