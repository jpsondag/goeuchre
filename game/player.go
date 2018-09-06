package game

import "github.com/jpsondag/goeuchre/card"

// Player is an ID + the hand.
type Player struct {
	Name string
	Hand card.Hand
}

// NewPlayer returns an initialized player.
func NewPlayer(n string) Player {
	return Player{Name: n}
}
