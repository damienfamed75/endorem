package scene

import (
	"fmt"
	"log"

	"github.com/SolarLune/resolv/resolv"

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

	music        *r.Music
	musicTrigger *testing.Plane
	exitTrigger  *testing.Plane

	heartFull  r.Texture2D
	heartEmpty r.Texture2D

	boss   *resolv.Space
	ground *physics.Space
	world  *physics.Space

	camera *common.EndoCamera
}

// Preload is used to load in assets and entities
func (d *DemoScene) Preload() {
	r.InitAudioDevice()

	d.music = r.LoadMusicStream("assets/sounds/Broke_For_Free_-_01_-_Night_Owl.mp3")
	d.music.SetLoopCount(5)
	d.music.SetVolume(0.7)
	// d.music = r.LoadSound("assets/sounds/Broke_For_Free_-_01_-_Night_Owl.mp3")

	d.ground = physics.NewSpace()
	d.world = physics.NewSpace()
	d.boss = resolv.NewSpace()

	d.heartFull = r.LoadTexture("assets/heart1.png")
	d.heartEmpty = r.LoadTexture("assets/heart2.png")

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

		// invisible walle
		testing.NewPlane(850, 0, 50, 200, r.Yellow),
		testing.NewPlane(150, 0, 10, 200, r.Yellow),
	)
	d.ground.AddTags(common.TagGround)

	d.musicTrigger = testing.NewPlane(277, 500, 200, 50, r.Gold)
	d.exitTrigger = testing.NewPlane(160, 0, 10, 300, r.Red)

	// Add ground elements to the world space.
	d.world.Add(*d.ground...)

	// Create player & camera
	d.player = player.NewPlayer(200, 50, func() {}, d.ground)
	d.player.AddTags(common.TagPlayer)

	d.camera = common.NewEndoCamera(d.player.Collision)
	d.player.AddTags(common.TagPlayer)

	// Add player to world space.
	d.world.Add(*d.player.Space...)

	// Add enemies and boss to space
	d.boss.Add(
		enemy.NewFungalBoss(700, 1000-96-128, d.world),
	)
	// Add enemies and boss to space

}

// Update frames
func (d *DemoScene) Update(dt float32) {
	// Update the camera and player.
	d.camera.Update(d.player.Update(dt))

	d.music.UpdateStream()
	// d.slime.Update(dt)

	if d.musicTrigger.Overlaps(d.player.RayRec()) {
		d.music.PlayStream()
		// var dir float32
		// if d.player.RayRec().X > d.slime.RayRec().X {
		// 	dir = 1
		// } else {
		// 	dir = -1
		// }

		// d.player.TakeDamage(dir)
	}
	// if d.exitTrigger.Overlaps(d.player.RayRec()) {
	// d.Unload()
	// r.CloseWindow()
	// exit
	// }
	for i := range *d.boss {
		(*d.boss)[i].(common.Entity).Update(dt)
	}
}

// Draw frames
func (d *DemoScene) Draw() {
	r.BeginMode2D(d.camera.Camera2D)
	// r.ClearBackground(r.Gray)
	r.ClearBackground(r.NewColor(31, 14, 28, 255))

	r.DrawTexture(d.Background, 0, 0, r.White)
	r.DrawTextureEx(d.Background, r.NewVector2(float32(d.Background.Width), 0), 180, 1, r.White)
	r.DrawTexture(d.Foreground, 0, 0, r.White)

	d.player.Draw()

	// d.slime.Draw()

	d.debugDraw()
	// Draw all ground elements
	// for i := range *d.ground {
	// 	(*d.ground)[i].(Drawer).Draw()
	// }

	// Draw boss
	for i := range *d.boss {
		(*d.boss)[i].(common.Entity).Draw()
	}

	for i := 0; i < d.player.MaxHealth; i++ {
		fmt.Printf("HELT: [%v] of [%v]\n", i, d.player.MaxHealth)
		var heart r.Texture2D
		if i > d.player.Health {
			heart = d.heartEmpty
			// r.DrawTexture(d.heartEmpty, int(d.heartEmpty.Width)*i, r.GetScreenHeight()-int(d.heartEmpty.Height), r.White)
		} else {
			heart = d.heartFull
			// r.DrawTexture(d.heartFull, int(d.heartFull.Width)*i, r.GetScreenHeight()-int(d.heartFull.Height), r.White)
		}

		// pos := r.GetWorldToScreen2D(
		// pos := r.GetScreenToWorld2D(
		// 	r.NewVector2(float32(int(heart.Width)*i), float32(r.GetScreenHeight()-int(heart.Height))),
		// 	d.camera.Camera2D,
		// )
		pos := r.NewVector2(
			d.player.RayRec().MinPosition().X+float32((int(heart.Width/2)*i))-d.player.RayRec().Width/4,
			d.player.RayRec().MaxPosition().Y+float32(heart.Height/2),
		)

		r.DrawTextureEx(heart, pos, 0, 0.5, r.White)
		// r.DrawTexture(heart, int(d.camera.Position.X)+int(heart.Width)*(i+1), int(d.camera.Position.Y)+int(heart.Height), r.White)
		// r.DrawTexture(heart, d.player.RayRec().X*(i+1), int(d.camera.Offset.Y)+int(heart.Height), r.White)
		// r.DrawTexture(heart, int(d.camera.Offset.X)+int(heart.Width)*(i+1), int(d.camera.Offset.Y)+int(heart.Height), r.White)
		// r.DrawTexture(heart, int(heart.Width)*(i+1), r.GetScreenHeight()-int(heart.Height), r.White)
	}

	r.EndMode2D()
}

func (d *DemoScene) Unload() {
	d.player.Sprite.Unload()
	d.player.MaskObj.Sprite.Unload()
	d.Background.Unload()
	d.Foreground.Unload()
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
