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
	SpriteStand r.Texture2D
	SpriteDuck  r.Texture2D
	Collision   *resolv.Rectangle
	Hitbox      *resolv.Rectangle
	Health      int
	IsDead      bool
	Facing      common.Direction

	SpeedX float32
	SpeedY float32

	Ground      *resolv.Space
	Transitions *resolv.Space
	onGround    bool

	isAttacking     bool
	isCrouched      bool
	deathFunc       func()
	healthBefore    time.Time
	attackBefore    time.Time
	attackTimer     time.Duration
	invincibleTimer time.Duration
	state           common.State

	*resolv.Space
}

func setupPlayer() *Player {
	return &Player{
		SpriteStand:     r.LoadTexture("assets/playerTest.png"),
		SpriteDuck:      r.LoadTexture("assets/playerDuck.png"),
		Space:           resolv.NewSpace(),
		Facing:          common.Right,
		Health:          3,
		healthBefore:    time.Now(),
		attackBefore:    time.Now(),
		attackTimer:     time.Duration(common.GlobalConfig.Player.AttackTimer),
		invincibleTimer: time.Duration(common.GlobalConfig.Player.InvincibleTimer),
		state:           common.StateIdle,
	}
}

// NewPlayer creates a player struct, loading the player sprite texture and generates
// the collision space for the player
func NewPlayer(x, y int, deathFunc func(), ground *resolv.Space, transitions *resolv.Space) *Player {
	p := setupPlayer()

	// Set the death function that'll be called when the player dies.
	p.deathFunc = deathFunc

	// Setup collision and trigger boxes for the player.
	p.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		p.SpriteStand.Width, p.SpriteStand.Height,
	)

	p.Hitbox = resolv.NewRectangle(
		0, 0, p.SpriteStand.Height, p.SpriteStand.Width,
	)

	// Set all spaces to have self referencial data.
	p.SetData(p)
	// Set the hurtbox to store differing data.
	p.Collision.SetData(HurtboxData)
	p.Collision.AddTags(TagHurtbox)
	p.Hitbox.AddTags(TagHitbox)

	// Add to collision boxes to player space.
	p.Add(p.Collision)

	// Saves the ground
	// for collision detection with player
	p.Ground = ground
	p.Transitions = transitions

	return p
}

// movePlayer handles key binded events involving the movement of the character
func (p *Player) movePlayer() r.Vector2 {

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

	// Controller Events
	if r.IsKeyDown(r.KeyD) {
		p.SpeedX += accel
		p.Facing = common.Right
		p.state = common.StateRight
	}
	if r.IsKeyDown(r.KeyA) {
		p.SpeedX -= accel
		p.Facing = common.Left
		p.state = common.StateLeft
	}

	// Speed Limit
	if p.SpeedX > maxSpd {
		p.SpeedX = maxSpd
	}
	if p.SpeedX < -maxSpd {
		p.SpeedX = -maxSpd
	}

	// JUMPING
	down := p.Ground.Resolve(p.Collision, 0, 4)
	onGround := down.Colliding()

	if r.IsKeyPressed(r.KeyW) && onGround {
		p.SpeedY = -8
	}

	x := int32(p.SpeedX)
	y := int32(p.SpeedY)

	// Check wall collision
	if res := p.Ground.Resolve(p.Collision, x, 0); res.Colliding() {
		x = res.ResolveX
		p.SpeedX = 0
	}

	res := p.Ground.Resolve(p.Collision, 0, y+4)

	if y < 0 || (res.Teleporting && res.ResolveY < -p.Collision.H/2) {
		res = resolv.Collision{}
	}
	if !res.Colliding() {
		res = p.Ground.Resolve(p.Collision, 0, y)
	}

	if res.Colliding() {
		y = res.ResolveY
		p.SpeedY = 0
	}
	p.Collision.X += x
	p.Collision.Y += y

	// Transition Check
	if p.IsColliding(p.Transitions) {

	}

	// Crouching
	// Changes to crouch sprite and hurtboxes
	if r.IsKeyDown(r.KeyS) {
		p.isCrouched = true
		p.state = common.StateCrouch
	} else {
		p.isCrouched = false
	}

	return r.NewVector2(float32(x), float32(y))
}

func (p *Player) checkAttack() {
	// If the last attack was performed too little of time ago then return.
	if time.Since(p.attackBefore) <= time.Millisecond*p.attackTimer {
		return
	}

	// Attacking
	if r.IsMouseButtonPressed(r.MouseLeftButton) {
		// Re-add hurtbox to the enemy space and set position to enemy.
		p.attack()
	} else {
		// Remove hurtbox from enemy space.
		p.Remove(p.Hitbox)
		p.isAttacking = false
	}
}

func (p *Player) attack() {
	// Based on the direction the player is facing, set the position of the
	// hitbox in front of the player.
	if p.Facing == common.Left {
		p.Hitbox.SetXY(p.Collision.X-(p.Hitbox.W/2), p.Collision.Y+p.Collision.H/3.0)
	} else {
		p.Hitbox.SetXY(p.Collision.X, p.Collision.Y+p.Collision.H/3.0)
	}

	p.attackBefore = time.Now() // Reset timer
	p.Add(p.Hitbox)
	p.isAttacking = true
}

func (p *Player) Update() (r.Vector2, r.Vector2) {
	p.state = common.StateIdle

	diff := p.movePlayer()
	p.checkAttack()
	//p.checkInAir(ground)

	return diff, r.NewVector2(float32(p.Collision.X), float32(p.Collision.Y))
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	//p.Collision.SetXY()
	x, y := p.Collision.GetXY()

	if p.isCrouched {
		p.Collision.H = p.SpriteDuck.Height
		p.Collision.W = p.SpriteDuck.Width
		r.DrawTexture(p.SpriteDuck, int(x), int(y), r.White)
	} else {
		p.Collision.H = p.SpriteStand.Height
		p.Collision.W = p.SpriteStand.Width
		r.DrawTexture(p.SpriteStand, int(x), int(y), r.White)
	}
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
		p.state.String(),
		int(p.Collision.X), int(p.Collision.Y+p.Collision.H), 10,
		r.White,
	)
	r.DrawText(
		p.Facing.String(),
		int(p.Collision.X), int(p.Collision.Y+p.Collision.H+10), 10,
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
