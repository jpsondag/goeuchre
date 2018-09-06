package card

import (
	"fmt"
)

// Suit is the type for card suits (Heart/Club/Spade/Diamond)
type Suit int

// The 4 suits.
const (
	Heart   Suit = 1
	Club    Suit = 2
	Diamond Suit = 3
	Spade   Suit = 4
)

// Suits returns all 4 suits.
func Suits() [4]Suit {
	return [4]Suit{Heart, Club, Diamond, Spade}
}

// Card represents a (suit, value). Ace has value of 1.
type Card struct {
	Value int
	Suit  Suit
}

// Syntactic sugar for values of cards.
const (
	AceLow  int = 1
	Jack    int = 11
	Queen   int = 12
	King    int = 13
	AceHigh     = 14
)

func (s Suit) String() string {
	switch s {
	case Heart:
		return "H"

	case Spade:
		return "S"

	case Diamond:
		return "D"
	case Club:
		return "C"
	}
	return "impossible to reach"
}
func (c Card) String() string {
	return fmt.Sprintf("%v%s ", c.Value, c.Suit)
}
