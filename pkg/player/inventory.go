package player

import "github.com/damienfamed75/endorem/pkg/item"

type Inventory struct {
	Effectors map[string]item.Effector
	Usables   map[string]item.Usable
}

func NewInventory() *Inventory {
	return &Inventory{
		Effectors: make(map[string]item.Effector),
		Usables:   make(map[string]item.Usable),
	}
}

func (i *Inventory) AddItem(newItem interface{}) {
	switch t := newItem.(type) {
	case item.Effector:
		i.Effectors[t.String()] = t
	case item.Usable:
		i.Usables[t.String()] = t
	}
}

func (i *Inventory) RemoveItem(itemName string) {
	if _, ok := i.Effectors[itemName]; ok {
		delete(i.Effectors, itemName)
	} else if _, ok := i.Usables[itemName]; ok {
		delete(i.Usables, itemName)
	}
}
