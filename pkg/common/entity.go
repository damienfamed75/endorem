package common

import "github.com/SolarLune/resolv/resolv"

type Entity interface {
	Update(dt float32)
	Draw()

	resolv.Shape
}
