// +build !windows

// Until we find out that linux breaks, this version of the camera will build on
// anything that's NOT Windows.

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

// NewEndoCamera creates a default offset of the player's position.
func NewEndoCamera(playerColl *resolv.Rectangle) *EndoCamera {
	defaultZoom := GlobalConfig.Game.Camera.DefaultZoom

	// Get the center coordinates of the player collision
	cx, cy := playerColl.Center()
	xOff, yOff := -float32(cx)*defaultZoom, -float32(cy)*defaultZoom

	return &EndoCamera{
		Camera2D: r.Camera2D{
			Offset: r.NewVector2(
				xOff+float32(r.GetScreenWidth()),
				yOff+float32(r.GetScreenHeight()),
			),
			Rotation: 0,
			Zoom:     defaultZoom,
		},
		LerpAmount: GlobalConfig.Game.Camera.DefaultSpeed,
	}
}

// Update changes the offset position of the camera and the target.
func (e *EndoCamera) Update(curr r.Vector2) {
	// Update camera offset coordinates for it to move.
	xOff, yOff := -float32(curr.X+4)*e.Zoom, -float32(curr.Y+8)*e.Zoom
	e.Offset = r.NewVector2(
		xOff+float32(r.GetScreenWidth()),
		yOff+float32(r.GetScreenHeight()),
	)

	// Reset the camera's target to the player's current position.
	// Using a lerp to make the camera movement smoother.
	e.Target = e.Target.Lerp(curr, e.LerpAmount)
}
