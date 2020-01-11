package main

import (
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"

	r "github.com/lachee/raylib-goplus/raylib"
)

func main() {
	// Load in the global configuration for all future items to reference.
	common.LoadConfig()

	g := NewGame()

	r.InitWindow(
		common.GlobalConfig.ScreenWidth(),
		common.GlobalConfig.ScreenHeight(),
		"Endorem",
	)
	defer r.CloseWindow()

	r.SetTargetFPS(60)

	tPlane := testing.NewPlane()
	tPlayer := player.NewPlayer(0, 468, g.GameOver)
	basicEnemy := enemy.NewBasic(100, 468)

	// camera requires player's position for offset.
	cam := NewEndoCamera(tPlayer.Collision) // TODO - move camera to Game

	// Add everything to the world space.
	g.world.Add(tPlane, tPlayer, basicEnemy)

	for !r.WindowShouldClose() {
		// Update the camera's position based on the player's movement.
		cam.Update(tPlayer.Update(tPlane.Space))

		basicEnemy.Update()

		enemies := g.world.FilterByTags(common.TagEnemy)

		// Player has been touched by an enemy.
		if sh := tPlayer.GetCollidingShapes(enemies); sh.GetData() == player.HurtboxData {
			tPlayer.TakeDamage()
		}

		for _, en := range *enemies {
			if en.GetData() == nil {
				continue
			}

			// Check the type of the enemy space data.
			// If it's a string, then it's a Hitbox.
			// If it's a reference to itself then it's a Hurtbox.
			switch t := en.GetData().(type) {
			case *enemy.Basic: // Hurtbox
				// If the hurtbox is colliding a player hitbox then take damage.
				if t.IsColliding(tPlayer.Hitbox) {
					t.TakeDamage()
				}
			}
		}

		r.BeginDrawing()
		r.BeginMode2D(r.Camera2D(*cam)) // Begin drawing with camera.
		r.ClearBackground(r.Black)

		r.DrawText("Endorem hello", 20, 20, 40, r.GopherBlue)

		tPlane.Draw()
		tPlayer.Draw()
		basicEnemy.Draw()

		r.EndMode2D()
		r.EndDrawing()
	}
}
