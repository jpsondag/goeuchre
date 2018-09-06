package card

import (
	"math/rand"
	"strings"
)

// Deck represents a deck of cards.
type Deck struct {
	Cards []Card
}

// DeepCopy returns a new instance of this object.
func (d *Deck) DeepCopy() Deck {
	ret := Deck{make([]Card, len(d.Cards))}
	copy(ret.Cards, d.Cards)
	return ret
}

// FiftyTwoCardDeck returns a full deck of cards.
func FiftyTwoCardDeck() Deck {
	deck := Deck{make([]Card, 0, 52)}
	c := Card{0, 0}
	for i := AceLow; i <= King; i++ {
		for _, suit := range Suits() {
			c.Suit = suit
			c.Value = i
			deck.Cards = append(deck.Cards, c)
		}
	}
	return deck
}

func swap(s []Card, x int, y int) {
	// Cache before overwriting.
	old := s[y]

	// Replace cached with new.
	s[y] = s[x]

	// Put cached value in x.
	s[x] = old
}

// Shuffle shuffles the cards in the deck.
func (d *Deck) Shuffle() {
	// This shuffle isn't great.
	for idx := range d.Cards {
		// Choose a random new location for this card.
		newIndex := rand.Intn(len(d.Cards))
		swap(d.Cards, idx, newIndex)
	}
}

// Pop removes the top n cards from the deck. If n < len(cards) remaining cards are returned.
func (d *Deck) Pop(n int) []Card {
	var ret []Card
	if n < len(d.Cards) {
		ret = d.Cards[:n]
	} else {
		ret = d.Cards[:]
	}
	d.Cards = d.Cards[len(ret):]
	return ret
}

func (d *Deck) dealHand(numCards int) Hand {
	if len(d.Cards) < numCards {
		panic(strings.Join([]string{"not enough cards in deck found", string(len(d.Cards)), " need ", string(numCards)}, ""))
	}
	hand := Hand{d.Cards[0:numCards]}
	d.Cards = d.Cards[numCards:]
	return hand
}

// Deal creates 'n' Hands from the deck. If the deck
func (d *Deck) Deal(numHands int, cardsPerHand int) []Hand {
	hands := make([]Hand, 0, numHands)
	for i := 0; i < numHands; i++ {
		hands = append(hands, d.dealHand(cardsPerHand))
	}
	return hands
}
