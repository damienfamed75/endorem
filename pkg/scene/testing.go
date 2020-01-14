package scene

import (
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// TestingScene is a testing scene for scene management.
type TestingScene struct {
	player *player.Player

	enemies *resolv.Space
	ground  *resolv.Space
	world   *resolv.Space
	camera  *common.EndoCamera
}

// Preload is used to load in assets and entities
func (s *TestingScene) Preload() {
	s.world = resolv.NewSpace()
	s.ground = resolv.NewSpace()

	s.enemies = resolv.NewSpace()

	// Add all ground to the ground space.
	s.ground.Add(
		testing.NewPlane(0, 500, 800, 100, r.Orange),
		testing.NewPlane(200, 400, 80, 70, r.DarkGreen),
	)
	s.ground.AddTags(common.TagGround)

	// Add the ground elements to the world space.
	s.world.Add(s.ground)

	s.player = player.NewPlayer(0, 468, func() {}, s.ground)
	s.player.AddTags(common.TagPlayer)

	s.camera = common.NewEndoCamera(s.player.Collision)
	s.camera.Zoom = 1

	// Add the player to the world space.
	s.world.Add(s.player)

	// Add enemies to the enemy space. Must be of common.Entity
	s.enemies.Add(
		enemy.NewBasic(100, 468, s.player.Collision),
		enemy.NewSlime(300, 400, s.world),
	)

	// Add enemies to the world space.
	s.world.Add(s.enemies)
}

// Update frames
func (s *TestingScene) Update(dt float32) {
	// Update the camera and player.
	s.camera.Update(s.player.Update())

	// Update all the enemies.
	for i := range *s.enemies {
		(*s.enemies)[i].(common.Entity).Update(dt)
	}

	// Loop through all the enemies and detect collisions with the player.
	for _, en := range *s.enemies {
		if en.GetData() == nil {
			continue
		}

		// Check the type of the enemy space data.
		// If it's a string, then it's a Hitbox.
		// If it's a reference to itself then it's a Hurtbox.
		switch t := en.GetData().(type) {
		case *enemy.Slime: // Hurtbox
			// If the hurtbox is colliding a player hitbox then take damage.
			if t.FilterByTags(enemy.TagHurtbox).IsColliding(s.player.FilterByTags(player.TagHitbox)) {
				t.TakeDamage()
				// If the player is colliding with the enemy then they should take damage.
			} else if s.player.FilterByTags(player.TagHurtbox).IsColliding(t.FilterOutByTags(enemy.TagAttackZone)) {
				s.player.TakeDamage()
			}
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

	// Draw all ground elements.
	for i := range *s.ground {
		(*s.ground)[i].(Drawer).Draw()
	}

	// Draw all the enemies.
	for i := range *s.enemies {
		(*s.enemies)[i].(common.Entity).Draw()
	}

	s.player.Draw()

	r.EndMode2D()
}

// Unload everything in TestingScene
func (s *TestingScene) Unload() {
}

// String returns name of TestingScene
func (s *TestingScene) String() string {
	return "testing scene"
}
