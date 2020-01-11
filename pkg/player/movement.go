package player

import (
	"log"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard character collision and image
type Player struct {
	Space  *resolv.Space
	Sprite r.Texture2D

	onGround bool
}

func NewPlayer() *Player {
	spr := r.LoadTexture("assets/playerTest.png")
	playerSpace := resolv.NewSpace()

	playerSpace.Add(
		resolv.NewRectangle(0, 468, spr.Width, spr.Height),
	)
	return &Player{
		Space:  playerSpace,
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

	if r.IsKeyPressed(r.KeySpace) && p.onGround {
		p.Space.Move(0, -20)
		log.Print(p.onGround)
	}
}

func (p *Player) CheckInAir(ground *resolv.Space) {
	if p.Space.IsColliding(ground) {
		p.onGround = true
	} else {
		p.onGround = false
		p.Space.Move(0, 1)
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()
	x, y := p.Space.GetXY()
	r.DrawTexture(p.Sprite, int(x), int(y), r.White)
}
