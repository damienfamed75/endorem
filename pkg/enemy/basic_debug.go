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
		int(b.GetX()), int(b.GetY()-(b.Sprite.Width/2)), 10,
		r.White,
	)
	// Draw state.
	r.DrawText(
		b.state.String(),
		int(b.GetX()), int(b.GetY()+b.Sprite.Height), 10,
		r.White,
	)
	// Draw facing direction.
	r.DrawText(
		b.Facing.String(),
		int(b.GetX()), int(b.GetY()+b.Sprite.Height+10), 10,
		r.White,
	)
	r.DrawText(
		fmt.Sprintf("PSen: %v Atk: %v", b.PlayerSeen, b.ShouldAttack),
		int(b.GetX()), int(b.GetY()+b.Sprite.Height+20), 10,
		r.White,
	)

	// Draw the collision box for debugging reasons.
	r.DrawRectangleLines(
		int(b.GetX()), int(b.GetY()),
		int(b.GetX()), int(b.Sprite.Height),
		r.Red,
	)

	enemyCenterBottom := r.NewVector2(
		float32(b.GetX())+float32(b.Sprite.Width/2),
		float32(b.GetY()+b.Sprite.Height),
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
