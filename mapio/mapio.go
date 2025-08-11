package mapio

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Position struct for map data
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

// Sprite represents an object placed on the map
type Sprite struct {
	Type int
	Pos  Pos
}

// MapData contains all map information
type MapData struct {
	Width   int
	Height  int
	Tiles   [][]int
	Nodes   []Node
	Paths   []Path
	Sprites []Sprite
	NPCs    []NPC
}

// NPC represents a placed NPC with dialogue. VoiceKey reserved for future voice integration.
type NPC struct {
	Name       string
	Pos        Pos
	Dialogues  []string
	VoiceKey   string // placeholder for future audio key/asset id
	SpritePath string
}

// NewMapData creates a new empty map with specified dimensions
func NewMapData(width, height int) *MapData {
	tiles := make([][]int, height)
	for y := range tiles {
		tiles[y] = make([]int, width)
	}

	return &MapData{
		Width:   width,
		Height:  height,
		Tiles:   tiles,
		Nodes:   []Node{},
		Paths:   []Path{},
		Sprites: []Sprite{},
		NPCs:    []NPC{},
	}
}

// LoadMapFromFile reads map data from a text file
func LoadMapFromFile(filename string) (*MapData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	mapData := &MapData{
		Nodes:   []Node{},
		Paths:   []Path{},
		Sprites: []Sprite{},
		NPCs:    []NPC{},
	}

	scanner := bufio.NewScanner(file)
	isReadingSprites := false
	isReadingNodes := false
	isReadingPaths := false
	isReadingNPCs := false

	y := 0
	var maxWidth int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Look for section headers
		switch line {
		case "---SPRITES---":
			isReadingSprites = true
			isReadingNodes = false
			isReadingPaths = false
			isReadingNPCs = false
			continue
		case "---NODES---":
			isReadingSprites = false
			isReadingNodes = true
			isReadingPaths = false
			isReadingNPCs = false
			continue
		case "---PATHS---":
			isReadingSprites = false
			isReadingNodes = false
			isReadingPaths = true
			isReadingNPCs = false
			continue
		case "---NPCS---":
			isReadingSprites = false
			isReadingNodes = false
			isReadingPaths = false
			isReadingNPCs = true
			continue
		}

		// Process data based on current section
		if isReadingSprites {
			sprite, err := parseSpriteLine(line)
			if err != nil {
				fmt.Printf("Warning: Invalid sprite data: %s\n", line)
				continue
			}
			mapData.Sprites = append(mapData.Sprites, *sprite)
		} else if isReadingNodes {
			node, err := parseNodeLine(line)
			if err != nil {
				fmt.Printf("Warning: Invalid node data: %s\n", line)
				continue
			}
			mapData.Nodes = append(mapData.Nodes, *node)
		} else if isReadingPaths {
			path, err := parsePathLine(line)
			if err != nil {
				// Paths might have different formats, so we're more lenient
				continue
			}
			mapData.Paths = append(mapData.Paths, *path)
		} else if isReadingNPCs {
			npc, err := parseNPCLine(line)
			if err != nil {
				fmt.Printf("Warning: Invalid NPC data: %s\n", line)
				continue
			}
			mapData.NPCs = append(mapData.NPCs, *npc)

		} else {
			// Process map tile data
			if mapData.Tiles == nil {
				// Initialize tiles on first row
				mapData.Tiles = [][]int{}
			}

			values := strings.Split(line, ",")
			row := make([]int, len(values))

			for x, value := range values {
				value = strings.TrimSpace(value)
				if value == "" {
					continue
				}

				intValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing map value at (%d,%d): %v", x, y, err)
				}
				row[x] = intValue
			}

			mapData.Tiles = append(mapData.Tiles, row)
			if len(row) > maxWidth {
				maxWidth = len(row)
			}
			y++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	mapData.Width = maxWidth
	mapData.Height = y

	fmt.Printf("Loaded map: %dx%d with %d nodes, %d paths, %d sprites, %d NPCs\n",
		mapData.Width, mapData.Height, len(mapData.Nodes), len(mapData.Paths), len(mapData.Sprites), len(mapData.NPCs))

	return mapData, nil
}

