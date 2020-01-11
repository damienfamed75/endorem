package scene

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"
	r "github.com/lachee/raylib-goplus/raylib"
)

// MenuScene is the game's initial scene before a player
// begins their run
type MenuScene struct {
	player *player.Player

	ground *resolv.Space
	world  *resolv.Space
	camera r.Camera2D
}

//Preload is used to load in assets and entities
func (s *MenuScene) Preload() {
	s.world = resolv.NewSpace()
	s.ground = resolv.NewSpace()

	//Add all ground and walls to spaces
	s.ground.Add(
		testing.NewPlane(0, 250, 250, 50),
		testing.NewPlane(300, 250, 100, 50),
		testing.NewPlane(0, 0, 50, 200),
		testing.NewPlane(350, 0, 50, 250),
	)

	// Create player and camera
	s.player = player.NewPlayer(200, 268, func() {}, s.ground)
	defaultZoom := common.GlobalConfig.Game.Camera.DefaultZoom
	s.camera = r.Camera2D{
		Zoom: defaultZoom,
	}

	s.world.Add(s.ground, s.player)
}

// Update frames
func (s *MenuScene) Update(dt float32) {
	// s.camera.Update(s.player.Update())
	s.player.Update()
}

// Draw frames
func (s *MenuScene) Draw() {
	r.BeginMode2D(s.camera)
	r.ClearBackground(r.Black)

	// Draw ground elements.
	for i := range *s.ground {
		(*s.ground)[i].(Drawer).Draw()
	}

	s.player.Draw()

	r.EndMode2D()
}

//Unload everything in MenuScene
func (s *MenuScene) Unload() {

}

func (s *MenuScene) String() string {
	return "menu scene"
}
