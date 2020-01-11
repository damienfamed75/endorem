package enemy

import (
	"strconv"
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Basic is a testing enemy that is very basic in attacks and features.
type Basic struct {
	Health       int
	SpeedX       float32
	SpeedY       float32
	IsDead       bool
	Origin       r.Vector2
	Destinations [2]r.Vector2 // left and right destinations
	Sprite       r.Texture2D
	Collision    *resolv.Rectangle
	AttackZone   *resolv.Rectangle
	Hitbox       *resolv.Rectangle

	state           common.State
	isAttacking     bool
	speedMultiplier float32   // Multiplier of the enemy's horizontal movement.
	attackBefore    time.Time // How often the enemy can attack (milliseconds)
	attackTimer     time.Duration
	healthBefore    time.Time
	invincibleTimer time.Duration

	*resolv.Space
}

func setupBasic() *Basic {
	return &Basic{
		Sprite:          r.LoadTexture("assets/basicenemy.png"),
		Space:           resolv.NewSpace(),
		Health:          2 + common.GlobalConfig.Enemy.AddedHealth,
		state:           common.StateIdle,
		speedMultiplier: common.GlobalConfig.Enemy.MoveSpeedMultiplier,
		attackTimer:     time.Duration(common.GlobalConfig.Enemy.AttackTimer),
		attackBefore:    time.Now(),
		healthBefore:    time.Now(),
		invincibleTimer: time.Duration(common.GlobalConfig.Enemy.InvincibleTimer),
	}
}

// NewBasic returns a configured basic enemy at the given coordinates.
func NewBasic(x, y int) *Basic {
	b := setupBasic()

	b.Origin = r.NewVector2(float32(x), float32(y))
	b.Destinations = [2]r.Vector2{
		r.NewVector2(
			float32(int32(x)-b.Sprite.Height*2)+float32(b.Sprite.Width/2),
			float32(y)+float32(b.Sprite.Height),
		),
		r.NewVector2(
			float32(int32(x)+b.Sprite.Height*2)+float32(b.Sprite.Width/2),
			float32(y)+float32(b.Sprite.Height),
		),
	}

	var (
		// These variables are only used when instantiating a new enemy so they
		// do not belong in the enemy's structure, and since it's only used in
		// this particular enemy at the moment, they are not variables in the
		// configuration file yet.
		attackZoneWidth  int32 = 2
		attackZoneHeight int32 = 2
	)

	// Set the hit and hurt boxes.
	b.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		b.Sprite.Width, b.Sprite.Height,
	)
	b.AttackZone = resolv.NewRectangle(
		int32(x), int32(y),
		b.Sprite.Height*attackZoneWidth, b.Sprite.Height*attackZoneHeight,
	)
	b.Hitbox = resolv.NewRectangle(
		0, 0, b.Sprite.Height, b.Sprite.Width,
	)

	// Center the attacking zone to the center of the character.
	b.AttackZone.Move(
		int32(float32(-b.Sprite.Width)*(float32(attackZoneWidth)-0.5)),
		-b.Sprite.Height/2,
	)

	// Add the collision boxes to the enemy space.
	b.Add(b.Collision)
	b.SetData(b)

	// Set the hitbox data to be different from the hitbox data.
	b.Hitbox.SetData(HitboxData)

	// Tag this enemy as an enemy.
	b.AddTags(common.TagEnemy)
	b.Hitbox.AddTags(common.TagEnemy)

	return b
}

// Update is non drawing related functionality with the enemy.
func (b *Basic) Update() {
	// Debugging:
	// Timer for attacks every half second.
	if time.Since(b.attackBefore) >= time.Millisecond*b.attackTimer {
		// Reset timer.
		b.attackBefore = time.Now()

		// Flip attack value.
		b.isAttacking = !b.isAttacking
		if b.isAttacking {
			// Re-add hurtbox to the enemy space and set position to enemy.
			b.Hitbox.SetXY(b.Collision.X, b.Collision.Y+b.Collision.H/3.0)
			b.Add(b.Hitbox)
		} else {
			// Remove hurtbox from enemy space.
			b.Remove(b.Hitbox)
			b.state = common.StateIdle
		}
	}
}

// Draw is used for raylib exclusive drawing function calls.
func (b *Basic) Draw() {
	// Draw the enemy texture.
	r.DrawTexture(b.Sprite, int(b.Collision.X), int(b.Collision.Y), r.White)

	b.debugDraw() //DEBUG
}

func (b *Basic) TakeDamage() {
	if time.Since(b.healthBefore) >= time.Millisecond*b.invincibleTimer {
		b.healthBefore = time.Now()

		b.Health--
		if b.Health <= 0 {
			b.IsDead = true
			b.state = common.StateDead
			b.Clear()
		}
	}
}

func (b *Basic) debugDraw() {
	// Draw health.
	r.DrawText(
		"HP: "+strconv.Itoa(b.Health),
		int(b.Collision.X), int(b.Collision.Y-(b.Collision.W/2)), 10,
		r.White,
	)
	// Draw state.
	r.DrawText(
		b.state.String(),
		int(b.Collision.X), int(b.Collision.Y+b.Collision.H), 10,
		r.White,
	)

	// Draw the collision box for debugging reasons.
	r.DrawRectangleLines(
		int(b.Collision.X), int(b.Collision.Y),
		int(b.Collision.W), int(b.Collision.H),
		r.Red,
	)
	r.DrawRectangleLines(
		int(b.AttackZone.X), int(b.AttackZone.Y),
		int(b.AttackZone.W), int(b.AttackZone.H),
		r.Yellow,
	)

	enemyCenterBottom := r.NewVector2(
		float32(b.Collision.X)+float32(b.Collision.W/2),
		float32(b.Collision.Y+b.Collision.H),
	)

	r.DrawLineEx(
		enemyCenterBottom,
		b.Destinations[0],
		3, r.DarkBlue,
	)
	r.DrawLineEx(
		enemyCenterBottom,
		b.Destinations[1],
		3, r.Maroon,
	)

	// If the enemy is attacking then draw the debug collision box.
	if b.isAttacking {
		r.DrawRectangleLines(
			int(b.Hitbox.X), int(b.Hitbox.Y),
			int(b.Hitbox.W), int(b.Hitbox.H),
			r.Green,
		)
	}
}
