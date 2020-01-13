// +build windows

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

	return &EndoCamera{
		Camera2D: r.Camera2D{
			// Center the camera on the player's position.
			Offset: r.NewVector2(
				float32(r.GetScreenWidth()/2)-float32(playerColl.W/2),
				float32(r.GetScreenHeight()/2)-float32(playerColl.H/2),
			),
			Rotation: 0,
			Zoom:     defaultZoom,
		},
		LerpAmount: GlobalConfig.Game.Camera.DefaultSpeed,
	}
}

// Update changes the offset position of the camera and the target.
func (e *EndoCamera) Update(diff, curr r.Vector2) {
	// Note: For Windows we don't need the camera offset to change.

	// Update camera offset coordinates for it to move.
	e.Target = e.Target.Lerp(curr, 0.1)
}
