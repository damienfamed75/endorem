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

func (r *Rigidbody) Update(dt float32) {

	// r.Velocity.Y *= dt
	// r.Velocity.X *= dt
	// Gravity
	r.Velocity.Y += r.gravity

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

	// Misc.
	x, y := int32(r.Velocity.X), int32(r.Velocity.Y)

	// Ground check
	down := r.ground.Resolve(r.collisions, 0, y+1)
	// down := r.ground.Resolve(r.collisions, 0, y+1)
	r.onGround = down.Colliding()

	// X
	if res := r.ground.Resolve(r.collisions, x, 0); res.Colliding() {
		// if res.ResolveX > r.maxSpeedX {
		// 	res.ResolveX = r.maxSpeedX
		// }
		// if res.ResolveX < -r.maxSpeedX {
		// 	res.ResolveX = -r.maxSpeedX
		// }
		if res.Teleporting {
			// fmt.Printf("teleport(x) resx[%v] spdx[%v]\n", res.ResolveX, x)
			// if res.ResolveX > 32 || res.ResolveX < -32 {
			// 	res = r.ground.Resolve(r.collisions, 0, y)
			// 	fmt.Printf("stuck check resy[%v] spdy[%v]\n", res.ResolveY, y)
			// 	r.Velocity.Y = 0
			// 	r.Move(0, res.ResolveY)
			// }
			res.ResolveX = -x
			x = 0
			r.Velocity.X = float32(-x)
		} else {
			// fmt.Printf("resolve (x) resx[%v] spdx[%v]\n", res.ResolveX, x)
			// x *= -1
			x = res.ResolveX
		}

		res.ShapeA.Move(res.ResolveX, 0)
		// res.ShapeB.Move(-res.ResolveX, 0)

		// x = res.ResolveX
		// r.Velocity.X = 0
	}

	r.Space.Move(x, 0)

	// Y
	if res := r.ground.Resolve(r.collisions, 0, y); res.Colliding() {
		// if res.ResolveY > r.maxSpeedY {
		// 	res.ResolveY = r.maxSpeedY
		// }
		// if res.ResolveY < -r.maxSpeedY {
		// 	res.ResolveY = -r.maxSpeedY
		// }

		// if res.Teleporting {
		// 	fmt.Printf("teleport(y) resy[%v] spdy[%v]\n", res.ResolveY, y)
		// 	res.ResolveY = -y
		// }

		// y = res.ResolveY
		// r.Velocity.Y = 0

		if res.Teleporting {
			// fmt.Printf("teleport(y) resy[%v] spdy[%v]\n", res.ResolveY, y)
			res.ResolveY = -y
			y = 0
			r.Velocity.Y = float32(-y)
		} else {
			// fmt.Printf("resolve (y) resy[%v] spdy[%v]\n", res.ResolveY, y)
			// x *= -1
			y = res.ResolveY
		}

		res.ShapeA.Move(0, res.ResolveY)
	}

	r.Space.Move(0, y)
}
