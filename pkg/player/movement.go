package player

import (
	"log"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard playable character, including functions that allow
// for movement and action
type Player struct {
	Sprite    r.Texture2D
	Collision *resolv.Rectangle
	Hitbox    *resolv.Rectangle

	SpeedX float32
	SpeedY float32

	onGround    bool
	isAttacking bool

	*resolv.Space
}

// NewPlayer creates a player struct, loading the player sprite texture and generates
// the collision space for the player
func NewPlayer(x, y int) *Player {
	// spr := r.LoadTexture("assets/playerTest.png")
	// playerSpace := resolv.NewSpace()

	// playerSpace.Add(
	// 	resolv.NewRectangle(0, 468, spr.Width, spr.Height),
	// )
	// return &Player{
	// 	Collision: playerSpace,
	// 	Sprite:    spr,
	// }
	p := &Player{
		Sprite: r.LoadTexture("assets/playerTest.png"),
		Space:  resolv.NewSpace(),
	}

	p.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		p.Sprite.Width, p.Sprite.Height,
	)

	p.Hitbox = resolv.NewRectangle(
		0, 0, p.Sprite.Height, p.Sprite.Width,
	)

	//Add to collision boxes to player space.
	p.Add(p.Collision)

	return p
}

// movePlayer handles key binded events involving the movement of the character
func (p *Player) movePlayer(ground *resolv.Space) {

	// Left/Right Movement
	//	p.SpeedY += 0.5

	friction := float32(0.5)
	accel := 0.5 + friction

	maxSpd := float32(3)

	if p.SpeedX > friction {
		p.SpeedX -= friction
	} else if p.SpeedX < -friction {
		p.SpeedX += friction
	} else {
		p.SpeedX = 0
	}

	if r.IsKeyDown(r.KeyD) {
		p.SpeedX += accel
	}
	if r.IsKeyDown(r.KeyA) {
		p.SpeedX -= accel
	}

	if p.SpeedX > maxSpd {
		p.SpeedX = maxSpd
	}
	if p.SpeedX < -maxSpd {
		p.SpeedX = -maxSpd
	}
	x := int32(p.SpeedX)
	//y := int32(p.SpeedY)

	if res := ground.Resolve(p.Collision, x, 0); res.Colliding() {
		x = res.ResolveX
		p.SpeedX = 0
	}

	p.Collision.X += x

	// Jumping
	if r.IsKeyPressed(r.KeyW) && p.onGround {
		p.Collision.Move(0, -20)
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

func (p *Player) checkAttack() {
	// Attacking
	if r.IsMouseButtonPressed(r.MouseLeftButton) {
		// Re-add hurtbox to the enemy space and set position to enemy.
		p.Hitbox.SetXY(p.Collision.X, p.Collision.Y+p.Collision.H/3.0)
		p.Add(p.Hitbox)
		p.isAttacking = true
	} else {
		// Remove hurtbox from enemy space.
		p.Remove(p.Hitbox)
		p.isAttacking = false
	}
}

// checkInAir determines if the player is colliding with the ground, and if not they will
// fall towards the ground
func (p *Player) checkInAir(ground *resolv.Space) {
	if p.Collision.IsColliding(ground) {
		p.onGround = true
	} else {
		p.onGround = false
		p.Collision.Move(0, 1)
	}
}

func (p *Player) Update(ground *resolv.Space) {
	p.movePlayer(ground)
	p.checkAttack()
	p.checkInAir(ground)
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()
	x, y := p.Collision.GetXY()
	r.DrawTexture(p.Sprite, int(x), int(y), r.White)

	p.debugDraw()
}

func (p *Player) debugDraw() {
	r.DrawRectangleLines(
		int(p.Collision.X), int(p.Collision.Y),
		int(p.Collision.W), int(p.Collision.H),
		r.Red,
	)
	if p.isAttacking {
		r.DrawRectangleLines(
			int(p.Hitbox.X), int(p.Hitbox.Y),
			int(p.Hitbox.W), int(p.Hitbox.H),
			r.Green,
		)
	}
}
