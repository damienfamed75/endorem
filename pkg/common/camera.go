package common

import (
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// EndoCamera is a custom name for the raylib Camera2D
type EndoCamera struct {
	LerpAmount float32
	r.Camera2D
}

// Menu camera is a static camera that doesn't follow the player
// could potentially be used as a
type MenuCamera struct {
	r.Camera2D
}

// type EndoCamera r.Camera2D

// NewEndoCamera creates a default offset of the player's position.
func NewEndoCamera(playerColl *resolv.Rectangle) *EndoCamera {
	// Get the center coordinates of the player collision
	xOff, yOff := playerColl.Center()
	defaultZoom := GlobalConfig.Game.Camera.DefaultZoom
	offsetMultiplier := defaultZoom - 0.5 // 0.5 tries to center the Y of cam
	return &EndoCamera{
		Camera2D: r.Camera2D{
			Offset: r.NewVector2(
				float32(xOff+int32(r.GetScreenWidth())),
				float32(yOff-int32(float32(r.GetScreenHeight())*offsetMultiplier)),
			),
			Rotation: 0,
			Zoom:     defaultZoom,
		},
		LerpAmount: GlobalConfig.Game.Camera.DefaultSpeed,
	}
}

// Update changes the offset position of the camera and the target.
func (e *EndoCamera) Update(diff, curr r.Vector2) {
	// Update camera offset coordinates for it to move.
	e.Offset.X -= diff.X * e.Zoom
	e.Offset.Y -= diff.Y * e.Zoom

	// Reset the camera's target to the player's current position.
	// Using a lerp to make the camera movement smoother.
	e.Target = e.Target.Lerp(curr, 0.1)
}

func NewMenuCamera() *MenuCamera {
	defaultZoom := GlobalConfig.Game.Camera.DefaultZoom
	return &MenuCamera{
		Camera2D: r.Camera2D{
			Zoom: defaultZoom,
		},
	}
}
