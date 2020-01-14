package enemy

import (
	"time"

	"github.com/SolarLune/resolv/resolv"

	"github.com/damienfamed75/endorem/pkg/common"
)

// Update is non drawing related functionality with the enemy.
func (b *Basic) Update(float32) {
	px, py := b.player.GetXY()

	dist := resolv.Distance(
		b.GetX(), b.GetY(),
		px, py,
	)

	b.PlayerSeen = dist < common.GlobalConfig.Enemy.VisionDistance
	b.ShouldAttack = dist < b.AttackDistance

	b.move()

	if b.ShouldAttack {
		b.attack()
	} else if !b.Hitbox.HasTags(TagAttackZone) {
		b.Hitbox.AddTags(TagAttackZone)
	}

	b.Rigidbody.Update()
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
			//b.Hitbox.SetXY(b.GetX(), b.GetY()+b.Sprite.Height/3.0)
			// Based on the direction the player is facing, set the position of the
			// hitbox in front of the player.

			// if b.Facing == common.Left {
			// 	b.Hitbox.SetXY(b.GetX()-(b.Hitbox.W/2), b.GetY()+b.Sprite.Height/3.0)
			// } else {
			// 	b.Hitbox.SetXY(b.GetX(), b.GetY()+b.Sprite.Height/3.0)
			// }

			b.Hitbox.SetXY(b.GetX()+(b.Hitbox.W/2*int32(b.getPlayerDirection()))-(b.Sprite.Width/2), b.Hitbox.Y)
			b.Hitbox.ClearTags()
			//	b.Add(b.Hitbox)
		} else {
			// Remove hurtbox from enemy space.
			//b.Remove(b.Hitbox)
			b.Hitbox.AddTags(TagAttackZone)
			//b.state = common.StateIdle
		}
	}
}

func (b *Basic) move() {
	// // idle walking.

	if !b.PlayerSeen {
		// if met destination on X, turn around
		b.idleWalk()
	} else {
		// TODO - chase player (day 2)
		b.chasePlayer()
	}

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

func (b *Basic) idleWalk() {
	// MoveIncrement not only changes the position, but influences time
	// b.MoveIncrement += 0.01
	// center := ((b.Destinations[1].X - b.Destinations[0].X) / 2)
	// b.Collision.X = int32(float64(center) + (math.Sin(b.MoveIncrement*math.Pi) * float64(center)))
}

func (b *Basic) chasePlayer() {
	// Move x-position towards player
	// TODO stop movement if in attack range
	b.Rigidbody.Velocity.X = b.travelSpeed * b.getPlayerDirection()

	// TODO Jump is obstacle is in enemy way
	// res := b.Resolve(b.Ground, int32(b.SpeedX), 0)

	// b.Collision.Y += int32(b.SpeedY)
	// if res.Teleporting {
	// 	b.Collision.Y -= b.jumpHeight
	// }

}
func (b *Basic) getPlayerDirection() float32 {
	px, _ := b.player.GetXY()
	if px > b.GetX() {
		return 1
	}

	return -1
}
func (b *Basic) tryToMove(x int) {

}
