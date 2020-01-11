package common

//go:generate stringer -type=State -linecomment -output=state_string.go

type State uint8

// Do not edit comments.
const (
	StateIdle    State = iota + 1 // idle
	StateLeft                     // left
	StateRight                    // right
	StateJumping                  // jumping
	StateFalling                  // falling
	StateAttack                   // attack
	StateCrouch                   // crouch
	StateDead                     // dead
)
