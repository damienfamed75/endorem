package player

import (
	"fmt"
	"log"
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/aseprite"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Player is the standard playable character, including functions that allow
// for movement and action
type Player struct {
	Ase    *aseprite.File
	Sprite r.Texture2D
	Color  r.Color

	MaskObj *Mask

	Collision *resolv.Rectangle
	Hitbox    *physics.Rectangle
	Health    int
	MaxHealth int
	IsDead    bool
	Facing    common.Direction
	*physics.Body

	damagePushback float32
	maxSpeedX      int32
	maxSpeedY      int32
	SpeedX         float32
	SpeedY         float32
	jumpHeight     float32
	madeJump       bool
	takingDamage   bool

	isAttacking     bool
	isCrouched      bool
	deathFunc       func()
	healthBefore    time.Time
	attackBefore    time.Time
	attackTimer     time.Duration
	invincibleTimer time.Duration
	state           common.State

	soundHurt       *r.Sound
	soundJump       *r.Sound
	soundDoubleJump *r.Sound

	*Inventory
}

func setupPlayer() *Player {
	return &Player{
		Sprite:          r.LoadTexture("assets/player.png"),
		MaskObj:         NewMask(),
		Inventory:       NewInventory(),
		Facing:          common.Right,
		Health:          3,
		MaxHealth:       3,
		maxSpeedX:       4,
		maxSpeedY:       8,
		damagePushback:  20,
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
	var err error

	p.Ase, err = aseprite.Open("assets/player.json")
	if err != nil {
		log.Fatal(err)
	}

	p.soundHurt = r.LoadSound("assets/sounds/Take_Damage_3.wav")
	p.soundJump = r.LoadSound("assets/sounds/Dry_Sword_Swing.wav")
	p.soundDoubleJump = r.LoadSound("assets/sounds/Dry_Sword_Swing.wav")

	p.soundJump.SetVolume(0.5)
	p.soundDoubleJump.SetVolume(0.5)
	p.soundDoubleJump.SetPitch(2.0)

	// Queues a default animation
	p.Ase.Play("idle")

	// Create Mask to follow player
	// Set the death function that'll be called when the player dies.
	p.deathFunc = deathFunc

	p.Body = physics.NewBody(
		float32(x), float32(y),
		float32(p.Ase.FrameBoundaries().Width), float32(p.Ase.FrameBoundaries().Height),
		float32(p.maxSpeedX), float32(p.maxSpeedY),
	)

	// Setup collision and trigger boxes for the player.
	p.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		int32(p.Ase.FrameBoundaries().Width), int32(p.Ase.FrameBoundaries().Height),
	)
	p.Collision.AddTags(common.TagCollision)

	// Set all spaces to have self referencial data.
	p.SetData(p)

	p.AddGround(*ground...)

	return p
}

// Update player
func (p *Player) Update(dt float32) r.Vector2 {
	p.Ase.Update(dt)
	p.state = common.StateIdle

	if p.MaskObj.Ase.IsPlaying("attack") {
		fmt.Println(p.Ase.CurrentAnimation.Name)
		p.Ase.Play("attack")
	} else {
		if !p.Ase.IsPlaying("damage") && !p.Ase.IsPlaying("run") {
			p.Ase.Play("idle")
		}
	}

	if p.soundHurt.IsPlaying() {
		p.Color = r.Red
	} else {
		p.Color = r.White
	}

	p.movePlayer()

	maskTar := r.Vector2{
		X: float32(p.Position().X),
		Y: float32(p.Position().Y),
	}
	p.MaskObj.checkDirection(maskTar, p.Facing)
	p.MaskObj.Update(dt)

	p.checkAttack()

	// p.Rigidbody.Update(dt)
	p.Body.Update(dt)
	return r.NewVector2(float32(p.Body.Position().X), float32(p.Body.Position().Y))
	// return r.NewVector2(float32(p.Collision.X), float32(p.Collision.Y))
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Player) Draw() {
	p.MaskObj.Draw()

	// srcX, srcY resemble the X and Y pixels where the active sprite is.
	srcX, srcY := p.Ase.FrameBoundaries().X, p.Ase.FrameBoundaries().Y
	w, h := p.Ase.FrameBoundaries().Width, p.Ase.FrameBoundaries().Height

	// src resembles the cropped out area that the sprite is in the spritesheet.
	var src r.Rectangle
	if p.Facing == common.Left {
		src = r.NewRectangle(float32(srcX), float32(srcY), float32(-w), float32(h))
	} else {
		src = r.NewRectangle(float32(srcX), float32(srcY), float32(w), float32(h))
	}

	// dest is the world position that the slime should appear in.
	dest := r.NewRectangle(float32(p.Position().X), float32(p.Position().Y), float32(w), float32(h))

	r.DrawTexturePro(
		p.Sprite, src, dest, r.NewVector2(0, 0), 0, p.Color,
	)

	// p.debugDraw()
}

func (p *Player) TakeDamage(dir float32) {
	// The player has their invincibility frames, and if they have run out of
	// time of that then they can take more damage.
	if time.Since(p.healthBefore) >= time.Millisecond*p.invincibleTimer {
		p.healthBefore = time.Now()

		p.Health--
		r.PlaySound(p.soundHurt)
		p.Velocity.X = p.damagePushback * dir
		p.Body.ResolveForces()

		p.Ase.Play("damage")
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
		p.Ase.Play("run")
	} else if p.Velocity.X < -friction {
		p.Body.Velocity.X += friction
		p.Ase.Play("run")
	} else {
		p.Body.Velocity.X = 0
		if !p.Ase.IsPlaying("damage") && !p.Ase.IsPlaying("attack") {
			p.Ase.Play("idle")
		}
	}

	// Controller Events
	// TODO For simplicity in testing the velocity is the maxSpeed of X
	if r.IsKeyDown(r.KeyRight) {
		p.Body.Velocity.X += float32(p.maxSpeedX)

		p.Facing = common.Right
		p.state = common.StateRight
	}
	if r.IsKeyDown(r.KeyLeft) {
		p.Body.Velocity.X -= float32(p.maxSpeedX)

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
	if r.IsKeyPressed(r.KeyUp) && p.Body.OnGround() {
		// fmt.Println("JUMP!")
		// if r.IsKeyPressed(r.KeyW) && p.OnGround() {
		p.Body.Velocity.Y = p.jumpHeight
		// p.Rigidbody.Velocity.Y = p.jumpHeight
		p.madeJump = true
		r.PlaySound(p.soundJump)
	} else if r.IsKeyPressed(r.KeyUp) && p.madeJump {
		p.Body.Velocity.Y = p.jumpHeight
		// p.Rigidbody.Velocity.Y = p.jumpHeight
		p.madeJump = false
		r.PlaySound(p.soundDoubleJump)
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
	// if p.Facing == common.Left {
	// 	p.Hitbox.SetPosition(r.NewVector2(
	// 		p.Body.Position().X-(p.Hitbox.Width/2),
	// 		p.Body.Position().Y+p.Body.Collider().Height/3.0,
	// 	))
	// } else {
	// 	p.Hitbox.SetPosition(r.NewVector2(
	// 		p.Body.Position().X,
	// 		p.Body.Position().Y+p.Body.Collider().Height/3.0,
	// 	))
	// }

	// p.attackBefore = time.Now() // Reset timerS
	// // p.Add(p.Hitbox)
	// p.Body.Add(p.Hitbox)
	// p.isAttacking = true
}
