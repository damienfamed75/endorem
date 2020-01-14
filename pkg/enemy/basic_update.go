package enemy

import (
	"math"
	"time"

	"github.com/damienfamed75/endorem/pkg/common"
)

// Update is non drawing related functionality with the enemy.
func (b *Basic) Update(float32) {
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
	// // idle walking.
	if !b.PlayerSeen {
		// if met destination on X, turn around

	} else {
		// TODO - chase player (day 2)
	}
	b.walk()
	// for i, d := range b.Destinations {
	// 	if b.current.X == d.X && b.LastDestination != i {
	// 		//log.Print("change direction")
	// 		b.direction *= -1
	// 		b.LastDestination = i
	// 		b.destinationMetTime = time.Now() // Reset wait timer.
	// 	} else {
	// 		target = d
	// 	}
	// }

}

func (b *Basic) walk() {
	b.MoveIncrement += 0.01
	center := ((b.Destinations[1].X - b.Destinations[0].X) / 2)
	b.Collision.X = int32(float64(center) + (math.Sin(b.MoveIncrement*math.Pi) * float64(center)))
}
