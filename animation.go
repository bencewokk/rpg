package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// Animation represents a sequence of frames with timing info.
type Animation struct {
	Name          string
	Frames        []*ebiten.Image
	FrameDuration float64 // seconds per frame
	Loop          bool
}

// AnimationDef is the JSON-loaded definition (paths instead of images).
type AnimationDef struct {
	Name          string   `json:"name"`
	Frames        []string `json:"frames"`
	FrameDuration float64  `json:"frame_duration"`
	Loop          bool     `json:"loop"`
}

// Manifest JSON structure: entity -> []AnimationDef
type AnimationManifest map[string][]AnimationDef

// AnimationManager holds all animations keyed by entity then name.
type AnimationManager struct {
	entities map[string]map[string]*Animation
}

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{entities: make(map[string]map[string]*Animation)}
}

// LoadManifest loads animations from a JSON file; on failure returns error (caller may fallback).
func (am *AnimationManager) LoadManifest(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	var manifest AnimationManifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		return err
	}

	// Load each animation's frames.
	for entity, defs := range manifest {
		if _, ok := am.entities[entity]; !ok {
			am.entities[entity] = make(map[string]*Animation)
		}
		for _, def := range defs {
			var frames []*ebiten.Image
			for _, framePath := range def.Frames {
				img := loadPNG(filepath.ToSlash(framePath))
				frames = append(frames, img)
			}
			if len(frames) == 0 {
				// Skip empty animations
				continue
			}
			// Default frame duration if omitted
			fd := def.FrameDuration
			if fd <= 0 {
				fd = 0.13
			}
			am.entities[entity][def.Name] = &Animation{
				Name:          def.Name,
				Frames:        frames,
				FrameDuration: fd,
				Loop:          def.Loop,
			}
		}
	}
	fmt.Printf("Loaded animations for %d entity types\n", len(am.entities))
	return nil
}

// Get returns an animation for an entity & name; nil if missing.
func (am *AnimationManager) Get(entity, name string) *Animation {
	if e, ok := am.entities[entity]; ok {
		return e[name]
	}
	return nil
}

// AnimationPlayer keeps per-instance playback state.
type AnimationPlayer struct {
	Anim       *Animation
	FrameIndex int
	Accum      float64
	Finished   bool
}

func (p *AnimationPlayer) SetAnimation(a *Animation, reset bool) {
	if a == nil {
		return
	}
	if p.Anim != a || reset {
		p.Anim = a
		p.FrameIndex = 0
		p.Accum = 0
		p.Finished = false
	}
}

func (p *AnimationPlayer) Update(dt float64) *ebiten.Image {
	if p.Anim == nil || len(p.Anim.Frames) == 0 {
		return nil
	}
	if p.Finished {
		return p.Anim.Frames[p.FrameIndex]
	}
	p.Accum += dt
	for p.Accum >= p.Anim.FrameDuration {
		p.Accum -= p.Anim.FrameDuration
		p.FrameIndex++
		if p.FrameIndex >= len(p.Anim.Frames) {
			if p.Anim.Loop {
				p.FrameIndex = 0
			} else {
				p.FrameIndex = len(p.Anim.Frames) - 1
				p.Finished = true
				break
			}
		}
	}
	return p.Anim.Frames[p.FrameIndex]
}

// Global animation manager instance (initialized in gameinit)
var animationManager *AnimationManager
