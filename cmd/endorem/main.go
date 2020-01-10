package main

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

func main() {
	r.InitWindow(800, 600, "Endorem")

	for !r.WindowShouldClose() {
		r.BeginDrawing()
		r.ClearBackground(r.RayWhite)

		r.DrawText("Endorem hello", 20, 20, 40, r.GopherBlue)

		r.EndDrawing()
	}

	defer r.CloseWindow()
}
