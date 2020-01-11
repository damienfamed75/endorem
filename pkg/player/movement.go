package player

import (
	"log"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard playable character, including functions that allow
// for movement and action
type Player struct {
	Space  *resolv.Space
	Sprite r.Texture2D

	onGround bool
}

// NewPlayer creates a player struct, loading the player sprite texture and generates
// the collision space for the player
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

// MovePlayer handles key binded events involving the movement of the character
func (p *Player) MovePlayer() {
	if r.IsKeyDown(r.KeyD) {
		p.Space.Move(1, 0)
	}
	if r.IsKeyDown(r.KeyA) {
		p.Space.Move(-1, 0)
	}

	// Jumping
	if r.IsKeyPressed(r.KeyW) && p.onGround {
		p.Space.Move(0, -20)
	}

	// Crouching
	// Changes to crouch sprite and hurtboxes
	if r.IsKeyDown(r.KeyS) {
		//TODO
		log.Print("Woah you are crouching")
	} else {
		//TODO
	}
}

//CheckInAir determines if the player is colliding with the ground, and if not they will
// fall towards the ground
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

	p.debugDraw()
}

func (p *Player) debugDraw() {

}
