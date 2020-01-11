package scene

import (
	"github.com/SolarLune/dngn"
	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ Scene = &LevelOne{}
)

type LevelOne struct {
	mapData *dngn.Room
	ground  *resolv.Space
	camera  *r.Camera2D
}

func (l *LevelOne) Preload() {
	l.camera = &r.Camera2D{
		Offset: r.NewVector2(
			0,
			0,
		),
		Rotation: 0,
		Zoom:     1,
	}

}

func (l *LevelOne) Update(dt float32) {

}

func (l *LevelOne) Draw() {

}

func (l *LevelOne) Unload() {

}

func (l *LevelOne) String() string {
	return "level one"
}
