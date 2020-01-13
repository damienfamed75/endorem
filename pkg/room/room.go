package room

// // findRooms scans through the map for rooms that meet a certain required size
// // and then returns the size specifications of each individual one along with a
// // selection of it on the map.
// func findRooms(sceneMap *dngn.Room) []RoomSpec {
// 	var (
// 		selection  = sceneMap.Select()
// 		foundRooms []RoomSpec
// 	)

// 	// Loop through the selection to find rooms.
// 	selection.By(func(x, y int) bool {
// 		if sceneMap.Get(x, y) == Wall {
// 			xa, ya := x+1, y+1

// 			var diagRoom int // The diagonal size of the measured area.

// 			// Pure diagonal traversal.
// 			for sceneMap.Get(xa, ya) != Wall && sceneMap.Get(xa, ya) != Door {
// 				xa++
// 				ya++
// 				diagRoom++

// 				// Break if the diagonal count is over the total map area.
// 				if diagRoom > sceneMap.Area() {
// 					return false
// 				}
// 			}

// 			// Big enough room.
// 			if diagRoom >= 3 {
// 				// Vertical traversal.
// 				// This traverses vertically in the room to get the taller ones.
// 				for sceneMap.Get(xa-1, ya) != Wall && sceneMap.Get(xa-1, ya) != Door {
// 					ya++

// 					// Break infinite loops.
// 					if ya > sceneMap.Area() {
// 						return false
// 					}
// 				}

// 				fmt.Printf("Room found\n")
// 				room := sceneMap.Select().ByArea(x, y, xa-x+1, ya-y+1).ByRune(Air)

// 				// room.Fill(':')

// 				// Remove the found room from the selection so then it doesn't
// 				// get scanned twice.
// 				selection.RemoveSelection(room)

// 				// Make sure that the room is valid size.
// 				if len(room.Cells) > 0 {
// 					foundRooms = append(foundRooms, RoomSpec{
// 						X: x, Y: y, X2: xa, Y2: ya,
// 						Size:      len(room.Cells) * len(room.Cells[0]),
// 						Selection: room,
// 					})
// 				}
// 			}

// 		}

// 		return false
// 	})

// 	return foundRooms
// }
