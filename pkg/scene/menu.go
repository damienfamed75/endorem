package scene

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/player"
)

// MenuScene is the game's initial scene before a player
// begins their run
type MenuScene struct {
	player *player.Player

	ground *resolv.Space
	world  *resolv.Space
	camera *common.EndoCamera
}

//Preload is used to load in assets and entities
func (s *MenuScene) Preload() {

}

// Update frames
func (s *MenuScene) Update(dt float32) {

}

// Draw frames
func (s *MenuScene) Draw() {

}

//Unload everything in MenuScene
func (s *MenuScene) Unload() {

}

func (s *MenuScene) String() string {
	return "menu scene"
}
