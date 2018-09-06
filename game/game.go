package game

import (
	"github.com/jpsondag/goeuchre/card"
)

// Config is used to configure the game, this is passed into the constructor.
type Config struct {
	AcesHigh   bool
	ScoreToWin int
	Deck       card.Deck
}

// Game is made up of a set of teams.
type Game struct {
	Teams    []Team
	Deck     card.Deck
	AcesHigh bool
}

// New creates new Game with given configuration.
func New(config Config, teams []Team) Game {
	g := Game{teams, config.Deck, config.AcesHigh}
	return g
}

// Convert to string.
func (g Game) String() string {
	s := ""
	for _, t := range g.Teams {
		s += string(t.Score)
	}
	return s
}
