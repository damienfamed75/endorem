package common

import (
	"fmt"

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

	// return sceneMap, rooms
}

func findRooms(sceneMap *dngn.Room) []RoomSpec {
	var (
		selection  = sceneMap.Select()
		foundRooms []RoomSpec
	)

	// Loop through the selection to find rooms.
	selection.By(func(x, y int) bool {
		if sceneMap.Get(x, y) == '#' {
			xa, ya := x+1, y+1

			var diagRoom int // The diagonal size of the measured area.

			// Pure diagonal traversal.
			for sceneMap.Get(xa, ya) != '#' && sceneMap.Get(xa, ya) != '-' {
				xa++
				ya++
				diagRoom++

				// Break if the diagonal count is over the total map area.
				if diagRoom > sceneMap.Area() {
					return false
				}
			}

			// Big enough room.
			if diagRoom >= 3 {
				// Vertical traversal.
				// This traverses vertically in the room to get the taller ones.
				for sceneMap.Get(xa-1, ya) != '#' && sceneMap.Get(xa-1, ya) != '-' {
					ya++

					// Break infinite loops.
					if ya > sceneMap.Area() {
						return false
					}
				}

				fmt.Printf("Room found\n")
				room := sceneMap.Select().ByArea(x, y, xa-x+1, ya-y+1).ByRune(' ')

				// room.Fill(':')

				// Remove the found room from the selection so then it doesn't
				// get scanned twice.
				selection.RemoveSelection(room)

				// Make sure that the room is valid size.
				if len(room.Cells) > 0 {
					foundRooms = append(foundRooms, RoomSpec{
						X: x, Y: y, X2: xa, Y2: ya,
						Size:      len(room.Cells) * len(room.Cells[0]),
						Selection: room,
					})
				}
			}

		}

		return false
	})

	return foundRooms
}

func newMap(sceneMap *dngn.Room) []RoomSpec {
	sceneMap.Select().Fill('#')

	// IMPORTANT -----------------------------------------------------------------
	// sceneMap.GenerateDrunkWalk(' ', 0.8)
	// sceneMap.GenerateRandomRooms(' ', 80, 4, 4, 6, 6, true)
	sceneMap.GenerateDrunkWalk(' ', 0.2)
	sceneMap.GenerateDrunkWalk(' ', 0.2)
	// sceneMap.GenerateRandomRooms(' ', 24, 4, 4, 5, 5, true)
	fmt.Printf("MAP: Map Generation 1 Complete!\n")
	for i := 0; i < 100; i++ {
		sceneMap.Select().Degrade(' ')
		// sceneMap.Select().Degrade('#')
	}

	for i := 0; i < 3; i++ {
		sceneMap.Select().ByRune(' ').ByNeighbor('#', 3, false).Fill('#')
	}

	// IMPORTANT -----------------------------------------------------------------
	sceneMap.GenerateBSP('#', '-', 100)
	// sceneMap.GenerateBSP('#', '-', 80)
	fmt.Printf("MAP: Map Generation 2 Complete!\n")

	rooms := findRooms(sceneMap)

	fmt.Printf("MAP: Treasure Room Marked!\n")

	sceneMap.Select().ByRune('-').By(func(x, y int) bool {
		// Ceiling or floor doors.
		if sceneMap.Get(x, y+1) == ' ' || sceneMap.Get(x, y+1) == ':' {
			sceneMap.Set(x, y, '^')

			offset := 2
			for (sceneMap.Get(x, y+offset) == ' ' || sceneMap.Get(x, y+offset) == ':') && (sceneMap.Get(x, y+offset+1) == ' ' || sceneMap.Get(x, y+offset+1) == ':') {
				sceneMap.Set(x, y+offset, '=')

				offset += 2
			}

			return true
		}

		// Wall doors with no ledges
		if sceneMap.Get(x+1, y+1) == ' ' || sceneMap.Get(x-1, y+1) == ' ' {
			// Door ledges.
			if sceneMap.Get(x+1, y+2) != '^' || sceneMap.Get(x+1, y+2) != '-' {
				sceneMap.Set(x+1, y+1, '#')
			}
			if sceneMap.Get(x-1, y+2) != '^' || sceneMap.Get(x-1, y+2) != '-' {
				sceneMap.Set(x-1, y+1, '#')
			}

			// Remove walls too close to doors.
			sceneMap.Set(x-2, y, ' ')
			sceneMap.Set(x+2, y, ' ')

			// Floating ledges to get to the door ledges.
			if sceneMap.Get(x-2, y+2) == ' ' && sceneMap.Get(x-2, y+3) != '^' {
				sceneMap.Set(x-2, y+2, '~')
			}
			if sceneMap.Get(x+2, y+2) == ' ' && sceneMap.Get(x+2, y+3) != '^' {
				sceneMap.Set(x-2, y+2, '~')
			}
		}

		return false
	})

	// With the rest of the left doors make sure that they have a top and bottom.
	sceneMap.Select().ByRune('-').By(func(x, y int) bool {
		ya := y + 1

		sceneMap.Set(x-1, y, ' ')
		sceneMap.Set(x+1, y, ' ')

		for sceneMap.Get(x, ya) == ' ' || sceneMap.Get(x, ya) == ':' {
			sceneMap.Set(x, ya, '#')
			ya++
		}

		ya = y - 1

		for sceneMap.Get(x, ya) == ' ' || sceneMap.Get(x, ya) == ':' {
			sceneMap.Set(x, ya, '#')
			ya--
		}

		return false
	})

	sceneMap.Select().ByRune('^').By(func(x, y int) bool {
		sceneMap.Set(x, y+1, ' ')
		sceneMap.Set(x, y-1, ' ')

		return false
	})

	fmt.Printf("MAP: Exceptions Marked!\n")

	drawMapBorders(sceneMap)

	return rooms
}

func drawMapBorders(sceneMap *dngn.Room) {
	// Draw borders around the map.
	sceneMap.DrawLine(0, 0, 0, len(sceneMap.Data), '#', 1, false)
	sceneMap.DrawLine(0, 0, len(sceneMap.Data[0]), 0, '#', 1, false)
	sceneMap.DrawLine(len(sceneMap.Data[0]), 0, len(sceneMap.Data[0]), len(sceneMap.Data), '#', 1, false)
	sceneMap.DrawLine(0, len(sceneMap.Data), len(sceneMap.Data[0]), len(sceneMap.Data), '#', 1, false)
}
