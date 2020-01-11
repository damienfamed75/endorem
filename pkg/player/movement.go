package player

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard character collision and image
type Player struct {
	Space  *resolv.Space
	Sprite r.Texture2D

	inAir bool
}

func NewPlayer() *Player {
	spr := r.LoadTexture("assets/playerTest.png")
	playerCol := resolv.NewSpace()

	playerCol.Add(
		resolv.NewRectangle(0, 468, spr.Width, spr.Height),
	)
	return &Player{
		Space:  playerCol,
		Sprite: spr,
	}
}

func (p *Player) MovePlayer() {
	if r.IsKeyDown(r.KeyD) {
		p.Space.Move(1, 0)
	}
	if r.IsKeyDown(r.KeyA) {
		p.Space.Move(-1, 0)
	}
	//Improve Jump
	if r.IsKeyPressed(r.KeySpace) {
		p.Space.Move(0, -5)
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()
	x, y := p.Space.GetXY()
	r.DrawTexture(p.Sprite, int(x), int(y), r.White)
}
