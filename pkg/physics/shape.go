package physics

import r "github.com/lachee/raylib-goplus/raylib"

type Shape interface {
	GetData() interface{}
	SetData(dat interface{})
	GetTags() []string
	HasTags(tags ...string) bool
	AddTags(tags ...string)
	RemoveTags(tags ...string)
	ClearTags()
	RayRec() r.Rectangle

	// raylib api
	// MaxPosition() r.Vector2
	// MinPosition() r.Vector2
	Position() r.Vector2
	// Center() r.Vector2
	// Size() r.Vector2
	Overlaps(r.Rectangle) bool
	// GetOverlapRec(r.Rectangle) r.Rectangle
	// LerpPosition(pos r.Vector2, amount float32) r.Rectangle
	// SetPosition(r.Vector2) r.Rectangle
	// Move(x, y float32) r.Rectangle
}
