package testing

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Platform struct {
	Color  r.Color
	Width  int32
	Height int32

	*resolv.Space
}

func NewPlatform(x, y int32, color r.Color) *Platform {
	platformSpace := resolv.NewSpace()

	platformSpace.Add(
		resolv.NewRectangle(x, y, 50, 5),
	)
	return &Platform{
		Space: platformSpace,
		Color: color,

		Width:  50,
		Height: 5,
	}
}

func (p *Platform) Draw() {
	x, y := p.Space.GetXY()

	r.DrawRectangle(
		int(x), int(y),
		int(p.Width), int(p.Height),
		p.Color,
	)
}
