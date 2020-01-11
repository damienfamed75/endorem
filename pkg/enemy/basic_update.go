package enemy

import (
	"time"

	"github.com/damienfamed75/endorem/pkg/common"
)

// Update is non drawing related functionality with the enemy.
func (b *Basic) Update() {
	b.move()
	b.attack()
}

func (b *Basic) attack() {
	// TODO - if player is in attackzone then try to attack.
	// Debugging:
	// Timer for attacks every half second.
	if time.Since(b.attackBefore) >= time.Millisecond*b.attackTimer {
		// Reset timer.
		b.attackBefore = time.Now()

		// Flip attack value.
		b.isAttacking = !b.isAttacking
		if b.isAttacking {
			// Re-add hurtbox to the enemy space and set position to enemy.
			// b.Hitbox.SetXY(b.Collision.X, b.Collision.Y+b.Collision.H/3.0)
			// Based on the direction the player is facing, set the position of the
			// hitbox in front of the player.
			if b.Facing == common.Left {
				b.Hitbox.SetXY(b.Collision.X-(b.Hitbox.W/2), b.Collision.Y+b.Collision.H/3.0)
			} else {
				b.Hitbox.SetXY(b.Collision.X, b.Collision.Y+b.Collision.H/3.0)
			}

			b.Add(b.Hitbox)
		} else {
			// Remove hurtbox from enemy space.
			b.Remove(b.Hitbox)
			b.state = common.StateIdle
		}
	}
}

func (b *Basic) move() {
	// idle walking.
	if !b.PlayerSeen {
		b.idleWalk()
	} else {
		// TODO - chase player (day 2)
	}
}

func (b *Basic) idleWalk() {
	// Wait for the enemy to sit for a bit at their destination.
	if time.Since(b.destinationMetTime) <= time.Millisecond*b.waitTime {
		return
	}

	friction := float32(0.5)
	accel := (0.5 + friction) * float32(b.direction)

	maxSpd := float32(1)
	if b.SpeedX > friction {
		b.SpeedX -= friction
	} else if b.SpeedX < -friction {
		b.SpeedX += friction
	} else {
		b.SpeedX = 0
	}

	// if met destination on X, turn around
	for i, d := range b.Destinations {
		if b.Collision.X == int32(d.X) && b.LastDestination != i {
			b.direction *= -1
			b.LastDestination = i
			b.destinationMetTime = time.Now() // Reset wait timer.
		}
	}
	b.SpeedX += accel

	if b.SpeedX > maxSpd {
		b.SpeedX = maxSpd
	}
	if b.SpeedX < -maxSpd {
		b.SpeedX = -maxSpd
	}

	if b.direction > 0 {
		b.Facing = common.Right
	} else {
		b.Facing = common.Left
	}

	x := int32(b.SpeedX)

	b.Collision.Move(x, 0)
}
