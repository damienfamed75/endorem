package room

import (
	"fmt"

	"github.com/SolarLune/dngn"
)

// newMap is the logic behind generating the rooms, doors, hatches, and platforms
// within the map. Everytime this functionis called there is a new seed used
// based on the unix time when clling this.
func newMap(sceneMap *dngn.Room) []RoomSpec {
	sceneMap.Select().Fill(Wall)

	// IMPORTANT -----------------------------------------------------------------
	sceneMap.GenerateDrunkWalk(Air, 0.8)
	// sceneMap.GenerateRandomRooms(Air, 80, 4, 4, 6, 6, true)
	sceneMap.GenerateDrunkWalk(Air, 0.2)
	sceneMap.GenerateDrunkWalk(Air, 0.2)
	// sceneMap.GenerateRandomRooms(Air, 24, 4, 4, 5, 5, true)
	fmt.Printf("MAP: Map Generation 1 Complete!\n")
	for i := 0; i < 100; i++ {
		sceneMap.Select().Degrade(Air)
		// sceneMap.Select().Degrade(Wall)
	}

	for i := 0; i < 3; i++ {
		sceneMap.Select().ByRune(Air).ByNeighbor(Wall, 3, false).Fill(Wall)
	}

	// IMPORTANT -----------------------------------------------------------------
	sceneMap.GenerateBSP(Wall, Door, 100)
	// sceneMap.GenerateBSP(Wall, Door, 80)
	fmt.Printf("MAP: Map Generation 2 Complete!\n")

	rooms := findRooms(sceneMap)

	fmt.Printf("MAP: Treasure Room Marked!\n")

	sceneMap.Select().ByRune(Door).By(func(x, y int) bool {
		// Ceiling or floor doors.
		if sceneMap.Get(x, y+1) == Air || sceneMap.Get(x, y+1) == ':' {
			sceneMap.Set(x, y, Hatch)

			offset := 2
			for (sceneMap.Get(x, y+offset) == Air || sceneMap.Get(x, y+offset) == ':') && (sceneMap.Get(x, y+offset+1) == Air || sceneMap.Get(x, y+offset+1) == ':') {
				sceneMap.Set(x, y+offset, FloatingPlatform2)

				offset += 2
			}

			return true
		}

		// Wall doors with no ledges
		if sceneMap.Get(x+1, y+1) == Air || sceneMap.Get(x-1, y+1) == Air {
			// Door ledges.
			if sceneMap.Get(x+1, y+2) != Hatch || sceneMap.Get(x+1, y+2) != Door {
				sceneMap.Set(x+1, y+1, Wall)
			}
			if sceneMap.Get(x-1, y+2) != Hatch || sceneMap.Get(x-1, y+2) != Door {
				sceneMap.Set(x-1, y+1, Wall)
			}

			// Remove walls too close to doors.
			sceneMap.Set(x-2, y, Air)
			sceneMap.Set(x+2, y, Air)

			// Floating ledges to get to the door ledges.
			if sceneMap.Get(x-2, y+2) == Air && sceneMap.Get(x-2, y+3) != Hatch {
				sceneMap.Set(x-2, y+2, FloatingPlatform1)
			}
			if sceneMap.Get(x+2, y+2) == Air && sceneMap.Get(x+2, y+3) != Hatch {
				sceneMap.Set(x-2, y+2, FloatingPlatform1)
			}
		}

		return false
	})

	// With the rest of the left doors make sure that they have a top and bottom.
	sceneMap.Select().ByRune(Door).By(func(x, y int) bool {
		ya := y + 1

		sceneMap.Set(x-1, y, Air)
		sceneMap.Set(x+1, y, Air)

		for sceneMap.Get(x, ya) == Air || sceneMap.Get(x, ya) == ':' {
			sceneMap.Set(x, ya, Wall)
			ya++
		}

		ya = y - 1

		for sceneMap.Get(x, ya) == Air || sceneMap.Get(x, ya) == ':' {
			sceneMap.Set(x, ya, Wall)
			ya--
		}

		return false
	})

	sceneMap.Select().ByRune(Hatch).By(func(x, y int) bool {
		sceneMap.Set(x, y+1, Air)
		sceneMap.Set(x, y-1, Air)

		return false
	})

	// Remove 50% of the doors.
	sceneMap.Select().ByRune(Door).ByPercentage(.5).By(func(x, y int) bool {
		sceneMap.Set(x, y, Air)
		return false
	})

	fmt.Printf("MAP: Exceptions Marked!\n")

	drawMapBorders(sceneMap)

	return rooms
}

func drawMapBorders(sceneMap *dngn.Room) {
	// Draw borders around the map.
	sceneMap.DrawLine(0, 0, 0, len(sceneMap.Data), Wall, 1, false)
	sceneMap.DrawLine(0, 0, len(sceneMap.Data[0]), 0, Wall, 1, false)
	sceneMap.DrawLine(len(sceneMap.Data[0]), 0, len(sceneMap.Data[0]), len(sceneMap.Data), Wall, 1, false)
	sceneMap.DrawLine(0, len(sceneMap.Data), len(sceneMap.Data[0]), len(sceneMap.Data), Wall, 1, false)
}
