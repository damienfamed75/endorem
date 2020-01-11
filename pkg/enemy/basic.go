package enemy

import (
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Basic is a testing enemy that is very basic in attacks and features.
type Basic struct {
	Health     int
	SpeedX     float32
	SpeedY     float32
	Sprite     r.Texture2D
	Collision  *resolv.Rectangle
	AttackZone *resolv.Rectangle
	Hitbox     *resolv.Rectangle

	begin       time.Time
	isAttacking bool
	// Multiplier of the enemy's horizontal movement.
	speedMultiplier float32
	// How often the enemy can attack (in milliseconds)
	attackTimer time.Duration

	*resolv.Space
}

func setupBasic() *Basic {
	return &Basic{
		Sprite:          r.LoadTexture("assets/basicenemy.png"),
		Space:           resolv.NewSpace(),
		Health:          2 + common.GlobalConfig.Enemy.AddedHealth,
		speedMultiplier: common.GlobalConfig.Enemy.MoveSpeedMultiplier,
		attackTimer:     time.Duration(common.GlobalConfig.Enemy.AttackTimer),
	}
}

// NewBasic returns a configured basic enemy at the given coordinates.
func NewBasic(x, y int) *Basic {
	b := setupBasic()

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
	b.Add(b.Collision, b.AttackZone)

	// Tag this enemy as an enemy.
	b.Collision.AddTags(common.TagEnemy)
	b.Hitbox.AddTags(common.TagEnemy)

	// DEBUG - timer for enemy's attack atm. Will be removed.
	b.begin = time.Now()

	return b
}

// Update is non drawing related functionality with the enemy.
func (b *Basic) Update() {
	// Debugging:
	// Timer for attacks every half second.
	if time.Since(b.begin) >= time.Millisecond*b.attackTimer {
		// Reset timer.
		b.begin = time.Now()

		// Flip attack value.
		b.isAttacking = !b.isAttacking
		if b.isAttacking {
			// Re-add hurtbox to the enemy space and set position to enemy.
			b.Hitbox.SetXY(b.Collision.X, b.Collision.Y+b.Collision.H/3.0)
			b.Add(b.Hitbox)
		} else {
			// Remove hurtbox from enemy space.
			b.Remove(b.Hitbox)
		}
	}
}

// Draw is used for raylib exclusive drawing function calls.
func (b *Basic) Draw() {
	// Draw the enemy texture.
	r.DrawTexture(b.Sprite, int(b.Collision.X), int(b.Collision.Y), r.White)

	b.debugDraw() //DEBUG
}

func (b *Basic) debugDraw() {
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

	// If the enemy is attacking then draw the debug collision box.
	if b.isAttacking {
		r.DrawRectangleLines(
			int(b.Hitbox.X), int(b.Hitbox.Y),
			int(b.Hitbox.W), int(b.Hitbox.H),
			r.Green,
		)
	}
}
