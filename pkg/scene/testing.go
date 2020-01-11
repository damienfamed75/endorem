package scene

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"

	r "github.com/lachee/raylib-goplus/raylib"
)

// TestingScene is a testing scene for scene management.
type TestingScene struct {
	enemy  *enemy.Basic
	player *player.Player

	ground *resolv.Space
	world  *resolv.Space
	camera *common.EndoCamera
}

// Preload is used to load in assets and entities
func (s *TestingScene) Preload() {
	s.world = resolv.NewSpace()
	s.ground = resolv.NewSpace()

	s.ground.Add(
		testing.NewPlane(),
	)

	s.enemy = enemy.NewBasic(100, 468)
	s.player = player.NewPlayer(0, 468, func() {})
	s.camera = common.NewEndoCamera(s.player.Collision)

	s.world.Add(s.ground)
	s.world.Add(s.player, s.enemy)
}

// Update frames
func (s *TestingScene) Update(dt float32) {
	s.camera.Update(s.player.Update(s.ground))
	s.enemy.Update()

	enemies := s.world.FilterByTags(common.TagEnemy)

	for _, en := range *enemies {
		if en.GetData() == nil {
			continue
		}

		// Check the type of the enemy space data.
		// If it's a string, then it's a Hitbox.
		// If it's a reference to itself then it's a Hurtbox.
		switch t := en.GetData().(type) {
		case *enemy.Basic: // Hurtbox
			enX, enY := t.Collision.Center()
			pX, pY := s.player.Collision.Center()

			// Calculate the distance from the enemy to the player.
			dist := resolv.Distance(enX, enY, pX, pY)

			t.PlayerSeen = dist < common.GlobalConfig.Enemy.VisionDistance
			t.ShouldAttack = dist < t.AttackDistance

			// If the hurtbox is colliding a player hitbox then take damage.
			if t.FilterByTags(enemy.TagHurtbox).IsColliding(s.player.Hitbox) {
				t.TakeDamage()
				// If the player is colliding with the enemy then they should take damage.
			} else if s.player.FilterByTags(player.TagHurtbox).IsColliding(t.FilterOutByTags(enemy.TagAttackZone)) {
				s.player.TakeDamage()
			}
		}
	}
}

// Draw frames
func (s *TestingScene) Draw() {
	r.BeginMode2D(s.camera.Camera2D)
	r.ClearBackground(r.Black)

	for _, g := range *s.ground {
		g.(Drawer).Draw()
	}

	s.player.Draw()
	s.enemy.Draw()

	r.EndMode2D()
}

// Unload everything in TestingScene
func (s *TestingScene) Unload() {
}

// String returns name of TestingScene
func (s *TestingScene) String() string {
	return "testing scene"
}
