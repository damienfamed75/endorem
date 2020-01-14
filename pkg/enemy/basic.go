package enemy

import (
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.Entity = &Basic{}
)

// Basic is a testing enemy that is very basic in attacks and features.
type Basic struct {
	Health          int
	SpeedX          float32 // speed of basic in X-direction
	SpeedY          float32
	IsDead          bool
	PlayerSeen      bool // If the enemy has spotted the enemy.
	ShouldAttack    bool
	AttackDistance  int32
	jumpHeight      int32
	Facing          common.Direction
	Origin          r.Vector2
	Destinations    [2]r.Vector2 // left and right destinations
	LastDestination int
	Sprite          r.Texture2D
	Collision       *resolv.Rectangle
	Hitbox          *resolv.Rectangle
	Ground          *resolv.Space

	direction          int8
	state              common.State
	isAttacking        bool
	speedMultiplier    float32   // Multiplier of the enemy's horizontal movement.
	attackBefore       time.Time // How often the enemy can attack (milliseconds)
	attackTimer        time.Duration
	healthBefore       time.Time
	invincibleTimer    time.Duration
	destinationMetTime time.Time
	waitTime           time.Duration

	MoveIncrement float64
	*resolv.Space
}

func setupBasic() *Basic {
	return &Basic{
		Sprite:             r.LoadTexture("assets/basicenemy.png"),
		Space:              resolv.NewSpace(),
		Health:             2 + common.GlobalConfig.Enemy.AddedHealth,
		SpeedX:             2,
		SpeedY:             2,
		AttackDistance:     30,
		direction:          1,
		jumpHeight:         8,
		Facing:             common.Right,
		state:              common.StateIdle,
		speedMultiplier:    common.GlobalConfig.Enemy.MoveSpeedMultiplier,
		attackTimer:        time.Duration(common.GlobalConfig.Enemy.AttackTimer),
		attackBefore:       time.Now(),
		healthBefore:       time.Now(),
		destinationMetTime: time.Now(),
		waitTime:           time.Duration(common.GlobalConfig.Enemy.WaitTime),
		invincibleTimer:    time.Duration(common.GlobalConfig.Enemy.InvincibleTimer),
	}
}

// NewBasic returns a configured basic enemy at the given coordinates.
func NewBasic(x, y int, ground *resolv.Space) *Basic {
	b := setupBasic()

	b.Origin = r.NewVector2(float32(x), float32(y))
	b.Destinations = [2]r.Vector2{
		r.NewVector2(
			float32(int32(x)-b.Sprite.Height*2)+float32(b.Sprite.Width/2),
			float32(y)+float32(b.Sprite.Height),
		),
		r.NewVector2(
			float32(int32(x)+b.Sprite.Height*2)+float32(b.Sprite.Width/2),
			float32(y)+float32(b.Sprite.Height),
		),
	}

	// Set the hit and hurt boxes.
	b.Collision = resolv.NewRectangle(
		int32(x), int32(y),
		b.Sprite.Width, b.Sprite.Height,
	)
	b.Hitbox = resolv.NewRectangle(
		0, 0, b.Sprite.Height, b.Sprite.Width,
	)
	b.Ground = ground

	// Add the collision boxes to the enemy space.
	b.Add(b.Collision, b.Hitbox)
	b.SetData(b)

	// Set the hitbox data to be different from the hitbox data.
	b.Collision.AddTags(TagHurtbox)
	b.Hitbox.SetData(HitboxData)

	// Tag this enemy as an enemy.
	b.AddTags(common.TagEnemy)
	b.Hitbox.AddTags(common.TagEnemy)

	return b
}

// Draw is used for raylib exclusive drawing function calls.
func (b *Basic) Draw() {
	// Draw the enemy texture.
	r.DrawTexture(b.Sprite, int(b.Collision.X), int(b.Collision.Y), r.White)

	b.debugDraw() //DEBUG
}

func (b *Basic) TakeDamage() {
	if time.Since(b.healthBefore) >= time.Millisecond*b.invincibleTimer {
		b.healthBefore = time.Now()

		b.Health--
		if b.Health <= 0 {
			b.IsDead = true
			b.state = common.StateDead
			b.Clear()
		}
	}
}
