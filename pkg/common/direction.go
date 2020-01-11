package common

//go:generate stringer -type=Direction -linecomment -output=direction_string.go

// Direction is used to tell the direction of an entity.
type Direction uint8

// Comments are used for generating the strings. Do not remove!
const (
	Left  Direction = iota + 1 // left
	Right                      // right
)
