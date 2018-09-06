package game

import (
	"github.com/jpsondag/goeuchre/card"
)

// Builder is Used for easily configuring a new game.
type Builder struct {
	acesHigh   bool
	scoreToWin int
	deck       card.Deck
	teams      []Team
}

// SetDeck initializes deck used for this game.
func (g *Builder) SetDeck(d card.Deck) *Builder {
	g.deck = d
	return g
}

// SetAcesHigh is used to tell if an ace is considered high or low for this game.
func (g *Builder) SetAcesHigh(v bool) *Builder {
	g.acesHigh = v
	return g
}

// SetScoreToWin initializes the score needed to finish the game.
func (g *Builder) SetScoreToWin(v int) *Builder {
	g.scoreToWin = v
	return g
}

// AddTeam adds a new team to this game.
func (g *Builder) AddTeam(t Team) *Builder {
	g.teams = append(g.teams, t)
	return g
}

// Build returns the game for this builder.
func (g Builder) Build() Game {
	c := Config{
		AcesHigh:   g.acesHigh,
		ScoreToWin: g.scoreToWin,
		Deck: g.deck,
	}
	return New(c, g.teams)
}
