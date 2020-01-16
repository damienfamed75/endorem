package scene

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/physics"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"

	r "github.com/lachee/raylib-goplus/raylib"
)

// TestingScene is a testing scene for scene management.
type TestingScene struct {
	player *player.Player

	g1 *testing.Plane
	g2 *testing.Plane
	g3 *testing.Plane

	ground  *physics.Space
	enemies *physics.Space
	boss    *resolv.Space

	world *physics.Space

	camera *common.EndoCamera
}

// Preload is used to load in assets and entities
func (s *TestingScene) Preload() {
	s.ground = physics.NewSpace()
	s.enemies = physics.NewSpace()
	s.world = physics.NewSpace()
	s.boss = resolv.NewSpace()

	// Add all ground to the ground space.
	s.ground.Add(
		testing.NewPlane(0, 500, 800, 100, r.Orange),
		testing.NewPlane(200, 450, 50, 50, r.Green),
		testing.NewPlane(400, 400, 100, 50, r.Green),
	)
	s.ground.AddTags(common.TagGround)

	// Add the ground elements to the world space.
	s.world.Add(*s.ground...)

	s.player = player.NewPlayer(0, 468, func() {}, s.ground)
	s.player.AddTags(common.TagPlayer)

	// s.player.Body.AddGround(*s.ground...)
	s.camera = common.NewEndoCamera(s.player.Collision)
	s.camera.Zoom = 1

	// Add the player to the world space.
	s.world.Add(*s.player.Space...)

	// Add enemies to the enemy space. Must be of common.Entity

	s.enemies.Add(
		//enemy.NewBasic(100, 468, s.world),
		enemy.NewSlime(300, 400, s.world),
	)

	s.boss.Add(
		enemy.NewFungalBoss(700, 468, s.world),
	)
	s.enemies.AddTags(common.TagEnemy)

	// Add enemies to the world space.
	s.world.Add(*s.enemies...)
}

// Update frames
func (s *TestingScene) Update(dt float32) {
	// Update the camera and player.
	s.camera.Update(s.player.Update(dt))

	// Update all the enemies.
	for i := range *s.enemies {
		(*s.enemies)[i].(common.Entity).Update(dt)
	}

	for i := range *s.boss {
		(*s.boss)[i].(common.Entity).Update(dt)
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
			if t.Collider().Overlaps(s.player.Hitbox.Rectangle) {
				// if t.FilterByTags(enemy.TagHurtbox).IsColliding(s.player.Hitbox) {
				t.TakeDamage()
				// If the player is colliding with the enemy then they should take damage.
				// } else if s.player.FilterByTags(player.TagHurtbox).IsColliding(t.FilterOutByTags(enemy.TagAttackZone)) {
			}

			if s.player.Collider().Overlaps(t.Body.Collider()) {
				s.player.TakeDamage()
			}
		case *enemy.Basic: // Hurtbox
			// If the hurtbox is colliding a player hitbox then take damage.
			// if t.FilterByTags(enemy.TagHurtbox).IsColliding(s.player.FilterByTags(player.TagHitbox)) {
			// 	t.TakeDamage()
			// 	// If the player is colliding with the enemy then they should take damage.
			// } else if s.player.FilterByTags(player.TagHurtbox).IsColliding(t.FilterOutByTags(enemy.TagAttackZone)) {
			// 	s.player.TakeDamage()
			// }
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

	// Draw boss
	for i := range *s.boss {
		(*s.boss)[i].(common.Entity).Draw()
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
