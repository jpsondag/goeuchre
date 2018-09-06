package card_test

import (
	"fmt"
	"testing"

	"github.com/jpsondag/goeuchre/card"
)

// Make test hand:
//    3H, 4C, KS, 9D
func testHand() card.Hand {
	hand := card.Hand{Cards: make([]card.Card, 0, 4)}

	hand.Cards = append(hand.Cards,
		card.Card{Value: 2, Suit: card.Heart},
		card.Card{Value: 4, Suit: card.Club},
		card.Card{Value: card.King, Suit: card.Spade},
		card.Card{Value: 9, Suit: card.Diamond})
	return hand
}

func TestFind(t *testing.T) {
	hand := testHand()
	if hand.Find(card.Card{Value: 2, Suit: card.Heart}) == -1 {
		t.Error("Should have 2H")
	}
	if hand.Find(card.Card{Value: 3, Suit: card.Heart}) != -1 {
		t.Error("Should not have 3H")
	}
}

func TestPlay(t *testing.T) {
	// Create a hand.
	hand := testHand()

	if len(hand.Cards) != 4 {
		t.Error("Should have 4 cards in hand")
	}
	// Play the 2 of Hearts.
	if !hand.Play(card.Card{Value: 2, Suit: card.Heart}) {
		t.Error("Should have 2 of hearts")
	}
	fmt.Println(hand)
	if len(hand.Cards) != 3 {
		t.Error("Should have 3 cards in hand")
	}

	// Playing again should not work
	if hand.Play(card.Card{Value: 2, Suit: card.Heart}) {
		t.Error("Should not have 2 of hearts")
	}
}
