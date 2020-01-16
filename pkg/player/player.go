package player

import (
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard playable character, including functions that allow
// for movement and action
type Player struct {
	MaskObj     *Mask
	SpriteStand r.Texture2D
	SpriteDuck  r.Texture2D
	Collision   *resolv.Rectangle
	Hitbox      *physics.Rectangle
	Health      int
	MaxHealth   int
	IsDead      bool
	Facing      common.Direction
	*physics.Body

	maxSpeedX  int32
	maxSpeedY  int32
	SpeedX     float32
	SpeedY     float32
	jumpHeight float32
	madeJump   bool

	// Ground *resolv.Space

	isAttacking     bool
	isCrouched      bool
	deathFunc       func()
	healthBefore    time.Time
	attackBefore    time.Time
	attackTimer     time.Duration
	invincibleTimer time.Duration
	state           common.State

	*Inventory
	//*resolv.Space
	// *common.Rigidbody
}

func setupPlayer() *Player {
	return &Player{
		MaskObj:     NewMask(),
		Inventory:   NewInventory(),
		SpriteStand: r.LoadTexture("assets/playerTest.png"),
		SpriteDuck:  r.LoadTexture("assets/playerDuck.png"),
		//Space:           resolv.NewSpace(),
		Facing:          common.Right,
		Health:          3,
		maxSpeedX:       4,
		maxSpeedY:       8,
		jumpHeight:      -6,
		madeJump:        false,
		healthBefore:    time.Now(),
		attackBefore:    time.Now(),
		attackTimer:     time.Duration(common.GlobalConfig.Player.AttackTimer),
		invincibleTimer: time.Duration(common.GlobalConfig.Player.InvincibleTimer),
		state:           common.StateIdle,
	}
}

// NewPlayer creates a player struct, loading the player sprite texture and generates
// the collision space for the player
func NewPlayer(x, y int, deathFunc func(), ground *physics.Space) *Player {
	p := setupPlayer()
	p.MaskObj.setMovePattern("test")
	// Create Mask to follow player
	// Set the death function that'll be called when the player dies.
	p.deathFunc = deathFunc

	p.Body = physics.NewBody(
		float32(x), float32(y),
		float32(p.SpriteStand.Width), float32(p.SpriteStand.Height),
		float32(p.maxSpeedX), float32(p.maxSpeedY),
	)

	// Setup collision and trigger boxes for the player.
	p.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		p.SpriteStand.Width, p.SpriteStand.Height,
	)
	p.Collision.AddTags(common.TagCollision)
	p.Hitbox = physics.NewRectangle(
		0, 0, float32(p.SpriteStand.Height), float32(p.SpriteStand.Width),
	)

	// Set all spaces to have self referencial data.
	p.SetData(p)

	p.Hitbox.AddTags(TagHitbox)

	// Add to collision boxes to player space.
	p.Add(p.Hitbox)

	p.AddGround(*ground...)

	return p
}

// Update player
func (p *Player) Update(dt float32) r.Vector2 {
	p.state = common.StateIdle

	p.movePlayer()

	maskTar := r.Vector2{
		X: float32(p.Position().X),
		Y: float32(p.Position().Y),
	}
	p.MaskObj.checkDirection(maskTar, p.Facing)
	p.MaskObj.Update()

	p.checkAttack()

	// p.Rigidbody.Update(dt)
	p.Body.Update(dt)
	return r.NewVector2(float32(p.Body.Position().X), float32(p.Body.Position().Y))
	// return r.NewVector2(float32(p.Collision.X), float32(p.Collision.Y))
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	p.MaskObj.Draw()

	if p.isCrouched {
		r.DrawTexture(p.SpriteDuck, int(p.Position().X), int(p.Position().Y), r.White)
	} else {
		r.DrawTexture(
			p.SpriteStand,
			int(p.Body.Position().X),
			int(p.Body.Position().Y),
			r.White,
		)
		// r.DrawTexture(p.SpriteStand, int(p.GetX()), int(p.GetY()), r.White)
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

// movePlayer handles key binded events involving the movement of the character
func (p *Player) movePlayer() {

	// Left/Right Movement

	friction := float32(0.5)
	// accel := 0.5 + friction

	// Slows down player movement after key is released
	if p.Velocity.X > friction {
		p.Body.Velocity.X -= friction
		// p.Rigidbody.Velocity.X -= friction
	} else if p.Velocity.X < -friction {
		p.Body.Velocity.X += friction
		// p.Rigidbody.Velocity.X += friction
	} else {
		p.Body.Velocity.X = 0
		// p.Rigidbody.Velocity.X = 0
	}

	// Controller Events
	// TODO For simplicity in testing the velocity is the maxSpeed of X
	if r.IsKeyDown(r.KeyD) {
		p.Body.Velocity.X += float32(p.maxSpeedX)
		// p.Rigidbody.Velocity.X += float32(p.maxSpeedX)
		p.Facing = common.Right
		p.state = common.StateRight
	}
	if r.IsKeyDown(r.KeyA) {
		p.Body.Velocity.X -= float32(p.maxSpeedX)
		// p.Rigidbody.Velocity.X -= float32(p.maxSpeedX)
		p.Facing = common.Left
		p.state = common.StateLeft
	}

	// JUMPING
	// if the player isn't crouched allow for jump events
	if !p.isCrouched {
		p.playerJump()
	}

	// Crouching
	// Changes to crouch sprite and hurtboxes

	// if r.IsKeyDown(r.KeyS) {
	// 	p.Collision.H = p.SpriteDuck.Height
	// 	p.Collision.W = p.SpriteDuck.Width
	// 	p.isCrouched = true
	// 	p.state = common.StateCrouch
	// } else {
	// 	p.Collision.H = p.SpriteStand.Height
	// 	p.Collision.W = p.SpriteStand.Width
	// 	p.isCrouched = false
	// }

	// return r.NewVector2(float32(x), float32(y))
}

func (p *Player) playerJump() {
	if r.IsKeyPressed(r.KeyW) && p.Body.OnGround() {
		// fmt.Println("JUMP!")
		// if r.IsKeyPressed(r.KeyW) && p.OnGround() {
		p.Body.Velocity.Y = p.jumpHeight
		// p.Rigidbody.Velocity.Y = p.jumpHeight
		p.madeJump = true
	} else if r.IsKeyPressed(r.KeyW) && p.madeJump {
		p.Body.Velocity.Y = p.jumpHeight
		// p.Rigidbody.Velocity.Y = p.jumpHeight
		p.madeJump = false
	}
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
		p.Body.Remove(p.Hitbox)
		p.isAttacking = false
	}
}

func (p *Player) attack() {
	// Based on the direction the player is facing, set the position of the
	// hitbox in front of the player.
	if p.Facing == common.Left {
		p.Hitbox.SetPosition(r.NewVector2(
			p.Body.Position().X-(p.Hitbox.Width/2),
			p.Body.Position().Y+p.Body.Collider().Height/3.0,
		))
	} else {
		p.Hitbox.SetPosition(r.NewVector2(
			p.Body.Position().X,
			p.Body.Position().Y+p.Body.Collider().Height/3.0,
		))
	}

	p.attackBefore = time.Now() // Reset timerS
	// p.Add(p.Hitbox)
	p.Body.Add(p.Hitbox)
	p.isAttacking = true
}
