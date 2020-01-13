package room

import (
	"fmt"

	"github.com/SolarLune/dngn"
	"github.com/damienfamed75/endorem/pkg/common"
)

func sprinkleInEnemies(sceneMap *dngn.Room) {
	selection := sceneMap.Select().ByRune(Wall).ByNeighbor(Air, 1, false).
		ByPercentage(common.GlobalConfig.Enemy.SpawnDensity)

	spawnAreas := sceneMap.Select().ByRune(NoEnemies).Expand(2, true).Expand(2, true)

	selection.RemoveSelection(spawnAreas).By(func(x, y int) bool {
		if sceneMap.Get(x, y-1) == Air {
			sceneMap.Set(x, y-1, Enemy)
		}

		return false
	})

	fmt.Println(sceneMap.DataToString())
}
