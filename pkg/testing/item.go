package testing

import "github.com/damienfamed75/endorem/pkg/item"

var (
	_ item.Effector = &Item{}
)

type Item struct {
}

func (*Item) AddEffect(*item.EffectData) {

}

func (*Item) RemoveEffect(*item.EffectData) {

}

func (*Item) String() string {
	return "testing item"
}
