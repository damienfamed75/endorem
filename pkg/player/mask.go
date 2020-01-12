package player

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Mask struct {
	Sprite      r.Texture2D
	movePattern string
	current     r.Vector2
	target      r.Vector2
	Facing      common.Direction

	Hitbox *resolv.Rectangle
	state  common.State

	*resolv.Space
}

func setupMask() *Mask {
	return &Mask{
		Sprite: r.LoadTexture("assets/mask.png"),
		Facing: common.Right,

		state: common.StateIdle,
		Space: resolv.NewSpace(),
	}
}

func NewMask() *Mask {
	m := setupMask()

	m.SetData(m)
	//TODO HITBOX

	return m
}
func (m *Mask) setMovePattern(moveType string) {
	m.movePattern = "figureEight"
}
func (m *Mask) checkDirection(diff r.Vector2, pFacing common.Direction) {

	var newTarget r.Vector2
	if pFacing == common.Right {
		newTarget.X = diff.X - 8
		newTarget.Y = diff.Y - 16

	} else if pFacing == common.Left {
		newTarget.X = diff.X + 16
		newTarget.Y = diff.Y - 16
	}
	m.target = newTarget
}
func (m *Mask) Update() {
	m.current = m.current.Lerp(m.target, 0.1)
}
func (m *Mask) Draw() {
	r.DrawTexture(m.Sprite, int(m.current.X), int(m.current.Y), r.White)
	//log.Print("m draw")
}
