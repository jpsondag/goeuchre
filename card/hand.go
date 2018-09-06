package card

import (
	"strings"
)

// Hand is a set of cards.
type Hand struct {
	Cards []Card
}

// Play scans 'Cards' for 'card'
//   if found, the first instance of card is removed and true is returned
//   if not found false is returned and no change is made.
func (h *Hand) Play(card Card) bool {
	idx := h.Find(card)
	if idx == -1 {
		return false
	}
	// Swap with last element.
	h.Cards[idx] = h.Cards[len(h.Cards)-1]
	h.Cards = h.Cards[0 : len(h.Cards)-1]
	return true
}

// Find returns the index in this hand that contains 'card' otherwise returns -1.
func (h *Hand) Find(card Card) int {
	for idx, c := range h.Cards {
		if c == card {
			return idx
		}
	}
	return -1
}

// Count returns the number of cards in the hand with the given suit.
func (h *Hand) Count(s Suit) int {
	ret := 0
	for _, c := range h.Cards {
		if c.Suit == s {
			ret++
		}
	}
	return ret
}

func (h Hand) String() string {
	b := strings.Builder{}
	for _, c := range h.Cards {
		b.WriteString(c.String())
	}
	return b.String()
}
