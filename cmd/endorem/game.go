package main

import (
	"fmt"

	"github.com/damienfamed75/endorem/pkg/scene"
)

type Game struct {
	Scenes       []scene.Scene
	CurrentScene int
}

func setupGame() *Game {
	return &Game{
		CurrentScene: -1,
	}
}

func NewGame() *Game {
	g := setupGame()

	return g
}

func (g *Game) RegisterScenes(ss ...scene.Scene) {
	g.Scenes = append(g.Scenes, ss...)
}

func (g *Game) FindScene(sceneName string) (scene.Scene, int) {
	for i, s := range g.Scenes {
		if s.String() == sceneName {
			return s, i
		}
	}

	return nil, -1
}

func (g *Game) SwitchScene(index int) {
	g.Scenes[g.CurrentScene].Unload() // Unload the current scene.
	// Debug log to unload scene.

	g.Scenes[index].Preload() // Preload new scene.
	g.CurrentScene = index    // Switch the current scene to given.
}

func (g *Game) Start(defaultScene string) {
	if s, i := g.FindScene(defaultScene); s != nil {
		s.Preload()
		g.CurrentScene = i
	}
}

// Update the current scene.
func (g *Game) Update(dt float32) {
	g.Scenes[g.CurrentScene].Update(dt)
}

// Draw the current scene.
func (g *Game) Draw() {
	g.Scenes[g.CurrentScene].Draw()
}

// GameOver handles a game over.
func (g *Game) GameOver() {
	fmt.Println("GAME OVER")
}
