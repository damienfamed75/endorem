package common

import (
	"github.com/SolarLune/dngn"
	"github.com/SolarLune/paths"
	r "github.com/lachee/raylib-goplus/raylib"
)

const (
	PlayerSpawn = iota*-1 - 1
	BossSpawn
)

func validateMap(sceneMap *dngn.Room, start, end r.Vector2) bool {
	grid := paths.NewGridFromRuneArrays(sceneMap.Data)

	// Set the walls as non-walkable areas.
	for _, cell := range grid.GetCellsByRune('#') {
		cell.Walkable = false
	}

	for _, cell := range grid.GetCellsByRune(' ') {
		cell.Cost = 1
	}

	path := grid.GetPath(
		grid.Get(int(start.X), int(start.Y)),
		grid.Get(int(end.X), int(end.Y)),
		false,
	)

	return path.Valid()
}
