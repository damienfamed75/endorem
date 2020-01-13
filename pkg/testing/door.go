package testing

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Door struct {
	Color  r.Color
	X, Y   int32
	X2, Y2 int32

	*resolv.Space
}

func NewDoor(x, y, x2, y2 int32, color r.Color) *Door {
	doorSpace := resolv.NewSpace()

	doorSpace.Add(
		resolv.NewLine(x, y, x2, y2),
	)

	doorSpace.AddTags("door")

	return &Door{
		X: x, Y: y,
		X2: x2, Y2: y2,
		Color: color,
		Space: doorSpace,
	}
}

func (d *Door) Draw() {
	r.DrawLineEx(
		r.NewVector2(float32(d.X), float32(d.Y)),
		r.NewVector2(float32(d.X2), float32(d.Y2)),
		5, d.Color,
	)
}
