package physics

import (
	"github.com/damienfamed75/endorem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

type Body struct {
	Velocity r.Vector2

	onGround  bool
	gravity   float32
	maxSpeed  r.Vector2
	ground    *Space
	collision *Rectangle
	// ground   []r.Rectangle

	*Space
	// *Rectangle
	// r.Rectangle
}

func NewBody(x, y, w, h, maxSpdX, maxSpdY float32) *Body {
	b := &Body{
		Velocity:  r.NewVector2(0, 0),
		maxSpeed:  r.NewVector2(maxSpdX, maxSpdY),
		gravity:   common.GlobalConfig.Game.Gravity * 40,
		ground:    NewSpace(),
		Space:     NewSpace(),
		collision: NewRectangle(x, y, w, h),
	}

	b.Space.Add(b.collision)
	// b.Rectangle = NewRectangle(x, y, w, h)

	// b.Rectangle = r.NewRectangle(x, y, w, h)

	return b
}

func (b *Body) SetGravity(g float32) {
	b.gravity = g
}

func (b *Body) OnGround() bool {
	return b.onGround
}

func (b *Body) Collider() r.Rectangle {
	return b.collision.Rectangle
}

func (b *Body) Position() r.Vector2 {
	return b.collision.Rectangle.Position()
}

func (b *Body) SetGround(ground *Space) {
	b.ground = ground
}

func (b *Body) GetGround() *Space {
	return b.ground
}

func (b *Body) AddGround(ground ...Shape) {
	b.ground.Add(ground...)
}

func (b *Body) speedCheck() {
	// Cap player movement speed.
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
}

func (b *Body) ResolveForces() {
	tmpXRec := b.collision.Rectangle.Move(b.Velocity.X, 0)
	tmpYRec := b.collision.Rectangle.Move(0, b.Velocity.Y)

	// Limit the player to touching one object on each axis at a time.
	// This means that numbers won't get messed up when touching two
	// ground objects at the same time.
	var colx, coly bool
	for i := range *b.ground {

		// If the player hasn't collided with anything on the x-axis yet.
		if !colx {
			if (*b.ground)[i].Overlaps(tmpXRec) {
				overlap := (*b.ground)[i].RayRec().GetOverlapRec(tmpXRec)
				colx = true

				if b.Velocity.X > 0 {
					b.Velocity.X -= overlap.Width
				} else {
					b.Velocity.X += overlap.Width
				}
			}
		}

		// If the player hasn't collided with anything on the y-axis yet.
		if !coly {
			if (*b.ground)[i].Overlaps(tmpYRec) {
				overlap := (*b.ground)[i].RayRec().GetOverlapRec(tmpYRec)
				coly = true
				b.onGround = true

				if b.Velocity.Y > 0 {
					b.Velocity.Y -= overlap.Height
				} else {
					b.Velocity.Y += overlap.Height
				}
			}
		}
	}

	b.collision.Rectangle.X += b.Velocity.X
	b.collision.Rectangle.Y += b.Velocity.Y
}

func (b *Body) Update(dt float32) {
	// Add default gravity effect on velocity.
	b.Velocity.Y += b.gravity * dt

	b.speedCheck()

	// If the player's Y velocity is up and is greater than the default
	// gravity effect of the player then it should be considered on the ground.
	if b.Velocity.Y < -(b.gravity * dt) {
		b.onGround = false
	}

	b.ResolveForces()

	// b.Rectangle.X += b.Velocity.X
	// b.Rectangle.Y += b.Velocity.Y
}
