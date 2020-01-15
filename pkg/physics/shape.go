package physics

import r "github.com/lachee/raylib-goplus/raylib"

type Shape interface {
	HasTags(tags ...string) bool
	AddTags(tags ...string)
	ClearTags()

	// raylib api
	Move(x, y float32) r.Rectangle
	Overlaps(r.Rectangle) bool
	GetOverlapRec(r.Rectangle) r.Rectangle
}
