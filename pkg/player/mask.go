package player

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Mask struct {
	X, Y   int
	Sprite r.Texture2D
	Hitbox *resolv.Rectangle
	Facing common.Direction

	state common.State

	*resolv.Space
}

func setupMask() *Mask {
	return &Mask{
		Sprite: r.LoadTexture("assets/mask.png"),
		Facing: common.Right,
		state:  common.StateIdle,
		Space:  resolv.NewSpace(),
	}
}

func NewMask() *Mask {
	m := setupMask()

	m.SetData(m)
	//TODO HITBOX

	return m
}
func (m *Mask) setMovePattern(moveType string) {

}
func (m *Mask) followPlayer(x, y int32) {
	m.X = int(x - 8)
	m.Y = int(y - 16)
}

func (m *Mask) Draw() {
	r.DrawTexture(m.Sprite, int(m.X), int(m.Y), r.White)
	//log.Print("m draw")
}
