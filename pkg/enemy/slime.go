package enemy

import (
	"fmt"
	"strconv"
	"time"

	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Slime struct {
	Sprite r.Texture2D

	Health int
	IsDead bool

	player *resolv.Space
	ground *resolv.Space

	onGround         bool
	attacking        bool
	travelSpeed      float32
	attackSpeed      float32
	attackDistance   int32
	speedY           float32
	speedX           float32
	gravity          float32
	travelJumpHeight float32
	attackJumpHeight float32
	maxSpeedX        int32
	maxSpeedY        int32

	playerSeen    bool
	jumpTimeBegin time.Time
	jumpTimer     time.Duration

	invincibleTimer time.Duration
	healthBefore    time.Time

	*common.Rigidbody
}

func setupSlime() *Slime {
	return &Slime{
		Sprite:           r.LoadTexture("assets/slime.png"),
		Health:           1 + common.GlobalConfig.Enemy.AddedHealth,
		attackDistance:   60,
		gravity:          common.GlobalConfig.Game.Gravity,
		jumpTimer:        1000,
		travelSpeed:      1,
		attackSpeed:      4,
		maxSpeedX:        8,
		maxSpeedY:        6,
		attackJumpHeight: -4,
		travelJumpHeight: -4,
		jumpTimeBegin:    time.Now(),
		invincibleTimer:  time.Duration(common.GlobalConfig.Enemy.InvincibleTimer),
	}
}

// NewSlime creates a slime at the given position.
func NewSlime(x, y int, world *resolv.Space) *Slime {
	s := setupSlime()

	// Store the important spaces in the world.
	s.player = world.FilterByTags(common.TagPlayer)
	s.ground = world.FilterByTags(common.TagGround)

	collision := resolv.NewRectangle(
		int32(x), int32(y), s.Sprite.Width, s.Sprite.Height,
	)
	collision.AddTags(TagHurtbox, common.TagCollision)

	s.Rigidbody = common.NewRigidbody(
		int32(x), int32(y),
		s.maxSpeedX, s.maxSpeedY, s.ground,
		collision,
	)

	s.SetData(s)

	s.AddTags(common.TagEnemy, TagHurtbox)

	return s
}

func (s *Slime) TakeDamage() {
	if time.Since(s.healthBefore) >= time.Millisecond*s.invincibleTimer {
		s.healthBefore = time.Now()

		s.Health--
		if s.Health <= 0 {
			s.IsDead = true
			s.Clear()
		}
	}
}

// Update gets called every frame and tells if the slime is going to
// sit still, idle travel around, jump to the player, or attack.
func (s *Slime) Update(float32) {
	px, py := s.player.GetXY()

	dist := resolv.Distance(
		s.GetX(), s.GetY(),
		px, py,
	)

	s.playerSeen = dist < common.GlobalConfig.Enemy.VisionDistance

	if s.attacking {
		s.attacking = !s.OnGround()
	}

	if s.playerSeen {
		// if the slime sees the player and is in attacking distance
		// then jump at the player.
		if dist < s.attackDistance {
			s.attack()
		} else {
			// if the slime can see the player but isn't close enough
			// to attack then try to follow the player.
			s.followPlayer()
		}
	} else {
		s.Rigidbody.Velocity.X = 0
	}

	// Update the slime's position according to gravity and checks collisions
	s.Rigidbody.Update()
}

func (s *Slime) attack() {
	// Try to make a small jump.
	s.jump(s.attackJumpHeight)

	if !s.OnGround() && !s.attacking {
		s.attacking = true
		s.Rigidbody.Velocity.X = s.attackSpeed * s.getPlayerDirection()
	}
}

func (s *Slime) followPlayer() {
	// Try to make a small jump.
	s.jump(s.travelJumpHeight)

	if !s.OnGround() && !s.attacking {
		s.Rigidbody.Velocity.X = s.travelSpeed * s.getPlayerDirection()
	} else if !s.attacking {
		s.Rigidbody.Velocity.X = 0
	}
}

func (s *Slime) getPlayerDirection() float32 {
	px, _ := s.player.GetXY()
	// if px > s.GetX() {
	if px > s.GetX() {
		return 1
	}

	return -1
}

// jump is the slime's main form of movement.
func (s *Slime) jump(height float32) bool {
	// If the slime is on the ground and has waited long enough to jump
	// again then perform a jump.
	if s.OnGround() && time.Since(s.jumpTimeBegin) > time.Millisecond*s.jumpTimer {
		// if s.onGround && time.Since(s.jumpTimeBegin) > time.Millisecond*s.jumpTimer {
		// Reset jump timer.
		s.jumpTimeBegin = time.Now()

		// Set the vertical speed of the slime.
		s.Rigidbody.Velocity.Y = height
		// Update that the slime is not on the ground anymore.
		s.onGround = false
	}

	return !s.onGround
}

// Draw the sprite texture at the collision box coordinates.
func (s *Slime) Draw() {
	// sx, sy := s.GetXY()
	// r.DrawTexture(s.Sprite, int(sx), int(sy), r.White)
	r.DrawTexture(s.Sprite, int(s.GetX()), int(s.GetY()), r.White)

	// Draw debug messages about the entity's current information.
	s.debugDraw()
}

func (s *Slime) debugDraw() {
	// Draw health.
	r.DrawText(
		"HP: "+strconv.Itoa(s.Health),
		int(s.GetX()), int(s.GetY()-(s.Sprite.Width/2)), 10,
		r.White,
	)

	r.DrawText(
		fmt.Sprintf("G[%v]", s.OnGround()),
		int(s.GetX()), int(s.GetY()+(s.Sprite.Height)), 10,
		r.White,
	)

	r.DrawText(
		fmt.Sprintf("V[%v,%v]", s.Velocity.X, s.Velocity.Y),
		int(s.GetX()), int(s.GetY()+(s.Sprite.Height)+10), 10,
		r.White,
	)

	// Draw the collision box for debugging reasons.
	r.DrawRectangleLines(
		int(s.GetX()), int(s.GetY()),
		int(s.Sprite.Width), int(s.Sprite.Height),
		r.Red,
	)
}
