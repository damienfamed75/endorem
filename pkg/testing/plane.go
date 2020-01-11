package testing

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Plane is a surface for testing features for the game. It's very barebones and
// only includes a collision box that directly affects the visual shape also.
type Plane struct {
	// Collision affects the shape of the rectangle directly since it's a testing
	// object, there is no reason why they should be separate.
	Collision *resolv.Rectangle
	Color     r.Color
}

// NewPlane returns the default shape of the testing plane which is meant for an
// 800x600 display.
func NewPlane() *Plane {
	return &Plane{
		Collision: resolv.NewRectangle(0, 500, 800, 100),
		Color:     r.Orange,
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Plane) Draw() {
	rec := r.NewRectangle(
		float32(p.Collision.X),
		float32(p.Collision.Y),
		float32(p.Collision.W),
		float32(p.Collision.H),
	)

	r.DrawRectangleLinesEx(rec, 2, p.Color)
}
