package scene

import (
	"fmt"

	"github.com/SolarLune/dngn"
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/physics"
	"github.com/damienfamed75/endorem/pkg/player"
	"github.com/damienfamed75/endorem/pkg/room"
	"github.com/damienfamed75/endorem/pkg/testing"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ Scene = &LevelOne{}
)

type LevelOne struct {
	mapData *dngn.Room
	rooms   []room.RoomSpec
	player  *player.Player
	ground  *physics.Space
	// ground      *resolv.Space
	doors       *resolv.Space
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
	l.ground = physics.NewSpace()
	l.doors = resolv.NewSpace()

	l.mapData, l.rooms = room.GenerateMap(1)
	mapScale := 34

	var spawnRoom room.RoomSpec
	for i := range l.rooms {
		if l.rooms[i].Size == -1 {
			spawnRoom = l.rooms[i]
		}
	}

	// +1 to Y so player doesn't shoot up the ceiling collider.
	x, y := (spawnRoom.X * mapScale), ((spawnRoom.Y + 1) * mapScale)

	l.player = player.NewPlayer(x, y, func() {}, physics.NewSpace())
	l.camera = common.NewEndoCamera(l.player.Collision)

	fmt.Printf("player inventory before [%v]\n", l.player.Inventory)

	l.player.Inventory.AddItem(&testing.Item{})

	// Show that the player has gotten an item that does nothing.
	fmt.Printf("player inventory after [%v]\n", l.player.Inventory)

	l.mapData.Select().By(func(x, y int) bool {
		switch l.mapData.Get(x, y) {
		case room.Wall: // Wall
			l.ground.Add(
				// testing.NewSolidPlane(
				// 	int32(x*mapScale), int32(y*mapScale),
				// 	int32(mapScale), int32(mapScale),
				// 	r.Aqua,
				// ),
				testing.NewPlane(
					int32(x*mapScale), int32(y*mapScale),
					int32(mapScale), int32(mapScale),
					r.Aqua,
				),
			)
		case room.Door: // Door
			l.doors.Add(
			// testing.NewDoor(
			// 	int32((x*mapScale)+(mapScale/2)), int32(y*mapScale),
			// 	int32((x*mapScale)+(mapScale/2)), int32((y*mapScale)+mapScale),
			// 	r.Orange,
			// ),
			)
		case room.Hatch: // Hatches
			l.doors.Add(
			// testing.NewPlane(
			// 	int32(x*mapScale), int32(y*mapScale),
			// 	int32(mapScale), int32(mapScale),
			// 	r.Gold,
			// ),
			)
		case room.FloatingPlatform1: // Floating Platform 1
			l.ground.Add(
			// testing.NewSolidPlane(
			// 	int32(x*mapScale), int32(y*mapScale),
			// 	int32(mapScale), int32(mapScale),
			// 	r.GopherBlue,
			// ),
			)
		case room.FloatingPlatform2: // Floating Platform 2
			l.ground.Add(
			// testing.NewSolidPlane(
			// 	int32(x*mapScale), int32(y*mapScale),
			// 	int32(mapScale), int32(mapScale),
			// 	r.SkyBlue,
			// ),
			)
		}

		return false
	})

	l.player.Body.AddGround(*l.ground...)
	// l.player.Body.AddGround(l.ground)
	// l.world.Add(l.player)
}

func (l *LevelOne) Update(dt float32) {
	l.camera.Update(l.player.Update(dt))
}

func (l *LevelOne) Draw() {
	// r.BeginMode2D(l.overviewCam)
	r.BeginMode2D(l.camera.Camera2D)
	r.ClearBackground(r.Black)

	for i := range *l.ground {
		(*l.ground)[i].(Drawer).Draw()
	}
	for i := range *l.doors {
		(*l.doors)[i].(Drawer).Draw()
	}

	l.player.Draw()

	r.EndMode2D()
}

func (l *LevelOne) Unload() {

}

func (l *LevelOne) String() string {
	return "level one"
}
