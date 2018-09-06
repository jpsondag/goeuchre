package game

// Team is atleast one player and a score.
type Team struct {
	Players []Player
	Score   int
}

// NewTeam returns a new team with x players.
func NewTeam() Team {
	return Team{make([]Player, 0), 0}
}

// AddPlayer adds a player to the team.
func (t *Team) AddPlayer(p Player) {
	t.Players = append(t.Players, p)
}

// Name returns the team name which is just the player names
func (t *Team) Name() string {
	ret := ""
	for _, p := range t.Players {
		ret += "-" + p.Name
	}
	return ret
}
