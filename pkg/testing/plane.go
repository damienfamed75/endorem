package testing

import (
	"log"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Collision interface {
	HandleCollision(playerPos *resolv.Rectangle, playerHeight int32, speed *r.Vector2) (x, y int32)
}

// Plane is a surface for testing features for the game. It's very barebones and
// only includes a collision box that directly affects the visual shape also.
type Plane struct {
	// Collision affects the shape of the rectangle directly since it's a testing
	// object, there is no reason why they should be separate.
	Color  r.Color
	Width  int32
	Height int32

	*resolv.Space
}

var _ Collision = (*Plane)(nil)

// NewPlane returns the default shape of the testing plane which is meant for an
// 800x600 display.
func NewPlane(x, y, w, h int32, color r.Color) *Plane {
	planeSpace := resolv.NewSpace()
	planeSpace.SetData(planeSpace)
	planeSpace.Add(
		//resolv.NewRectangle(0, 500, 800, 100),
		resolv.NewRectangle(x, y, w, h),
	)

	return &Plane{
		Space: planeSpace,
		Color: color,

		Width:  w,
		Height: h,
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Plane) Draw() {
	x, y := p.Space.GetXY()

	rec := r.NewRectangle(
		float32(x),
		float32(y),
		float32(p.Width),
		float32(p.Height),
	)

	r.DrawRectangleLinesEx(rec, 2, p.Color)
}
func (p *Plane) HandleCollision(playerPos *resolv.Rectangle, playerHeight int32, speed *r.Vector2) (x, y int32) {
	playerX := int32(speed.X)
	playerY := int32(speed.Y)
	if res := p.Resolve(playerPos, playerX, 0); res.Colliding() {
		playerX = res.ResolveX
		speed.X = 0
	}

	res := p.Resolve(playerPos, 0, playerY+4)

	if playerY < 0 || (res.Teleporting && res.ResolveY < -playerHeight/2) {
		res = resolv.Collision{}
	}
	if !res.Colliding() {
		res = p.Resolve(playerPos, 0, playerY)
	}

	if res.Colliding() {
		playerY = res.ResolveY

		speed.Y = 0
	}
	log.Print("collide")
	return playerX, playerY
}
