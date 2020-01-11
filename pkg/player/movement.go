package player

import (
	"strconv"
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard playable character, including functions that allow
// for movement and action
type Player struct {
	Sprite    r.Texture2D
	Collision *resolv.Rectangle
	Hitbox    *resolv.Rectangle
	Health    int
	IsDead    bool

	SpeedX float32
	SpeedY float32

	onGround        bool
	isAttacking     bool
	state           string
	deathFunc       func()
	healthBefore    time.Time
	invincibleTimer time.Duration

	*resolv.Space
}

// NewPlayer creates a player struct, loading the player sprite texture and generates
// the collision space for the player
func NewPlayer(x, y int, deathFunc func()) *Player {
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
		Sprite:          r.LoadTexture("assets/playerTest.png"),
		Space:           resolv.NewSpace(),
		Health:          3,
		healthBefore:    time.Now(),
		invincibleTimer: time.Duration(common.GlobalConfig.Player.InvincibleTimer),
		deathFunc:       deathFunc,
		state:           common.StateIdle,
	}

	p.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		p.Sprite.Width, p.Sprite.Height,
	)

	p.Hitbox = resolv.NewRectangle(
		0, 0, p.Sprite.Height, p.Sprite.Width,
	)
	p.Hitbox.SetData(HitboxData)
	p.Collision.SetData(HurtboxData)

	//Add to collision boxes to player space.
	p.Add(p.Collision)

	return p
}

// movePlayer handles key binded events involving the movement of the character
func (p *Player) movePlayer(ground *resolv.Space) r.Vector2 {

	// Left/Right Movement
	p.SpeedY += 0.5

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
		p.state = common.StateRight
	}
	if r.IsKeyDown(r.KeyA) {
		p.SpeedX -= accel
		p.state = common.StateLeft
	}

	if p.SpeedX > maxSpd {
		p.SpeedX = maxSpd
	}
	if p.SpeedX < -maxSpd {
		p.SpeedX = -maxSpd
	}

	// Jumping
	down := ground.Resolve(p.Collision, 0, 4)
	onGround := down.Colliding()

	if r.IsKeyPressed(r.KeyW) && onGround {
		p.SpeedY = -8
	}

	x := int32(p.SpeedX)
	y := int32(p.SpeedY)

	// if res := ground.Resolve(p.Collision, x, 0); res.Colliding() {
	// 	x = res.ResolveX
	// 	p.SpeedX = 0
	// }

	p.Collision.X += x

	res := ground.Resolve(p.Collision, 0, y+4)

	if y < 0 || (res.Teleporting && res.ResolveY < -p.Collision.H/2) {
		res = resolv.Collision{}
	}
	if !res.Colliding() {
		res = ground.Resolve(p.Collision, 0, y)
	}

	if res.Colliding() {
		y = res.ResolveY
		p.SpeedY = 0
	}
	p.Collision.Y += y

	// Crouching
	// Changes to crouch sprite and hurtboxes
	if r.IsKeyDown(r.KeyS) {
		//TODO
		p.state = common.StateCrouch
	} else {
		//TODO
	}

	return r.NewVector2(float32(x), float32(y))
}

func (p *Player) checkAttack() {
	// Attacking
	if r.IsMouseButtonPressed(r.MouseLeftButton) {
		// Re-add hurtbox to the enemy space and set position to enemy.
		p.Hitbox.SetXY(p.Collision.X, p.Collision.Y+p.Collision.H/3.0)
		p.Add(p.Hitbox)
		p.isAttacking = true
		p.state = common.StateAttack
	} else {
		// Remove hurtbox from enemy space.
		p.Remove(p.Hitbox)
		p.isAttacking = false
	}
}

func (p *Player) Update(ground *resolv.Space) (r.Vector2, r.Vector2) {
	p.state = common.StateIdle

	diff := p.movePlayer(ground)
	p.checkAttack()
	//p.checkInAir(ground)

	return diff, r.NewVector2(float32(p.Collision.X), float32(p.Collision.Y))
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()
	x, y := p.Collision.GetXY()
	r.DrawTexture(p.Sprite, int(x), int(y), r.White)

	p.debugDraw()
}

func (p *Player) TakeDamage() {
	// The player has their invincibility frames, and if they have run out of
	// time of that then they can take more damage.
	if time.Since(p.healthBefore) >= time.Millisecond*p.invincibleTimer {
		p.healthBefore = time.Now()

		p.Health--
		if p.Health <= 0 {
			p.deathFunc()
			p.state = common.StateDead // :c
		}
	}
}

func (p *Player) debugDraw() {
	if p.SpeedY < 0 {
		p.state = common.StateJumping
	} else if p.SpeedY > 0 {
		p.state = common.StateFalling
	}

	// Draw health.
	r.DrawText(
		"HP: "+strconv.Itoa(p.Health),
		int(p.Collision.X), int(p.Collision.Y-(p.Collision.W/2)), 10,
		r.White,
	)
	// Draw state.
	r.DrawText(
		p.state,
		int(p.Collision.X), int(p.Collision.Y+p.Collision.H), 10,
		r.White,
	)

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
