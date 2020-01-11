package main

import (
	"github.com/SolarLune/resolv/resolv"
)

type Game struct {
	world *resolv.Space
}

func setupGame() *Game {
	return &Game{
		world: resolv.NewSpace(),
	}
}

func NewGame() *Game {
	g := setupGame()

	return g
}

func (g *Game) Update() {

}
