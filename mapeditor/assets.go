package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// Define light theme colors
var (
	lightGray     = color.RGBA{240, 240, 240, 255}
	mediumGray    = color.RGBA{200, 200, 200, 255}
	darkGray      = color.RGBA{128, 128, 128, 255}
	white         = color.RGBA{255, 255, 255, 255}
	lightBlue     = color.RGBA{173, 216, 230, 255}
	lightGreen    = color.RGBA{144, 238, 144, 255}
	lightBrown    = color.RGBA{222, 184, 135, 255}
	voidColor     = color.RGBA{64, 64, 64, 255}
)

type AssetManager struct {
	tileTextures [4]*ebiten.Image // For tile types 0, 1, 2, 3
	fallbackColors [4]color.RGBA
	assetsLoaded bool
}

func NewAssetManager() AssetManager {
	return AssetManager{
		fallbackColors: [4]color.RGBA{
			voidColor,    // 0 = void
			darkGray,     // 1 = mountains  
			lightGreen,   // 2 = plains/grass
			lightBrown,   // 3 = dry
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
