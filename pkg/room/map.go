package room

import (
	"github.com/SolarLune/dngn"
)

// RoomSpec is enough information to fill and locate rooms around the map.
type RoomSpec struct {
	X, Y, X2, Y2 int
	Size         int
	Selection    dngn.Selection
}

// GenerateMap creates a new completely random map.
func GenerateMap(bossroom int) (*dngn.Room, []RoomSpec) {
	sceneMap := dngn.NewRoom(60, 30)
	// sceneMap := dngn.NewRoom(40, 20)

	// treasure room (x,y,x2,y2)
	// boss room (x,y,x2,y2)
	// player spawn (x,y)
	rooms := newMap(sceneMap)

	switch bossroom {
	case 1:
		return InsertBossOneRoom(sceneMap, rooms)
	default:
		return sceneMap, rooms
	}
}
