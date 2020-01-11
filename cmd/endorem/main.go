package main

import (
	"github.com/damienfamed75/endorem/pkg/testing"
	r "github.com/lachee/raylib-goplus/raylib"
)

func main() {
	r.InitWindow(800, 600, "Endorem")
	defer r.CloseWindow()

	tPlane := testing.NewPlane()

	for !r.WindowShouldClose() {
		r.BeginDrawing()
		r.ClearBackground(r.Black)

		r.DrawText("Endorem hello", 20, 20, 40, r.GopherBlue)

		tPlane.Draw()
		
		r.EndDrawing()
	}
}
