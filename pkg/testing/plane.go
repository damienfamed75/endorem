package testing

import (
	"github.com/damienfamed75/endorem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Plane is a surface for testing features for the game. It's very barebones and
// only includes a collision box that directly affects the visual shape also.
type Plane struct {
	// Collision affects the shape of the rectangle directly since it's a testing
	// object, there is no reason why they should be separate.
	Color  r.Color
	Width  int32
	Height int32

	*physics.Rectangle
	// Collision r.Rectangle
	// *resolv.Space
}

// NewPlane returns the default shape of the testing plane which is meant for an
// 800x600 display.
func NewPlane(x, y, w, h int32, color r.Color) *Plane {
	return &Plane{
		Color:     color,
		Rectangle: physics.NewRectangle(float32(x), float32(y), float32(w), float32(h)),

		Width:  w,
		Height: h,
	}
}

// Draw creates a rectangle using Raylib and draws the outline of it.
func (p *Plane) Draw() {
	// x, y := p.Space.GetXY()

	// rec := r.NewRectangle(
	// 	float32(x),
	// 	float32(y),
	// 	float32(p.Width),
	// 	float32(p.Height),
	// )

	// r.DrawRectangleLinesEx(rec, 2, p.Color)
	r.DrawRectangleLinesEx(p.Rectangle.Rectangle, 2, p.Color)
}
