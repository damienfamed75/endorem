package player

import (
	"time"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/endorem/pkg/common"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Mask handles the companion character that follows the player
type Mask struct {
	Sprite      r.Texture2D
	movePattern string
	current     r.Vector2
	target      r.Vector2
	Facing      common.Direction

	shotCooldown time.Time
	shotTimer    time.Duration
	shotRange    float32
	shotSpeed    int32

	Hitbox *resolv.Space
	state  common.State

	*resolv.Space
}

func setupMask() *Mask {
	return &Mask{
		Sprite: r.LoadTexture("assets/mask.png"),
		Facing: common.Right,

		shotCooldown: time.Now(),
		shotTimer:    time.Duration(500),
		shotRange:    500,
		shotSpeed:    10,
		Hitbox:       resolv.NewSpace(),
		state:        common.StateIdle,
		Space:        resolv.NewSpace(),
	}
}

// NewMask returns a configured Mask entity
func NewMask() *Mask {
	m := setupMask()

	m.SetData(m)

	return m
}

// setMovePattern will change how the mask will move around the player
func (m *Mask) setMovePattern(moveType string) {
	m.movePattern = "figureEight"
}

// checkDirecton will change the target
//TODO depending on movePattern
func (m *Mask) checkDirection(diff r.Vector2, pFacing common.Direction) {

	var newTarget r.Vector2
	if pFacing == common.Right {
		m.Facing = common.Right

		newTarget.X = diff.X - 8
		newTarget.Y = diff.Y - 16

	} else if pFacing == common.Left {
		m.Facing = common.Left

		newTarget.X = diff.X + 16
		newTarget.Y = diff.Y - 16
	}
	m.target = newTarget
}

// Update Mask
func (m *Mask) Update() {
	m.current = m.current.Lerp(m.target, 0.1)

	m.shoot()

	// Move shot bullets
	for _, shape := range *m.Hitbox {
		x, y := shape.GetXY()
		if shape.HasTags("left") {
			shape.SetXY(x-m.shotSpeed, y)
		} else if shape.HasTags("right") {
			shape.SetXY(x+m.shotSpeed, y)
		}

	}
}

func (m *Mask) shoot() {
	// Shoot when key is pressed and not on CD
	if r.IsKeyDown(r.KeyC) {
		if time.Since(m.shotCooldown) >= time.Millisecond*m.shotTimer {
			m.shotCooldown = time.Now()
			bullet := resolv.NewRectangle(
				int32(m.current.X), int32(m.current.Y), 5, 5)
			if m.Facing == common.Left {
				bullet.AddTags("left")
			} else {
				bullet.AddTags("right")
			}
			m.Hitbox.Add(bullet)
		}
	}

	// Remove bullet when out of range
	for _, bullet := range *m.Hitbox {
		x, _ := bullet.GetXY()

		if x > int32((m.current.X+m.shotRange)) || x < int32((m.current.X-m.shotRange)) {
			(*m.Hitbox).Remove(bullet)
		}
	}
}

// Draw Mask at new frame
func (m *Mask) Draw() {
	r.DrawTexture(m.Sprite, int(m.current.X), int(m.current.Y), r.White)

	m.debugDraw()
}

func (m *Mask) debugDraw() {
	for _, shape := range *m.Hitbox {
		x, y := shape.GetXY()
		r.DrawRectangleLines(
			int(x), int(y),
			int(5), int(5),
			r.Blue,
		)
	}
}
