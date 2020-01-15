package player

import (
	"fmt"
	"strconv"
	
	"github.com/damienfamed75/endorem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

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

	px, py := p.GetXY()

	r.DrawText(
		fmt.Sprintf("P[%v,%v]", px, py),
		int(p.Collision.X), int(p.Collision.Y+(p.Collision.H)+20), 10,
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

	// r.DrawRectangleLines(
	// 	int(p.GetX()), int(p.GetY()),
	// 	int(p.SpriteStand.Width), int(p.SpriteStand.Height),
	// 	r.Red,
	// )

	r.DrawRectangleLinesEx(p.Body.Rectangle, 2, r.Red)

	if p.isAttacking {
		r.DrawRectangleLines(
			int(p.Hitbox.X), int(p.Hitbox.Y),
			int(p.Hitbox.W), int(p.Hitbox.H),
			r.Green,
		)
	}
}
