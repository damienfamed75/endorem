package scene

import (
	"log"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/player"
	r "github.com/lachee/raylib-goplus/raylib"
)

// MenuScene is the game's initial scene before a player
// begins their run
type MenuScene struct {
	player *player.Player

	ground   *resolv.Space
	startRun *resolv.Space
	exit     *resolv.Space
	world    *resolv.Space
	camera   r.Camera2D
}

//Preload is used to load in assets and entities
func (s *MenuScene) Preload() {
	s.world = resolv.NewSpace()
	s.ground = resolv.NewSpace()
	s.startRun = resolv.NewSpace()
	s.exit = resolv.NewSpace()

	//Add all ground and walls to spaces
	// s.ground.Add(
	// 	// ground
	// 	testing.NewPlane(0, 250, 250, 50, r.Orange),
	// 	testing.NewPlane(300, 250, 100, 50, r.Orange),

	// 	//walls
	// 	testing.NewPlane(0, 0, 50, 200, r.Orange),
	// 	testing.NewPlane(350, 0, 50, 250, r.Orange),
	// )

	// // exits
	// s.startRun.Add(
	// 	testing.NewPlane(250, 290, 50, 10, r.Green),
	// )
	// s.exit.Add(
	// 	testing.NewPlane(0, 200, 10, 50, r.Green),
	// )

	// Create player and camera
	s.player = player.NewPlayer(200, 268, func() {}, s.ground)
	defaultZoom := common.GlobalConfig.Game.Camera.DefaultZoom
	s.camera = r.Camera2D{
		Zoom: defaultZoom,
	}

	s.world.Add(s.ground, s.startRun, s.exit, s.player)
}

// Update frames
func (s *MenuScene) Update(dt float32) {
	// s.camera.Update(s.player.Update())
	s.player.Update(dt)
	s.checkTransitions()
}

// Draw frames
func (s *MenuScene) Draw() {
	r.BeginMode2D(s.camera)
	r.ClearBackground(r.Black)

	// Draw ground elements.
	for i := range *s.ground {
		(*s.ground)[i].(Drawer).Draw()
	}
	// Draw transition elements.
	for i := range *s.exit {
		(*s.exit)[i].(Drawer).Draw()
	}
	for i := range *s.startRun {
		(*s.startRun)[i].(Drawer).Draw()
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

func (s *MenuScene) checkTransitions() {
	// Transition Check
	if s.player.IsColliding(s.exit) {
		// ADD GAME OVER
		log.Print("EXIT GAME")
	}

	if s.player.IsColliding(s.startRun) {
		//TRANSITION TO GAME START
		log.Print("START RUN")
	}
}
