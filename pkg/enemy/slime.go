package enemy

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/damienfamed75/aseprite"
	"github.com/damienfamed75/endorem/pkg/common"

	"github.com/SolarLune/resolv/resolv"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Slime struct {
	Ase    *aseprite.File
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
	targetAnimation  string

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
	var err error

	s.Ase, err = aseprite.Open("assets/slime.json")
	if err != nil {
		log.Fatal(err)
	}

	// Queues a default animation.
	s.Ase.Play("idle")

	// Store the important spaces in the world.
	s.player = world.FilterByTags(common.TagPlayer)
	s.ground = world.FilterByTags(common.TagGround)

	collision := resolv.NewRectangle(
		int32(x), int32(y), int32(s.Ase.FrameBoundaries().Width), int32(s.Ase.FrameBoundaries().Height),
	)
	collision.AddTags(TagHurtbox)

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
	if time.Since(s.healthBefore) >= time.Millisecond*s.invincibleTimer && !s.IsDead {
		s.healthBefore = time.Now()

		s.Health--
		// s.Rigidbody.Velocity.X = -(s.travelSpeed * s.getPlayerDirection())
		// s.Rigidbody.Velocity.Y = -s.travelSpeed
		// s.attacking = true
		// TODO push slime back from the player and play hurt animation.
		// s.playPriorityAnimation("damage")

		if s.Health <= 0 {
			s.IsDead = true
		}
	}
}

func (s *Slime) playPriorityAnimation(anim string) {
	s.Ase.Play(anim)
	s.targetAnimation = anim
}

func (s *Slime) playingTargetAnimation() bool {
	if s.Ase.IsPlaying(s.targetAnimation) {
		return s.Ase.AnimationFinished()
	}

	return false
}

// Update gets called every frame and tells if the slime is going to
// sit still, idle travel around, jump to the player, or attack.
func (s *Slime) Update(dt float32) {
	s.Ase.Update(dt)

	if s.IsDead {
		if s.waitAndPlay("death") {
			s.Clear()
		}

		return
	}

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
		if !s.playingTargetAnimation() {
			s.Ase.Play("idle")
		}

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
		s.attacking = true
		s.Rigidbody.Velocity.X = s.travelSpeed * s.getPlayerDirection()
	} else if !s.attacking {
		s.Rigidbody.Velocity.X = 0
	}
}

func (s *Slime) getPlayerDirection() float32 {
	px, _ := s.player.GetXY()

	if px > s.GetX() {
		return 1
	}

	return -1
}

// jump is the slime's main form of movement.
func (s *Slime) jump(height float32) {
	// If the slime is on the ground and has waited long enough to jump
	// again then perform a jump.
	if s.OnGround() && time.Since(s.jumpTimeBegin) > time.Millisecond*s.jumpTimer {
		if s.playingTargetAnimation() {
			return
		}

		if s.waitAndPlay("jump") {
			// Reset jump timer.
			s.jumpTimeBegin = time.Now()

			// Set the vertical speed of the slime.
			s.Rigidbody.Velocity.Y = height
			// Update that the slime is not on the ground anymore.
			s.onGround = false
		}
	} else if s.OnGround() {
		if !s.playingTargetAnimation() {
			s.Ase.Play("idle")
		}
	}
}

// waitAndPlay queues an animation to be played and
// returns false until it has finishes.
func (s *Slime) waitAndPlay(anim string) bool {
	// If the wanted animation is already playing.
	if s.Ase.IsPlaying(anim) {
		// And the animation has finished.
		if s.Ase.AnimationFinished() {
			return true
		}

		return false
	}

	// Play the animation if it's not playing.
	s.Ase.Play(anim)

	return false
}

// Draw the sprite texture at the collision box coordinates.
func (s *Slime) Draw() {
	// srcX, srcY resemble the X and Y pixels where the active sprite is.
	srcX, srcY := s.Ase.FrameBoundaries().X, s.Ase.FrameBoundaries().Y
	w, h := s.Ase.FrameBoundaries().Width, s.Ase.FrameBoundaries().Height

	// src resembles the cropped out area that the sprite is in the spritesheet.
	src := r.NewRectangle(float32(srcX), float32(srcY), float32(w), float32(h))

	// dest is the world position that the slime should appear in.
	dest := r.NewRectangle(float32(s.GetX()), float32(s.GetY()), float32(w), float32(h))

	r.DrawTexturePro(
		s.Sprite, src, dest, r.NewVector2(0, 0), 0, r.White,
	)

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
		int(s.Ase.FrameBoundaries().Width), int(s.Ase.FrameBoundaries().Width),
		r.Red,
	)
}
