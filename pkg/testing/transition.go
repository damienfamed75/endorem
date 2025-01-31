package testing

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

const (
	ExitTransition  string = "exitGame"
	SceneTransition        = "newScene"
)

type Transition struct {
	Color  r.Color
	Width  int32
	Height int32

	collided bool
	*resolv.Space
}

func NewTransition(x, y, w, h int32, TransitionType string) *Transition {
	transitionSpace := resolv.NewSpace()

	transitionSpace.Add(
		resolv.NewRectangle(x, y, w, h),
	)

	return &Transition{
		Space: transitionSpace,
		Color: r.Green,

		Width:  w,
		Height: h,
	}
}

func (t *Transition) Draw() {
	x, y := t.Space.GetXY()

	rec := r.NewRectangle(
		float32(x),
		float32(y),
		float32(t.Width),
		float32(t.Height),
	)

	r.DrawRectangleLinesEx(rec, 2, t.Color)
}
