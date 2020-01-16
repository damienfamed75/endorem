package enemy

import (
	"log"
	"math"
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
	"github.com/damienfamed75/aseprite"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Spores struct {
	Ase    *aseprite.File
	Sprite r.Texture2D

	sporeHeight int
	sporeWidth  int
	sporeMoveX  float64
	Speed       int32

	rowStart  int32
	rowHeight int32 // Spawn height of spores
	rowAmt    int   //amt of spores in a row
	rowApart  int32 // distance each spore is from another

	*resolv.Space
}

func setupSpores() *Spores {
	return &Spores{
		Sprite:      r.LoadTexture("assets/spore.png"),
		sporeHeight: 10,
		sporeWidth:  10,
		Speed:       -3,
		rowAmt:      5,
		rowApart:    80,
		Space:       resolv.NewSpace(),
	}
}

func NewSpores(bossX, bossY int32) *Spores {
	s := setupSpores()
	var err error

	s.Ase, err = aseprite.Open("assets/spore.json")
	if err != nil {
		log.Fatal(err)
	}

	s.rowStart = bossX
	s.rowHeight = bossY - 300 // rows start 300 pixels above boss

	s.Ase.Play("fire")

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
		rowSporeY := rand.Int31n(20) + s.rowHeight
		spore := resolv.NewRectangle(
			area, rowSporeY, 5, 5)

		// Update the x-axis for the next spore
		area -= s.rowApart

		s.Space.Add(spore)
	}
}

// Update the position of each spore in space
func (s *Spores) Update(dt float32) {
	s.Ase.Update(dt)

	s.sporeMoveX += 0.5
	// Updates the position of each indiviudal spore in Space
	for _, shape := range *s.Space {
		x, y := shape.GetXY()
		displacementX := (2 * int32(math.Sin(s.sporeMoveX*(math.Pi))))

		shape.SetXY(x+displacementX, y-s.Speed)
	}

	//TODO remove spores once they reach the ground, need to add ground to spores
}

func (s *Spores) Draw() {
	// srcX, srcY resemble the X and Y pixels where the active sprite is.
	srcX, srcY := s.Ase.FrameBoundaries().X, s.Ase.FrameBoundaries().Y
	w, h := s.Ase.FrameBoundaries().Width, s.Ase.FrameBoundaries().Height

	src := r.NewRectangle(float32(srcX), float32(srcY), float32(w), float32(h))
	for _, shape := range *s.Space {
		x, y := shape.GetXY()
		dest := r.NewRectangle(float32(x), float32(y), float32(w), float32(h))

		r.DrawTexturePro(
			s.Sprite, src, dest, r.NewVector2(0, 0), 0, r.White,
		)
	}
	// To check moving hitboxes of spores
	s.debugDraw()
}

func (s *Spores) debugDraw() {
	for _, shape := range *s.Space {
		x, y := shape.GetXY()
		r.DrawRectangleLines(
			int(x), int(y),
			int(s.sporeWidth), int(s.sporeHeight),
			r.Red,
		)
	}
}
