package scene

import (
	"log"

	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/enemy"
	"github.com/damienfamed75/endorem/pkg/physics"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"

	r "github.com/lachee/raylib-goplus/raylib"
)

// DemoScene is the level that will be the presentation level
// for the Weekly Game Jam (131)
type DemoScene struct {
	Foreground r.Texture2D
	Background r.Texture2D
	player     *player.Player

	ground *physics.Space

	world *physics.Space

	slime *enemy.Slime

	camera *common.EndoCamera
}

// Preload is used to load in assets and entities
func (d *DemoScene) Preload() {
	r.InitAudioDevice()

	d.ground = physics.NewSpace()
	d.world = physics.NewSpace()

	d.Foreground = r.LoadTexture("assets/foreground.png")
	d.Background = r.LoadTexture("assets/background.png")

	log.Print(d.Foreground.Height)
	log.Print(d.Foreground.Width)
	// Add all ground to ground space.
	d.ground.Add(
		// left
		testing.NewPlane(0, 168, 277, 736, r.Orange),
		testing.NewPlane(278, 168, 47, 442, r.Orange),
		// down
		testing.NewPlane(0, 904, 1000, 96, r.Orange),

		// right
		testing.NewPlane(941, 168, 59, 736, r.Orange),
		testing.NewPlane(380, 168, 563, 442, r.Orange),
	)
	d.ground.AddTags(common.TagGround)

	// Add ground elements to the world space.
	d.world.Add(*d.ground...)

	// Create player & camera
	d.player = player.NewPlayer(100, 50, func() {}, d.ground)
	d.player.AddTags(common.TagPlayer)

	d.camera = common.NewEndoCamera(d.player.Collision)
	d.player.AddTags(common.TagPlayer)

	// Add player to world space.
	d.world.Add(*d.player.Space...)

	// Add enemies and boss to space
	// d.slime = enemy.NewSlime(400, 50, d.world)

	// Add enemies and boss to space

}

// Update frames
func (d *DemoScene) Update(dt float32) {
	// Update the camera and player.
	d.camera.Update(d.player.Update(dt))
	// d.slime.Update(dt)

	// if d.slime.Overlaps(d.player.RayRec()) {
	// 	// overlap := d.slime.RayRec().GetOverlapRec(d.player.RayRec())
	// 	var dir float32
	// 	if d.player.RayRec().X > d.slime.RayRec().X {
	// 		dir = 1
	// 	} else {
	// 		dir = -1
	// 	}

	// 	d.player.TakeDamage(dir)
	// }
}

// Draw frames
func (d *DemoScene) Draw() {
	r.BeginMode2D(d.camera.Camera2D)
	// r.ClearBackground(r.Gray)
	r.ClearBackground(r.NewColor(31, 14, 28, 255))

	r.DrawTexture(d.Background, 0, 0, r.White)
	r.DrawTexture(d.Foreground, 0, 0, r.White)

	d.player.Draw()
	// d.slime.Draw()

	d.debugDraw()
	// Draw all ground elements
	// for i := range *d.ground {
	// 	(*d.ground)[i].(Drawer).Draw()
	// }

	r.EndMode2D()
}

func (d *DemoScene) Unload() {

}

func (d *DemoScene) String() string {
	return "demo level"
}
func (d *DemoScene) debugDraw() {
	// Used to correlate collisions to textures
	// r.DrawTexture(d.Background, 0, 0, r.White)
	// r.DrawTexture(d.Foreground, 0, 0, r.White)

	// Draw ground collision boxes
	// for _, shape := range *d.ground{
	// 	x,y := shape.
	// 	r.DrawRectangleLines(
	// 		int(x)
	// 	)
	// }

}
