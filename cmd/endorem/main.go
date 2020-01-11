package main

import (
	"github.com/SolarLune/resolv/resolv"
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

		for _, en := range *enemies {
			if en.GetData() == nil {
				continue
			}

			// Check the type of the enemy space data.
			// If it's a string, then it's a Hitbox.
			// If it's a reference to itself then it's a Hurtbox.
			switch t := en.GetData().(type) {
			case *enemy.Basic: // Hurtbox
				enX, enY := t.Collision.Center()
				pX, pY := tPlayer.Collision.Center()

				// Calculate the distance from the enemy to the player.
				dist := resolv.Distance(enX, enY, pX, pY)

				t.PlayerSeen = dist < common.GlobalConfig.Enemy.VisionDistance
				t.ShouldAttack = dist < t.AttackDistance

				// If the hurtbox is colliding a player hitbox then take damage.
				if t.FilterByTags(enemy.TagHurtbox).IsColliding(tPlayer.Hitbox) {
					t.TakeDamage()
					// If the player is colliding with the enemy then they should take damage.
				} else if tPlayer.FilterByTags(player.TagHurtbox).IsColliding(t.FilterOutByTags(enemy.TagAttackZone)) {
					tPlayer.TakeDamage()
				}
			}
		}

		r.BeginDrawing()
		r.BeginMode2D(r.Camera2D(*cam)) // Begin drawing with camera.
		r.ClearBackground(r.Black)

		tPlane.Draw()
		tPlayer.Draw()
		basicEnemy.Draw()

		r.EndMode2D()
		r.EndDrawing()
	}
}
