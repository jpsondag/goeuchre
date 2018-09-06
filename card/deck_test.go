package card_test

import (
	"testing"

	"github.com/jpsondag/goeuchre/card"
)

func TestDeck(t *testing.T) {
	d := card.FiftyTwoCardDeck()
	if len(d.Cards) != 52 {
		t.Error("Expected 52 cards got: ", len(d.Cards))
	}
}

func TestShuffle(t *testing.T) {
	d := card.FiftyTwoCardDeck()

	// Copy
	dCopy := d.DeepCopy()

	d.Shuffle()
	var diffs int
	for idx, c := range d.Cards {
		if dCopy.Cards[idx].Value != c.Value {
			diffs++
		}
	}
	if diffs == 0 {
		t.Error("Expected diffs got: ", diffs)
	}
}

func TestDeal(t *testing.T) {
	d := card.FiftyTwoCardDeck()

	// 5 hands 10 cards each.
	hands := d.Deal(5, 10)

	for _, hand := range hands {
		if len(hand.Cards) != 10 {
			t.Error("Expected 10 got: ", len(hand.Cards))
		}
	}
	if len(d.Cards) != 2 {
		t.Error("Expected 2 got: ", len(d.Cards))
	}
}
