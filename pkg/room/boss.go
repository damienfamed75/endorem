package room

import (
	"fmt"

	"github.com/SolarLune/dngn"
	r "github.com/lachee/raylib-goplus/raylib"
)

// readyDirection stores the name of the direction that is ready for a boss room
// or player spawn along with an x, y position for them to be placed.
type readyDirection struct {
	name string
	x, y int
}

// isRoomShapedRight checks to see if there is a valid spot to place a boss room
// and a player spawn in the generated room. If there isn't then the function
// returns false, which signifies that this room should be regenerated.
func isRoomShapedRight(sceneMap *dngn.Room) ([]readyDirection, bool) {
	var directions []readyDirection

	selection := sceneMap.Select()

	selection.ByRune(Wall).By(func(x, y int) bool {
		switch x {
		// case 0: // Build left?
		// if sceneMap.Get(x+1, y) == Air {
		// 	direction = "left"
		// }
		case sceneMap.Width - 1:
			fallthrough
		case sceneMap.Width: // Build right?
			if sceneMap.Get(x-1, y) == Air {
				directions = append(directions, readyDirection{"right", x, y})

				rightWall := selection.ByArea(x, 0, 1, sceneMap.Height)
				selection.RemoveSelection(rightWall)
			}
		}

		switch y {
		case 0: // Build up?
			if sceneMap.Get(x, y+1) == Air {
				directions = append(directions, readyDirection{"up", x, y})

				upperWall := selection.ByArea(0, 0, sceneMap.Width, 1)
				selection.RemoveSelection(upperWall)
			}
		case sceneMap.Height - 1: // Build down?
			fallthrough
		case sceneMap.Height: // Build down?
			if sceneMap.Get(x, y-1) == Air {
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

	directions, ok := isRoomShapedRight(sceneMap)
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
	newScene.Select().RemoveSelection(oldLvlSelection).Fill(Wall)

	switch direction {
	// case "left":
	// newScene.CopyFrom(sceneMap, 50, 0)
	// newScene.Select().ByArea(0, 0, 50, sceneMap.Height).Fill('*')
	// fallthrough
	case "right":
		newScene.CopyFrom(sceneMap, 0, 0)
		newScene.Set(xO, yO, BossDoor)
		newScene.Select().ByArea(sceneMap.Width, yO, roomWidth, roomHeight).Fill(NoEnemies)
	// case "up":
	// 	newScene.CopyFrom(sceneMap, 0, 100)
	case "down":
		newScene.CopyFrom(sceneMap, 0, 0)
		newScene.Set(xO, yO, BossDoor)
		newScene.Select().ByArea(xO, sceneMap.Height, roomWidth, roomHeight).Fill(NoEnemies)
	}

	// Ensure that the map doesn't have any missing floors, ceilings, or walls.
	drawMapBorders(newScene)

	fmt.Printf("MAP: Map Successfully Generated!\n")

	debugMap(newScene.DataToString())

	sceneMap = newScene

	// Add player spawn area.
	newScene = dngn.NewRoom(sceneMap.Width, sceneMap.Height+10+1)
	newScene.CopyFrom(sceneMap, 0, 10)

	newScene.Select().ByArea(0, 0, sceneMap.Width, 10).Fill(Wall)
	newScene.DrawLine(upDir.x, upDir.y+11, upDir.x, 2, NoEnemies, 2, false)
	newScene.DrawLine(upDir.x-1, upDir.y+13, upDir.x+1, upDir.y+13, Wall, 1, false)
	fmt.Println(newScene.DataToString())

	newScene.Select().ByRune(BossDoor).By(func(x, y int) bool {
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

	sprinkleInEnemies(newScene)

	fmt.Println(newScene.DataToString())

	return newScene, rooms
}
