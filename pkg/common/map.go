package common

import (
	"fmt"
	"strings"

	"github.com/SolarLune/dngn"
)

// GenerateMap creates a new completely random map.
func GenerateMap() (*dngn.Room, []RoomSpec) {
	sceneMap := dngn.NewRoom(40, 20)

	// treasure room (x,y,x2,y2)
	// boss room (x,y,x2,y2)
	// player spawn (x,y)
	rooms := newMap(sceneMap)

	return sceneMap, rooms
}

// RoomSpec is enough information to fill and locate rooms around the map.
type RoomSpec struct {
	x, y, x2, y2 int
	size         int
	selection    dngn.Selection
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

				room.Fill(':')

				// Remove the found room from the selection so then it doesn't
				// get scanned twice.
				selection.RemoveSelection(room)

				// Make sure that the room is valid size.
				if len(room.Cells) > 0 {
					foundRooms = append(foundRooms, RoomSpec{
						x: x, y: y, x2: xa, y2: ya,
						size:      len(room.Cells) * len(room.Cells[0]),
						selection: room,
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

	// sceneMap.GenerateDrunkWalk(' ', 0.5)
	// rooms := sceneMap.GenerateRandomRooms(' ', 80, 4, 4, 5, 5, true)
	sceneMap.GenerateRandomRooms(' ', 24, 4, 4, 5, 5, true)
	fmt.Printf("MAP: Map Generation 1 Complete!\n")
	for i := 0; i < 10; i++ {
		// sceneMap.Select().Degrade(' ')
		sceneMap.Select().Degrade('#')
	}

	for i := 0; i < 3; i++ {
		sceneMap.Select().ByRune(' ').ByNeighbor('#', 3, false).Fill('#')
	}

	sceneMap.GenerateBSP('#', '-', 80)
	fmt.Printf("MAP: Map Generation 2 Complete!\n")

	rooms := findRooms(sceneMap)

	fmt.Printf("MAP: Treasure Room Marked!\n")

	sceneMap.Select().ByRune('-').By(func(x, y int) bool {
		// Ceiling or floor doors.
		if sceneMap.Get(x, y+1) == ' ' {
			sceneMap.Set(x, y, '^')

			if sceneMap.Get(x, y+2) == ' ' && sceneMap.Get(x, y+3) == ' ' {
				sceneMap.Set(x, y+2, '=')

				if sceneMap.Get(x, y+4) == ' ' {
					sceneMap.Set(x, y+4, '=')
				}
			}

			return true
		}

		// Wall doors with no ledges
		if sceneMap.Get(x+1, y+1) == ' ' || sceneMap.Get(x-1, y+1) == ' ' {
			// Door ledges.
			sceneMap.Set(x+1, y+1, '#')
			sceneMap.Set(x-1, y+1, '#')

			// Floating ledges to get to the doore ledges.
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

		for sceneMap.Get(x, ya) == ' ' {
			sceneMap.Set(x, ya, '#')
			ya++
		}

		ya = y - 1

		for sceneMap.Get(x, ya) == ' ' {
			sceneMap.Set(x, ya, '#')
			ya--
		}

		return false
	})

	fmt.Printf("MAP: Exceptions Marked!\n")

	// Draw borders around the map.
	sceneMap.DrawLine(0, 0, 0, len(sceneMap.Data), '#', 1, false)
	sceneMap.DrawLine(0, 0, len(sceneMap.Data[0]), 0, '#', 1, false)
	sceneMap.DrawLine(len(sceneMap.Data[0]), 0, len(sceneMap.Data[0]), len(sceneMap.Data), '#', 1, false)
	sceneMap.DrawLine(0, len(sceneMap.Data), len(sceneMap.Data[0]), len(sceneMap.Data), '#', 1, false)

	fmt.Printf("MAP: Map Successfully Generated!\n")

	debugMap(sceneMap.DataToString())

	return rooms
}

func debugMap(data string) {
	data = strings.Replace(data, "X", "\033[1;34mX\033[0m", -1)
	data = strings.Replace(data, "Y", "\033[1;31mY\033[0m", -1)
	data = strings.Replace(data, ":", "\033[1;32m:\033[0m", -1)

	fmt.Println(data)
}
