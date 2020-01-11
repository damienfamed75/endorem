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
	walls  *resolv.Space
	world  *resolv.Space
	camera *common.MenuCamera
}

//Preload is used to load in assets and entities
func (s *MenuScene) Preload() {
	s.world = resolv.NewSpace()
	s.ground = resolv.NewSpace()

	//Add all ground and walls to spaces
	s.ground.Add(
		testing.NewPlane(0, 250, 250, 50),
		testing.NewPlane(300, 250, 100, 50),
	)

	// s.walls.Add()

	// Create player and camera
	s.player = player.NewPlayer(0, 268, func() {}, s.ground)
	s.camera = common.NewMenuCamera()

	s.world.Add(s.ground, s.player)
}

// Update frames
func (s *MenuScene) Update(dt float32) {
	// s.camera.Update(s.player.Update())
	s.player.Update()
}

// Draw frames
func (s *MenuScene) Draw() {
	r.BeginMode2D(s.camera.Camera2D)
	r.ClearBackground(r.Black)

	// Draw ground elements.
	for i := range *s.ground {
		(*s.ground)[i].(Drawer).Draw()
	}

	// Draw wall elements.
	// for i := range *s.walls {
	// 	(*s.walls)[i].(Drawer).Draw()
	// }

	s.player.Draw()

	r.EndMode2D()
}

//Unload everything in MenuScene
func (s *MenuScene) Unload() {

}

func (s *MenuScene) String() string {
	return "menu scene"
}
