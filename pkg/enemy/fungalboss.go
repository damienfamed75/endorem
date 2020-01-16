package enemy

import (
	"log"
	"time"

	"github.com/damienfamed75/aseprite"
	"github.com/damienfamed75/endorem/pkg/common"
	"github.com/damienfamed75/endorem/pkg/physics"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.Entity = &FungalBoss{}
)

// FungalBoss is the set boss of level one, this enemy will be placed
// at the end of the first level
type FungalBoss struct {
	Ase    *aseprite.File
	Sprite r.Texture2D

	Health  int
	isDead  bool
	Hurtbox *resolv.Rectangle

	player *physics.Space
	ground *physics.Space

	maxSpeedX int32
	maxSpeedY int32

	aoeAse          *aseprite.File
	aoeSprite       r.Texture2D
	attackAoeBefore time.Time
	attackAoeTimer  time.Duration // Time before making another attack(milliseconds)
	HitboxAOE       *resolv.Rectangle
	attackDistance  int32 // How far aoe will go
	attackSpeedAOE  float32
	attackRangeAOE  int32 // Calculated X-axis that distance will go
	isAttackingAOE  bool

	sporeBefore time.Time //How often the boss can shoot spores(milliseconds)
	sporeTimer  time.Duration
	spores      *Spores

	*resolv.Space
}

func setupFungalBoss() *FungalBoss {
	return &FungalBoss{
		Sprite:    r.LoadTexture("assets/mushroom.png"),
		Health:    1 + common.GlobalConfig.Enemy.AddedHealth,
		maxSpeedX: 0,
		maxSpeedY: 0,

		aoeSprite:       r.LoadTexture("assets/smoke.png"),
		attackSpeedAOE:  -10,
		isAttackingAOE:  false,
		attackAoeBefore: time.Now(),
		attackAoeTimer:  time.Duration(5000),
		attackDistance:  200,
		sporeBefore:     time.Now(),
		sporeTimer:      time.Duration(1000), //REMOVE HARDCODED #
		Space:           resolv.NewSpace(),
	}
}

func NewFungalBoss(x, y int, world *physics.Space) *FungalBoss {
	f := setupFungalBoss()
	var err error

	f.Ase, err = aseprite.Open("assets/mushroom.json")
	if err != nil {
		log.Fatal(err)
	}

	f.aoeAse, err = aseprite.Open("assets/smoke.json")
	if err != nil {
		log.Fatal(err)
	}

	// Queues a default animation.
	f.Ase.Play("idle")

	// Store important spaces in the world
	f.player = world.FilterByTags(common.TagPlayer)
	f.ground = world.FilterByTags(common.TagGround)

	// For now, the hurtbox is a simple rectangle (the stem of the mushboy maybe)
	f.Hurtbox = resolv.NewRectangle(
		//TODO hardcoded width and height until sprite is made
		int32(x), int32(y), 128, 128,
	)
	f.Hurtbox.AddTags(TagHurtbox)

	f.HitboxAOE = &resolv.Rectangle{}

	f.spores = NewSpores(int32(x), int32(y))
	f.spores.AddTags(HitboxData)

	f.SetData(f)
	f.AddTags(common.TagEnemy, TagHurtbox)

	f.attackRangeAOE = f.Hurtbox.X - f.attackDistance
	return f
}

// Update FungalBosses every frame, checking attack and validating rigidbody collision
func (f *FungalBoss) Update(dt float32) {
	f.Ase.Update(dt)
	f.aoeAse.Update(dt)
	f.determineAttack()

	// Movement of spores
	f.spores.Update(dt)

}

func (f *FungalBoss) determineAttack() {
	//TODO determine when certain attacks will happen
	// attack 1 must not have happened recently (spores still falling)
	// attack 2 must be in range
	// attack 1 - standard shooting of fungal spore, they will drop
	// down with time
	f.sporeAttack()

	if time.Since(f.attackAoeBefore) >= time.Millisecond*f.attackAoeTimer {
		f.attackAoeBefore = time.Now()
		//attack 2 - ground AOE - when player is close enough to boss
		// attack that comes out from boss to set distance, attacks only left
		//TODO if player is in range
		if !f.isAttackingAOE {
			f.aoeAttack()
		}
	}

	// Updates hitbox of aoe attack
	if f.isAttackingAOE {
		x, y := f.HitboxAOE.GetXY()
		newHitboxX := x + int32(f.attackSpeedAOE)

		// If the attack has not met the range, update the hitbox
		if newHitboxX <= int32(f.attackRangeAOE) {
			f.HitboxAOE.Move(0, 0)
			f.isAttackingAOE = false
		} else {
			f.HitboxAOE.SetXY(newHitboxX, y)
		}
	}
}

func (f *FungalBoss) sporeAttack() {
	if time.Since(f.sporeBefore) >= time.Millisecond*f.sporeTimer {
		// Reset timer.
		f.sporeBefore = time.Now()
		log.Print("Shoot Spore")
		f.spores.CreateRow()
	}
}

func (f *FungalBoss) aoeAttack() {
	log.Print("aoeAttack")
	f.HitboxAOE = resolv.NewRectangle(
		f.Hurtbox.X, f.Hurtbox.Y,
		50, 50,
	)
	f.isAttackingAOE = true
}

// Draw sprite texture at given collision coordinates
func (f *FungalBoss) Draw() {
	// srcX, srcY resemble the X and Y pixels where the active sprite is.
	srcX, srcY := f.Ase.FrameBoundaries().X, f.Ase.FrameBoundaries().Y
	w, h := f.Ase.FrameBoundaries().Width, f.Ase.FrameBoundaries().Height

	// src resembles the cropped out area that the sprite is in the spritesheet.
	src := r.NewRectangle(float32(srcX), float32(srcY), float32(w), float32(h))

	// dest is the world position that the slime should appear in.
	dest := r.NewRectangle(float32(f.Hurtbox.X), float32(f.Hurtbox.Y), float32(w), float32(h))

	r.DrawTexturePro(
		f.Sprite, src, dest, r.NewVector2(0, 0), 0, r.White,
	)

	f.spores.Draw()

	// Draw aoe attack

	// Draw debug messages about the entity's current information
	f.debugDraw()
}

func (f *FungalBoss) debugDraw() {
	//Hurtbox
	r.DrawRectangleLines(
		int(f.Hurtbox.X), int(f.Hurtbox.Y),
		int(f.Hurtbox.W), int(f.Hurtbox.H),
		r.White,
	)

	// Ground Attack
	if f.isAttackingAOE {
		// fmt.Printf("X[%v], Y[%v]", f.HitboxAOE.X, f.HitboxAOE.Y)
		r.DrawRectangleLines(
			int(f.HitboxAOE.X), int(f.HitboxAOE.Y),
			int(f.HitboxAOE.W), int(f.HitboxAOE.H),
			r.Red,
		)
	}
}
