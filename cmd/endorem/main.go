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

	// Add everything to the world space.
	g.world.Add(tPlane, tPlayer, basicEnemy)

	for !r.WindowShouldClose() {

		tPlayer.Update(tPlane.Space)
		basicEnemy.Update()

		enemies := g.world.FilterByTags(common.TagEnemy)

		// Player has been touched by an enemy.
		if enemies.IsColliding(tPlayer) {
			tPlayer.TakeDamage()
		}

		r.BeginDrawing()
		r.ClearBackground(r.Black)

		r.DrawText("Endorem hello", 20, 20, 40, r.GopherBlue)

		tPlane.Draw()
		tPlayer.Draw()
		basicEnemy.Draw()

		r.EndDrawing()
	}
}
