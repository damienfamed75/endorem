package enemy

import (
	"fmt"
	"strconv"

	r "github.com/lachee/raylib-goplus/raylib"
)

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
	// Draw facing direction.
	r.DrawText(
		b.Facing.String(),
		int(b.Collision.X), int(b.Collision.Y+b.Collision.H+10), 10,
		r.White,
	)
	r.DrawText(
		fmt.Sprintf("PSen: %v Atk: %v", b.PlayerSeen, b.ShouldAttack),
		int(b.Collision.X), int(b.Collision.Y+b.Collision.H+20), 10,
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
