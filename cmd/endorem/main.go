package main

import (
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/scene"

	r "github.com/lachee/raylib-goplus/raylib"
)

func main() {
	// Load in the global configuration for all future items to reference.
	common.LoadConfig()

	// Create a new empty game.
	g := NewGame()

	// Register all the available scenes in the game.
	// Should be something like...
	// testing scene
	// main menu
	// Game over screen?
	// level 1
	// level 2
	// level 3
	g.RegisterScenes(&scene.TestingScene{})
	g.RegisterScenes(&scene.MenuScene{})

	// Initialize raylib window.
	r.InitWindow(
		common.GlobalConfig.ScreenWidth(),
		common.GlobalConfig.ScreenHeight(),
		"Endorem",
	)
	defer r.CloseWindow()

	// Window settings
	r.SetTargetFPS(60)

	// Choose default scene.
	g.Start(common.GlobalConfig.Game.DefaultScene)

	// Game loop
	for !r.WindowShouldClose() {
		// Update the game with given deltatime
		g.Update(r.GetFrameTime())

		r.BeginDrawing()
		// Draw the game.
		g.Draw()

		r.EndDrawing()
	}
}
