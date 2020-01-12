package scene

import (
	"github.com/SolarLune/dngn"
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/testing"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ Scene = &LevelOne{}
)

type LevelOne struct {
	mapData     *dngn.Room
	rooms       []common.RoomSpec
	player      *player.Player
	ground      *resolv.Space
	world       *resolv.Space
	camera      *common.EndoCamera
	overviewCam r.Camera2D
}

func (l *LevelOne) Preload() {
	l.overviewCam = r.Camera2D{
		Rotation: 0,
		Zoom:     0.35,
	}

	l.world = resolv.NewSpace()
	l.ground = resolv.NewSpace()

	l.mapData, l.rooms = common.GenerateMap(1)
	mapScale := 34

	spawnRoom := l.rooms[0]

	x, y := (spawnRoom.X*mapScale)+int(34), (spawnRoom.Y*mapScale)+int(34*2)

	l.player = player.NewPlayer(x, y, func() {}, l.ground, resolv.NewSpace())
	l.camera = common.NewEndoCamera(l.player.Collision)

	l.mapData.Select().By(func(x, y int) bool {
		switch l.mapData.Get(x, y) {
		case '#': // Wall
			l.ground.Add(
				testing.NewSolidPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.Aqua,
				),
			)
		case '-': // Door
			l.ground.Add(
				testing.NewPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.Orange,
				),
			)
		case '^': // Hatches
			l.ground.Add(
				testing.NewPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.Gold,
				),
			)
		case '=': // Floating Platform 1
			l.ground.Add(
				testing.NewSolidPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.GopherBlue,
				),
			)
		case '~': // Floating Platform 2
			l.ground.Add(
				testing.NewSolidPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.SkyBlue,
				),
			)
		}

		return false
	})

	l.world.Add(l.ground, l.player)
}

func (l *LevelOne) Update(dt float32) {
	l.camera.Update(l.player.Update())
	// l.camera.Offset.Y += 2
}

func (l *LevelOne) Draw() {
	// r.BeginMode2D(l.overviewCam)
	r.BeginMode2D(l.camera.Camera2D)
	r.ClearBackground(r.Black)

	for i := range *l.ground {
		(*l.ground)[i].(Drawer).Draw()
	}
	l.player.Draw()
	r.DrawCircleLines(int(l.camera.Target.X), int(l.camera.Target.Y), 50, r.Red)
	// r.DrawCircle(int(l.camera.Offset.X), int(l.camera.Offset.Y), 50, r.Red)

	r.EndMode2D()
}

func (l *LevelOne) Unload() {

}

func (l *LevelOne) String() string {
	return "level one"
}
