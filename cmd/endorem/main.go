package main

import (
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/testing"

	r "github.com/lachee/raylib-goplus/raylib"
)

func main() {
	r.InitWindow(800, 600, "Endorem")
	defer r.CloseWindow()

	tPlane := testing.NewPlane()

	basicEnemy := enemy.NewBasic(100, 450)

	for !r.WindowShouldClose() {
		basicEnemy.Update()

		r.BeginDrawing()
		r.ClearBackground(r.Black)

		r.DrawText("Endorem hello", 20, 20, 40, r.GopherBlue)

		tPlane.Draw()
		basicEnemy.Draw()

		r.EndDrawing()
	}
}
