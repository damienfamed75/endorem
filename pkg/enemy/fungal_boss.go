package enemy

import (
	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

// FungalBoss is the set boss of level one, this enemy will be placed
// at the end of the first level
type FungalBoss struct {
	Sprite r.Texture2D

	Health int
	isDead bool

	player *resolv.Space
	ground *resolv.Space

	maxSpeedX      int32
	maxSpeedY      int32
	attackSpeed    float32
	attackDistance float32
	attackRangeAOE float32

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
func (f *FungalBoss) Update(float32) {
	f.attack()
	f.Rigidbody.Update()
}

func (f *FungalBoss) attack() {
	//TODO determine when certain attacks will happen
	// (i.e attack 2 must be in range)

	//TODO attack 1 - standard shooting of fungal spore, they will drop
	// down with time

	//TODO attack 2 - ground AOE - when player is close enough to boss
	// attack that comes out from boss to set distance
}

// Draw sprite texture at given collision coordinates
func (f *FungalBoss) Draw() {
	r.DrawTexture(f.Sprite, int(f.GetX()), int(f.GetY()), r.White)

	// Draw debug messages about the entity's current information
	f.debugDraw()
}

func (f *FungalBoss) debugDraw() {

}
