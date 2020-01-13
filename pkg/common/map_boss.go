package common

import (
	"fmt"

	"github.com/SolarLune/dngn"
	r "github.com/lachee/raylib-goplus/raylib"
)

type readyDirection struct {
	name string
	x, y int
}

func IsMapReadyForPlacedThings(sceneMap *dngn.Room) ([]readyDirection, bool) {
	// required := []string{"up"}
	var directions []readyDirection

	selection := sceneMap.Select()

	selection.ByRune('#').By(func(x, y int) bool {
		// if direction != "" {
		// 	return false
		// }

		switch x {
		// case 0: // Build left?
		// if sceneMap.Get(x+1, y) == ' ' {
		// 	direction = "left"
		// }
		case sceneMap.Width - 1:
			fallthrough
		case sceneMap.Width: // Build right?
			if sceneMap.Get(x-1, y) == ' ' {
				directions = append(directions, readyDirection{"right", x, y})

				rightWall := selection.ByArea(x, 0, 1, sceneMap.Height)
				selection.RemoveSelection(rightWall)
			}
		}

		switch y {
		case 0: // Build up?
			if sceneMap.Get(x, y+1) == ' ' {
				directions = append(directions, readyDirection{"up", x, y})

				upperWall := selection.ByArea(0, 0, sceneMap.Width, 1)
				selection.RemoveSelection(upperWall)
			}
		case sceneMap.Height - 1: // Build down?
			fallthrough
		case sceneMap.Height: // Build down?
			if sceneMap.Get(x, y-1) == ' ' {
				directions = append(directions, readyDirection{"down", x, y})

				bottWall := selection.ByArea(0, y, sceneMap.Width, 1)
				selection.RemoveSelection(bottWall)
			}
		}

		return false
	})

	for i := range directions {
		// If the required is met then return if there are more than 1 direction available.
		if directions[i].name == "up" {
			return directions, len(directions) > 1
		}
	}

	return directions, false
}

func InsertBossOneRoom(sceneMap *dngn.Room, rooms []RoomSpec) (*dngn.Room, []RoomSpec) {
	var (
		// direction string
		xO, yO int

		roomWidth, roomHeight = 20, 10

		upDir readyDirection
	)

	directions, ok := IsMapReadyForPlacedThings(sceneMap)
	if !ok {
		return restart()
	}

	var direction string
	for i := range directions {
		if directions[i].name != "up" {
			direction = directions[i].name
			xO, yO = directions[i].x, directions[i].y
		} else {
			upDir = directions[i]
		}
	}

	var newScene *dngn.Room

	if direction == "" {
		return restart()
	}

	newScene = dngn.NewRoom(sceneMap.Width+roomWidth+1, sceneMap.Height+roomHeight+1)

	oldLvlSelection := newScene.Select().ByArea(0, 0, sceneMap.Width, sceneMap.Height)
	newScene.Select().RemoveSelection(oldLvlSelection).Fill('#')

	switch direction {
	// case "left":
	// newScene.CopyFrom(sceneMap, 50, 0)
	// newScene.Select().ByArea(0, 0, 50, sceneMap.Height).Fill('*')
	// fallthrough
	case "right":
		newScene.CopyFrom(sceneMap, 0, 0)
		newScene.Set(xO, yO, '}')
		newScene.Select().ByArea(sceneMap.Width, yO, roomWidth, roomHeight).Fill(' ')
	// case "up":
	// 	newScene.CopyFrom(sceneMap, 0, 100)
	case "down":
		newScene.CopyFrom(sceneMap, 0, 0)
		newScene.Set(xO, yO, '}')
		newScene.Select().ByArea(xO, sceneMap.Height, roomWidth, roomHeight).Fill(' ')
	}

	// Ensure that the map doesn't have any missing floors, ceilings, or walls.
	drawMapBorders(newScene)

	fmt.Printf("MAP: Map Successfully Generated!\n")

	debugMap(newScene.DataToString())

	sceneMap = newScene

	_ = upDir
	newScene = dngn.NewRoom(sceneMap.Width, sceneMap.Height+10+1)
	newScene.CopyFrom(sceneMap, 0, 10)

	newScene.Select().ByArea(0, 0, sceneMap.Width, 10).Fill('#')
	newScene.DrawLine(upDir.x, upDir.y+12, upDir.x, 2, ' ', 2, false)
	fmt.Println(newScene.DataToString())

	newScene.Select().ByRune('}').By(func(x, y int) bool {
		xO, yO = x, y
		return true
	})

	// Insert the boss room and player spawn areas.
	rooms = append(rooms, RoomSpec{
		X:    upDir.x,
		Y:    upDir.y,
		Size: PlayerSpawn, // -1 indicates player spawn area.
	}, RoomSpec{
		X:    xO,
		Y:    yO,
		Size: BossSpawn, // -2 indicates boss spawn area.
	})

	// Ensure that the player can get from the spawn to the boss room.
	if !validateMap(
		newScene,
		r.NewVector2(float32(upDir.x), float32(upDir.y+1)),
		r.NewVector2(float32(xO), float32(yO)),
	) {
		return restart()
	}

	return newScene, rooms
}

func restart() (*dngn.Room, []RoomSpec) {
	newSceneMap := dngn.NewRoom(60, 30)
	newRooms := newMap(newSceneMap)
	return InsertBossOneRoom(newSceneMap, newRooms)
}