// SaveMapToFile writes map data to a text file
func SaveMapToFile(mapData *MapData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write map tiles
	for y := 0; y < mapData.Height; y++ {
		row := make([]string, mapData.Width)
		for x := 0; x < mapData.Width; x++ {
			if y < len(mapData.Tiles) && x < len(mapData.Tiles[y]) {
				row[x] = strconv.Itoa(mapData.Tiles[y][x])
			} else {
				row[x] = "0"
			}
		}
		writer.WriteString(strings.Join(row, ", ") + "\n")
	}

	// Write sprites section
	if len(mapData.Sprites) > 0 {
		writer.WriteString("---SPRITES---\n")
		for _, sprite := range mapData.Sprites {
			writer.WriteString(fmt.Sprintf("%d, %.1f, %.1f\n", sprite.Type, sprite.Pos.X, sprite.Pos.Y))
		}
	}

	// Write NPCs section (format: NPC, Name, X, Y, VoiceKeyPlaceholder, dialogue1|dialogue2|...)
	if len(mapData.NPCs) > 0 {
		writer.WriteString("---NPCS---\n")
		for _, n := range mapData.NPCs {
			voice := n.VoiceKey
			if voice == "" {
				voice = "-"
			}
			joined := strings.Join(n.Dialogues, "|")
			path := n.SpritePath
			if path == "" {
				path = "-"
			}
			writer.WriteString(fmt.Sprintf("NPC, %s, %.1f, %.1f, %s, %s, %s\n", n.Name, n.Pos.X, n.Pos.Y, voice, path, joined))
		}
	}

	// Write nodes section
	if len(mapData.Nodes) > 0 {
		writer.WriteString("---NODES---\n")
		for _, node := range mapData.Nodes {
			writer.WriteString(fmt.Sprintf("NODE, %d, %.1f, %.1f\n", node.ID, node.Pos.X, node.Pos.Y))
		}
	}

	// Write paths section
	if len(mapData.Paths) > 0 {
		writer.WriteString("---PATHS---\n")
		for _, path := range mapData.Paths {
			writer.WriteString(fmt.Sprintf("PATH, %d, %d, %.1f\n", path.NodeAID, path.NodeBID, path.Cost))
		}
	}

	return nil
}

// Helper functions for parsing lines

func parseSpriteLine(line string) (*Sprite, error) {
	values := strings.Split(line, ",")
	if len(values) != 3 {
		return nil, fmt.Errorf("invalid sprite format")
	}

	typeOf, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid sprite type: %v", err)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid sprite X: %v", err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid sprite Y: %v", err)
	}

	return &Sprite{
		Type: typeOf,
		Pos:  Pos{X: float32(x), Y: float32(y)},
	}, nil
}

func parseNodeLine(line string) (*Node, error) {
	values := strings.Split(line, ",")
	if len(values) != 4 {
		return nil, fmt.Errorf("invalid node format")
	}

	if strings.TrimSpace(values[0]) != "NODE" {
		return nil, fmt.Errorf("not a node line")
	}

	id, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %v", err)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid node X: %v", err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(values[3]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid node Y: %v", err)
	}

	return &Node{
		ID:  id,
		Pos: Pos{X: float32(x), Y: float32(y)},
	}, nil
}

func parsePathLine(line string) (*Path, error) {
	values := strings.Split(line, ",")
	if len(values) != 4 {
		return nil, fmt.Errorf("invalid path format")
	}

	if strings.TrimSpace(values[0]) != "PATH" {
		return nil, fmt.Errorf("not a path line")
	}

	nodeAID, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid nodeA ID: %v", err)
	}

	nodeBID, err := strconv.Atoi(strings.TrimSpace(values[2]))
	if err != nil {
		return nil, fmt.Errorf("invalid nodeB ID: %v", err)
	}

	cost, err := strconv.ParseFloat(strings.TrimSpace(values[3]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid cost: %v", err)
	}

	return &Path{
		NodeAID: nodeAID,
		NodeBID: nodeBID,
		Cost:    float32(cost),
	}, nil
}

