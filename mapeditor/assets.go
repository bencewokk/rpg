package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// Define light theme colors
var (
	lightGray  = color.RGBA{240, 240, 240, 255}
	mediumGray = color.RGBA{200, 200, 200, 255}
	darkGray   = color.RGBA{128, 128, 128, 255}
	white      = color.RGBA{255, 255, 255, 255}
	lightBlue  = color.RGBA{173, 216, 230, 255}
	lightGreen = color.RGBA{144, 238, 144, 255}
	lightBrown = color.RGBA{222, 184, 135, 255}
	voidColor  = color.RGBA{64, 64, 64, 255}
)

type AssetManager struct {
	tileTextures   [4]*ebiten.Image // For tile types 0, 1, 2, 3
	fallbackColors [4]color.RGBA
	assetsLoaded   bool
	npcSpritePaths []string // scanned NPC sprite paths (relative like import/Characters/...)
}

func NewAssetManager() AssetManager {
	return AssetManager{
		fallbackColors: [4]color.RGBA{
			voidColor,  // 0 = void
			darkGray,   // 1 = mountains
			lightGreen, // 2 = plains/grass
			lightBrown, // 3 = dry
		},
		assetsLoaded: false,
	}
}

func (a *AssetManager) LoadAssets(importPath string) error {
	fmt.Printf("Attempting to load assets from: %s\n", importPath)

	// Try to load actual textures
	grassPath := filepath.Join(importPath, "tiles", "Grass_S1.png")
	dryPath := filepath.Join(importPath, "tiles", "Dryland_S1.png")

	// Load grass texture for plains (type 2)
	if img, err := a.loadPNG(grassPath); err == nil {
		a.tileTextures[2] = a.resizeImage(img, tileSize, tileSize)
		fmt.Println("Loaded grass texture")
	}

	// Load dry texture for dry land (type 3)
	if img, err := a.loadPNG(dryPath); err == nil {
		a.tileTextures[3] = a.resizeImage(img, tileSize, tileSize)
		fmt.Println("Loaded dry texture")
	}

	// Create fallback textures for missing assets
	a.createFallbackTextures()

	// Scan NPC sprites
	a.scanNPCSprites(importPath)

	a.assetsLoaded = true
	return nil
}

func (a *AssetManager) loadPNG(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imgData, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(imgData), nil
}

func (a *AssetManager) resizeImage(img *ebiten.Image, width, height int) *ebiten.Image {
	resized := ebiten.NewImage(width, height)

	op := &ebiten.DrawImageOptions{}
	originalW, originalH := img.Size()
	scaleX := float64(width) / float64(originalW)
	scaleY := float64(height) / float64(originalH)
	op.GeoM.Scale(scaleX, scaleY)

	resized.DrawImage(img, op)
	return resized
}

func (a *AssetManager) createFallbackTextures() {
	for i := 0; i < 4; i++ {
		if a.tileTextures[i] == nil {
			// Create a simple colored square
			img := ebiten.NewImage(tileSize, tileSize)
			img.Fill(a.fallbackColors[i])
			a.tileTextures[i] = img
		}
	}
}

func (a *AssetManager) GetTileTexture(tileType int) *ebiten.Image {
	if tileType < 0 || tileType >= len(a.tileTextures) {
		return a.tileTextures[0] // Return void texture for invalid types
	}
	return a.tileTextures[tileType]
}

// scanNPCSprites walks the Characters directory under importPath and records all PNG files.
func (a *AssetManager) scanNPCSprites(importPath string) {
	base := filepath.Join(importPath, "Characters")
	var list []string
	filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		nameLower := strings.ToLower(d.Name())
		if strings.HasSuffix(nameLower, ".png") {
			rel, err := filepath.Rel(base, path)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			// Build relative path expected by game (strip any leading ../)
			spritePath := "import/Characters/" + rel
			list = append(list, spritePath)
		}
		return nil
	})
	sort.Strings(list)
	// Deduplicate
	dedup := make([]string, 0, len(list))
	last := ""
	for _, p := range list {
		if p != last {
			dedup = append(dedup, p)
			last = p
		}
	}
	a.npcSpritePaths = dedup
	fmt.Printf("Scanned %d NPC sprite(s)\n", len(a.npcSpritePaths))
}

// GetNPCSpritePaths returns the scanned list of NPC sprite paths.
func (a *AssetManager) GetNPCSpritePaths() []string { return a.npcSpritePaths }
