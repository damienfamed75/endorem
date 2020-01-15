package enemy

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Spores struct {
	Sprite      r.Texture2D
	sporeHeight int
	sporeWidth  int

	Speed int32

	rowStart  int32
	rowHeight int32 // Spawn height of spores
	rowAmt    int   //amt of spores in a row
	rowApart  int32 // distance each spore is from another

	*resolv.Space
}

func setupSpores() *Spores {
	return &Spores{
		Sprite:   r.LoadTexture("assets/playerDuck.png"),
		Speed:    -3,
		rowAmt:   5,
		rowApart: 80,
		Space:    resolv.NewSpace(),
	}
}

func NewSpores(bossX, bossY int32) *Spores {
	s := setupSpores()

	s.rowStart = bossX
	s.rowHeight = bossY - 200 // rows start 200 pixels above boss

	return s
}

// CreateRow will create a new row of spores
// and place them into the spores space
func (s *Spores) CreateRow() {

	// bases spore displacement on the boss position and time
	// how far each spore is (X-axis) in a row
	difference := rand.Int31n(80) + 50
	area := s.rowStart - difference

	// For the amount of spores, create a hitbox based on the
	// determined space apart
	for i := 0; i <= s.rowAmt; i++ {

		spore := resolv.NewRectangle(
			area, s.rowHeight, 5, 5)

		// Update the x-axis for the next spore
		area -= s.rowApart

		s.Space.Add(spore)
	}
}

// Update the position of each spore in space
func (s *Spores) Update() {
	// Updates the position of each indiviudal spore in Space
	for _, shape := range *s.Space {
		x, y := shape.GetXY()
		shape.SetXY(x, y-s.Speed)
	}

	//TODO remove spores once they reach the ground, need to add ground to spores
}

func (s *Spores) Draw() {

	// To check moving hitboxes of spores
	s.debugDraw()
}

func (s *Spores) debugDraw() {
	for _, shape := range *s.Space {
		x, y := shape.GetXY()
		r.DrawTexture(
			s.Sprite,
			int(x), int(y),
			r.White,
		)
	}
}
