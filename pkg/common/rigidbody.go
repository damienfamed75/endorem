package common

// rigidbody
//
// rigidbody is a resolv.Space
//
// configurable colliders
// configurable tags that the body can collide with
// handles collisions and corrects itself.

import (
	r "github.com/lachee/raylib-goplus/raylib"

	"github.com/SolarLune/resolv/resolv"
)

type Rigidbody struct {
	onGround   bool
	ground     *resolv.Space
	collisions *resolv.Space
	gravity    float32

	Velocity  r.Vector2
	maxSpeedX int32
	maxSpeedY int32

	*resolv.Space
}

func setupRigidbody() *Rigidbody {
	return &Rigidbody{
		Space:   resolv.NewSpace(),
		gravity: GlobalConfig.Game.Gravity,
	}
}

func NewRigidbody(x, y, maxSpeedX, maxSpeedY int32, ground *resolv.Space, colliders ...resolv.Shape) *Rigidbody {
	r := setupRigidbody()

	r.Add(colliders...)
	r.collisions = r.FilterByTags(TagCollision)
	r.SetXY(x, y)
	r.maxSpeedX, r.maxSpeedY = maxSpeedX, maxSpeedY

	r.ground = ground

	return r
}

func (r *Rigidbody) GetX() int32 {
	x, _ := r.GetXY()
	return x
}

func (r *Rigidbody) GetY() int32 {
	_, y := r.GetXY()
	return y
}

func (r *Rigidbody) OnGround() bool {
	return r.onGround
}

func (r *Rigidbody) SetGravity(g float32) {
	r.gravity = g
}

func (r *Rigidbody) Update() {
	// Gravity
	if !r.onGround {
		r.Velocity.Y += r.gravity
	}
	// max speed checks
	if r.Velocity.X > float32(r.maxSpeedX) {
		r.Velocity.X = float32(r.maxSpeedX)
	}
	if r.Velocity.X < -float32(r.maxSpeedX) {
		r.Velocity.X = -float32(r.maxSpeedX)
	}
	if r.Velocity.Y > float32(r.maxSpeedY) {
		r.Velocity.Y = float32(r.maxSpeedY)
	}

	// Ground check
	down := r.ground.Resolve(r.collisions, 0, r.maxSpeedY)
	r.onGround = down.Colliding()

	// Misc.
	x, y := int32(r.Velocity.X), int32(r.Velocity.Y)

	// X
	if res := r.ground.Resolve(r.collisions, x, 0); res.Colliding() {
		x = res.ResolveX
		r.Velocity.X = 0
	}

	r.Space.Move(x, 0)

	// Y
	if res := r.ground.Resolve(r.collisions, 0, y); res.Colliding() {
		y = res.ResolveY
		r.Velocity.Y = 0
	}

	r.Space.Move(0, y)
}
