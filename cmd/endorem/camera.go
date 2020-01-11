package main

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type EndoCamera r.Camera2D

// NewEndoCamera creates a default offset of the player's position.
func NewEndoCamera(playerColl *resolv.Rectangle) *EndoCamera {
	xOff, yOff := playerColl.X-(playerColl.W/2), playerColl.Y-(playerColl.H/2)
	return &EndoCamera{
		Offset: r.NewVector2(
			float32(xOff+int32(r.GetScreenWidth())),
			float32(yOff-int32(r.GetScreenHeight()/2)),
		),
		Rotation: 0,
		Zoom:     1,
	}
}

// Update changes the offset position of the camera and the target.
func (e *EndoCamera) Update(diff, curr r.Vector2) {
	e.Offset.X -= diff.X
	e.Offset.Y -= diff.Y

	e.Target = curr
}
