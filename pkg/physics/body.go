package physics

import (
	"github.com/damienfamed75/endorem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

type Body struct {
	Velocity r.Vector2

	onGround bool
	gravity  float32
	maxSpeed r.Vector2
	ground   *Space
	// ground   []r.Rectangle

	r.Rectangle
}

func NewBody(x, y, w, h, maxSpdX, maxSpdY float32) *Body {
	b := &Body{
		Velocity: r.NewVector2(0, 0),
		maxSpeed: r.NewVector2(maxSpdX, maxSpdY),
		gravity:  common.GlobalConfig.Game.Gravity * 40,
		ground:   NewSpace(),
	}

	b.Rectangle = r.NewRectangle(x, y, w, h)

	return b
}

func (b *Body) SetGravity(g float32) {
	b.gravity = g
}

func (b *Body) OnGround() bool {
	return b.onGround
}

func (b *Body) SetGround(ground *Space) {
	b.ground = ground
}

func (b *Body) AddGround(ground ...Shape) {
	b.ground.Add(ground...)
}

func (b *Body) Update(dt float32) {
	b.Velocity.Y += b.gravity * dt

	if b.Velocity.X > b.maxSpeed.X {
		b.Velocity.X = b.maxSpeed.X
	}
	if b.Velocity.X < -b.maxSpeed.X {
		b.Velocity.X = -b.maxSpeed.X
	}

	if b.Velocity.Y > b.maxSpeed.Y {
		b.Velocity.Y = b.maxSpeed.Y
	}
	if b.Velocity.Y < -b.maxSpeed.Y {
		b.Velocity.Y = -b.maxSpeed.Y
	}

	tmpXRec := b.Move(b.Velocity.X, 0)
	tmpYRec := b.Move(0, b.Velocity.Y)

	for i := range *b.ground {

		if (*b.ground)[i].Overlaps(tmpXRec) && b.Velocity.X > 0 {
			overlap := (*b.ground)[i].GetOverlapRec(tmpXRec)
			b.Velocity.X -= overlap.Width
		} else if (*b.ground)[i].Overlaps(tmpXRec) && b.Velocity.X < 0 {
			overlap := (*b.ground)[i].GetOverlapRec(tmpXRec)
			b.Velocity.X += overlap.Width
		}

		if (*b.ground)[i].Overlaps(tmpYRec) && b.Velocity.Y > 0 {
			overlap := (*b.ground)[i].GetOverlapRec(tmpYRec)
			b.Velocity.Y -= overlap.Height
			b.onGround = true
		} else if (*b.ground)[i].Overlaps(tmpYRec) && b.Velocity.Y < 0 {
			overlap := (*b.ground)[i].GetOverlapRec(tmpYRec)
			b.Velocity.Y += overlap.Height
			b.onGround = true
		}

		if b.Velocity.Y < -b.gravity {
			b.onGround = false
		}

	}

	b.Rectangle.X += b.Velocity.X
	b.Rectangle.Y += b.Velocity.Y
}
