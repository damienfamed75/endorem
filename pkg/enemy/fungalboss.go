package enemy

import (
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.Entity = &FungalBoss{}
)

// FungalBoss is the set boss of level one, this enemy will be placed
// at the end of the first level
type FungalBoss struct {
	Sprite r.Texture2D

	Health int
	isDead bool

	player *resolv.Space
	ground *resolv.Space

	maxSpeedX int32
	maxSpeedY int32

	attackSpeed    float32
	attackDistance float32
	attackRangeAOE float32

	attackBefore time.Time
	attackTimer  time.Duration // Time before making another attack(milliseconds)
	sporeBefore  time.Time     //How often the boss can shoot spores(milliseconds)
	sporeTimer   time.Duration
	spores       *Spores

	*common.Rigidbody
}

func setupFungalBoss() *FungalBoss {
	return &FungalBoss{
		Sprite:         r.LoadTexture("assets/fungalboss.png"),
		Health:         1 + common.GlobalConfig.Enemy.AddedHealth,
		maxSpeedX:      0,
		maxSpeedY:      0,
		attackSpeed:    10,
		attackDistance: 60,
		attackRangeAOE: 80,
		attackBefore:   time.Now(),
		attackTimer:    time.Duration(common.GlobalConfig.Enemy.AttackTimer),

		sporeBefore: time.Now(),
		sporeTimer:  time.Duration(2000), //REMOVE HARDCODED #
	}
}

func NewFungalBoss(x, y int, world *resolv.Space) *FungalBoss {
	f := setupFungalBoss()

	// Store important spaces in the world
	f.player = world.FilterByTags(common.TagPlayer)
	f.ground = world.FilterByTags(common.TagGround)

	collision := resolv.NewRectangle(
		int32(x), int32(y), f.Sprite.Width, f.Sprite.Height,
	)
	collision.AddTags(TagHurtbox, common.TagCollision)

	f.spores = NewSpores(int32(x), int32(y))

	// create Rigidbody
	f.Rigidbody = common.NewRigidbody(
		int32(x), int32(y),
		f.maxSpeedX, f.maxSpeedY, f.ground,
		collision,
	)

	f.SetData(f)
	f.AddTags(common.TagEnemy, TagHurtbox)

	return f
}

// Update FungalBosses every frame, checking attack and validating rigidbody collision
func (f *FungalBoss) Update(dt float32) {
	f.determineAttack()

	// Movement of spores
	f.spores.Update()

	f.Rigidbody.Update(dt)
}

func (f *FungalBoss) determineAttack() {
	//TODO determine when certain attacks will happen
	// attack 1 must not have happened recently (spores still falling)
	// attack 2 must be in range

	//TODO attack 1 - standard shooting of fungal spore, they will drop
	// down with time
	f.sporeAttack()
	//TODO attack 2 - ground AOE - when player is close enough to boss
	// attack that comes out from boss to set distance
	// f.aoeAttack()
}

func (f *FungalBoss) sporeAttack() {
	if time.Since(f.sporeBefore) >= time.Millisecond*f.sporeTimer {
		// Reset timer.
		f.sporeBefore = time.Now()
		//log.Print("Shoot Spore")
		f.spores.CreateRow()
	}
}

func (f *FungalBoss) aoeAttack() {

}

// Draw sprite texture at given collision coordinates
func (f *FungalBoss) Draw() {
	r.DrawTexture(f.Sprite, int(f.GetX()), int(f.GetY()), r.White)
	f.spores.Draw()

	// Draw debug messages about the entity's current information
	f.debugDraw()
}

func (f *FungalBoss) debugDraw() {

}