// parseNPCLine parses an NPC line of format:
// NPC, Name, X, Y, VoiceKey, dialogue1|dialogue2|...
// VoiceKey may be '-' placeholder.
func parseNPCLine(line string) (*NPC, error) {
	values := strings.Split(line, ",")
	if len(values) < 7 {
		return nil, fmt.Errorf("invalid NPC format")
	}
	if strings.TrimSpace(values[0]) != "NPC" {
		return nil, fmt.Errorf("not an NPC line")
	}
	name := strings.TrimSpace(values[1])
	x, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid NPC X")
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(values[3]), 32)
	if err != nil {
		return nil, fmt.Errorf("invalid NPC Y")
	}
	voiceKey := strings.TrimSpace(values[4])
	if voiceKey == "-" {
		voiceKey = ""
	}
	spritePath := strings.TrimSpace(values[5])
	if spritePath == "-" {
		spritePath = ""
	}
	dialogueField := strings.Join(values[6:], ",")
	dialogueField = strings.TrimSpace(dialogueField)
	dialogues := []string{}
	if dialogueField != "" {
		dialogues = strings.Split(dialogueField, "|")
	}
	for i, d := range dialogues {
		dialogues[i] = strings.TrimSpace(d)
	}
	return &NPC{Name: name, Pos: Pos{X: float32(x), Y: float32(y)}, Dialogues: dialogues, VoiceKey: voiceKey, SpritePath: spritePath}, nil
}

// GetTile safely gets a tile value at the specified coordinates
func (m *MapData) GetTile(x, y int) int {
	if y < 0 || y >= len(m.Tiles) || x < 0 || x >= len(m.Tiles[y]) {
		return 0 // Return void/default tile for out-of-bounds
	}
	return m.Tiles[y][x]
}

// SetTile safely sets a tile value at the specified coordinates
func (m *MapData) SetTile(x, y, tileType int) {
	if y < 0 || y >= len(m.Tiles) || x < 0 || x >= len(m.Tiles[y]) {
		return // Ignore out-of-bounds writes
	}
	m.Tiles[y][x] = tileType
}

// AddNode adds a new node to the map
func (m *MapData) AddNode(id int, x, y float32) {
	m.Nodes = append(m.Nodes, Node{
		ID:  id,
		Pos: Pos{X: x, Y: y},
	})
}

// RemoveNode removes a node by ID and all associated paths
func (m *MapData) RemoveNode(id int) {
	// Remove the node
	for i, node := range m.Nodes {
		if node.ID == id {
			m.Nodes = append(m.Nodes[:i], m.Nodes[i+1:]...)
			break
		}
	}

	// Remove all paths that reference this node
	filteredPaths := []Path{}
	for _, path := range m.Paths {
		if path.NodeAID != id && path.NodeBID != id {
			filteredPaths = append(filteredPaths, path)
		}
	}
	m.Paths = filteredPaths
}

// AddPath adds a path between two nodes
func (m *MapData) AddPath(nodeAID, nodeBID int, cost float32) {
	m.Paths = append(m.Paths, Path{
		NodeAID: nodeAID,
		NodeBID: nodeBID,
		Cost:    cost,
	})
}

// FindNodeByID finds a node by its ID
func (m *MapData) FindNodeByID(id int) *Node {
	for i, node := range m.Nodes {
		if node.ID == id {
			return &m.Nodes[i]
		}
	}
	return nil
}

// GetNextNodeID returns the next available node ID
func (m *MapData) GetNextNodeID() int {
	maxID := 0
	for _, node := range m.Nodes {
		if node.ID > maxID {
			maxID = node.ID
		}
	}
	return maxID + 1
}

// AddSprite adds a sprite to the map
func (m *MapData) AddSprite(spriteType int, x, y float32) {
	m.Sprites = append(m.Sprites, Sprite{
		Type: spriteType,
		Pos:  Pos{X: x, Y: y},
	})
}

// RemoveSprite removes a sprite at the specified position (within tolerance)
func (m *MapData) RemoveSprite(x, y, tolerance float32) bool {
	for i, sprite := range m.Sprites {
		dx := sprite.Pos.X - x
		dy := sprite.Pos.Y - y
		distance := dx*dx + dy*dy // squared distance
		if distance <= tolerance*tolerance {
			m.Sprites = append(m.Sprites[:i], m.Sprites[i+1:]...)
			return true
		}
	}
	return false
}
