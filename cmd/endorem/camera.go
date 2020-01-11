package main

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// EndoCamera is a custom name for the raylib Camera2D
type EndoCamera r.Camera2D

// NewEndoCamera creates a default offset of the player's position.
func NewEndoCamera(playerColl *resolv.Rectangle) *EndoCamera {
	// Get the center coordinates of the player collision
	xOff, yOff := playerColl.Center()
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
	// Update camera offset coordinates for it to move.
	e.Offset.X -= diff.X
	e.Offset.Y -= diff.Y

	// Reset the camera's target to the player's current position.
	e.Target = curr
}
