package item

import "github.com/SolarLune/resolv/resolv"

type EffectData struct {
	CurrentHealth *int
	MaxHealth     *int
	SpeedX        *float32
	SpeedY        *float32

	Collision *resolv.Rectangle
	// TODO Add mask
}

type Effector interface {
	AddEffect(*EffectData)
	RemoveEffect(*EffectData)

	String() string
}
