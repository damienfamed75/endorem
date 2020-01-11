package player

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard character collision and image
type Player struct {
	Collision *resolv.Rectangle
	Sprite    r.Texture2D

	x int
	y int
}

func NewPlayer() *Player {
	spr := r.LoadTexture("/assets/playerTest.png")
	return &Player{
		Collision: resolv.NewRectangle(0, 500, spr.Width, spr.Height),
		Sprite:    spr,
		x:         500,
		y:         500,
	}
}

func (p *Player) MovePlayer() {
	if r.IsKeyDown(r.KeyD) {
		p.x += 2.0
	}
	if r.IsKeyDown(r.KeyA) {
		p.x -= 2.0
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()

	r.DrawTexture(p.Sprite, p.x, p.y, r.White)

	r.DrawRectangleLines(
		p.x,
		p.y,
		100,
		100,
		r.White,
	)
}
