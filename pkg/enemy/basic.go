package enemy

import (
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Basic is a testing enemy that is very basic in attacks and features.
type Basic struct {
	Health     uint8
	Sprite     r.Texture2D
	Collision  *resolv.Rectangle
	AttackZone *resolv.Rectangle
	Hurtbox    *resolv.Rectangle

	begin       time.Time
	isAttacking bool

	*resolv.Space
}

func NewBasic(x, y int) *Basic {
	var (
		attackZoneWidth  int32 = 2
		attackZoneHeight int32 = 2
	)

	b := &Basic{
		Sprite: r.LoadTexture("assets/basicenemy.png"),
		Space:  resolv.NewSpace(),
		Health: 2,
	}

	// Set the hit and hurt boxes.
	b.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		b.Sprite.Width, b.Sprite.Height,
	)
	b.AttackZone = resolv.NewRectangle(
		int32(x), int32(y),
		b.Sprite.Height*attackZoneWidth, b.Sprite.Height*attackZoneHeight,
	)
	b.Hurtbox = resolv.NewRectangle(
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
	b.Hurtbox.AddTags(common.TagEnemy)

	b.begin = time.Now()

	return b
}

func (b *Basic) Update() {
	// Debugging:
	// Timer for attacks every half second.
	if time.Since(b.begin) >= time.Millisecond*500 {
		// Reset timer.
		b.begin = time.Now()

		// Flip attack value.
		b.isAttacking = !b.isAttacking
		if b.isAttacking {
			// Re-add hurtbox to the enemy space and set position to enemy.
			b.Hurtbox.SetXY(b.Collision.X, b.Collision.Y+b.Collision.H/3.0)
			b.Add(b.Hurtbox)
		} else {
			// Remove hurtbox from enemy space.
			b.Remove(b.Hurtbox)
		}
	}
}

func (b *Basic) Draw() {
	// Draw the enemy texture.
	r.DrawTexture(b.Sprite, int(b.Collision.X), int(b.Collision.Y), r.White)

	b.debugDraw()
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
			int(b.Hurtbox.X), int(b.Hurtbox.Y),
			int(b.Hurtbox.W), int(b.Hurtbox.H),
			r.Green,
		)
	}
}
