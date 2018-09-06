package euchre

import (
	"fmt"

	"github.com/jpsondag/goeuchre/card"
	"github.com/jpsondag/goeuchre/game"
)

type circularArray struct {
	Head  *game.Player
	array [4]*game.Player
}

func (c *circularArray) nextHead() *game.Player {
	iter := c.iter()
	c.Head = iter.Next()
	return c.Head
}

func (c *circularArray) iter() circularArrayIt {
	return newCircularArrayIt(c)
}

func newCircularArray(input [4]*game.Player) circularArray {
	return circularArray{input[0], input}
}

type circularArrayIt struct {
	array *circularArray
	Val   *game.Player
	last  int
	Done  bool
}

func newCircularArrayIt(array *circularArray) circularArrayIt {
	ret := circularArrayIt{array, array.Head, 0, false}
	index := ret.getIndex()
	if index != 0 {
		ret.last = index - 1
	} else {
		ret.last = len(array.array) - 1
	}
	return ret
}

func (c *circularArrayIt) getIndex() int {
	for i, p := range c.array.array {
		if p == c.Val {
			return i
		}
	}
	panic("iterator out of sync with underlying array")
}

func (c *circularArrayIt) Next() *game.Player {
	index := c.getIndex()
	if index == c.last {
		c.Done = true
		return nil
	}
	// Return next, possibly going back to start.
	if index >= len(c.array.array)-1 {
		index = 0
	} else {
		index = index + 1
	}
	c.Val = c.array.array[index]
	return c.Val
}

// Euchre encapsulates a euchre game.
type Euchre struct {
	Game game.Game

	//	Players [4]*game.Player
	Players circularArray
	Dealer  int
}

func euchreDeck() card.Deck {
	c := card.Card{Value: 0, Suit: 0}
	deck := card.Deck{}
	for i := 9; i <= card.AceHigh; i++ {
		for _, suit := range card.Suits() {
			c.Suit = suit
			c.Value = i
			deck.Cards = append(deck.Cards, c)
		}
	}
	return deck
}

// New returns a Euchre object.
func New() Euchre {
	b := game.Builder{}

	t1 := game.NewTeam()
	t1.AddPlayer(game.NewPlayer("JP"))
	t1.AddPlayer(game.NewPlayer("Ellen"))

	t2 := game.NewTeam()
	t2.AddPlayer(game.NewPlayer("Connie"))
	t2.AddPlayer(game.NewPlayer("Ed"))

	return Euchre{b.SetAcesHigh(true).SetScoreToWin(10).SetDeck(euchreDeck()).AddTeam(t1).AddTeam(t2).Build(), newCircularArray([4]*game.Player{nil, nil, nil, nil}), 0}
}

// Run runs the euchre game to completion.
func (euchre *Euchre) Run() {

	// Initialize dealer.
	p1 := &euchre.Game.Teams[0].Players[0]
	p2 := &euchre.Game.Teams[1].Players[0]
	p3 := &euchre.Game.Teams[0].Players[1]
	p4 := &euchre.Game.Teams[1].Players[1]
	players := newCircularArray([4]*game.Player{p1, p2, p3, p4})
	euchre.Players = players

	euchre.Dealer = 0

	for !euchre.gameOver() {
		euchre.RunRound()
	}
	for _, team := range euchre.Game.Teams {
		fmt.Println("Team: ", team.Name(), " final score: ", team.Score)
	}
}

func (euchre *Euchre) gameOver() bool {
	for _, team := range euchre.Game.Teams {
		if team.Score >= 10 {
			return true
		}
	}
	return false
}

// RunRound runs a single hand of euchre.
func (euchre *Euchre) RunRound() {
	d := euchre.Game.Deck.DeepCopy()

	dealer := euchre.getDealer()
	fmt.Println("Dealer: ", dealer.Name)

	const numHands = 4
	const cardsPerHand = 5
	d.Shuffle()
	hands := d.Deal(numHands, cardsPerHand)
	assignHands(euchre.Players, hands)
	euchre.printHands()

	flipped := d.Pop(1)[0]
	fmt.Println("Flipped ", flipped)

	orderedBy := euchre.order(flipped)
	trump := flipped.Suit
	if orderedBy == nil {
		orderedBy, trump = euchre.chooseTrump(flipped)
	}
	euchre.playHand(orderedBy, trump)
}

func (euchre *Euchre) playHand(orderedBy *game.Player, trump card.Suit) {
	// Make copy of players.
	players := euchre.Players
	teamToTricks := make(map[string]int)
	for len(players.Head.Hand.Cards) != 0 {
		table := [4]card.Card{nilCard(), nilCard(), nilCard(), nilCard()}
		winner := euchre.playTrick(&players, trump, table)

		fmt.Println(winner.Name, " won that trick")

		for _, team := range euchre.Game.Teams {
			for _, player := range team.Players {
				if player.Name == winner.Name {
					teamToTricks[team.Name()]++
				}
			}
			fmt.Println("Team ", team.Name(), " Tricks: ", teamToTricks[team.Name()])
		}

		for players.Head != winner {
			players.nextHead()
		}
	}

	winningTeam := ""
	winningTricks := -1
	for teamName, tricks := range teamToTricks {
		fmt.Println("Team ", teamName, " : ", tricks)
		if winningTricks < tricks {
			winningTeam = teamName
			winningTricks = tricks
		}
	}
	for i := range euchre.Game.Teams {
		team := &euchre.Game.Teams[i]
		if team.Name() == winningTeam {
			team.Score++
		}
		fmt.Println("Team ", team.Name(), " Score: ", team.Score)
	}

}

func (euchre *Euchre) playTrick(players *circularArray, trump card.Suit, table [4]card.Card) *game.Player {
	it := players.iter()
	lead := euchre.lead(it.Val, trump)
	if !it.Val.Hand.Play(lead) {
		panic("No available cards for lead?")
	}

	fmt.Println(it.Val.Name, " leads with  ", lead)
	it.Next()

	table[0] = lead
	fmt.Println("Table: ", table)

	contender := lead
	winner := it.Val
	i := 1
	for ; !it.Done; it.Next() {
		chosenCard := euchre.follow(it.Val, trump, table)
		if !it.Val.Hand.Play(chosenCard) {
			panic("follow should always return a valid card")
		}
		fmt.Println(it.Val.Name, " plays ", chosenCard)
		table[i] = chosenCard
		fmt.Println("Table: ", table)
		i++

		// Update contender/ winner.
		if cardLessThan(trump, contender, chosenCard) {
			winner = it.Val
			contender = chosenCard
		}
	}
	return winner
}

func (euchre *Euchre) lead(player *game.Player, trump card.Suit) card.Card {
	hand := player.Hand
	if hasLeftBower(hand, trump) && hasRightBower(hand, trump) {
		return rightBower(trump)
	}
	highestOffsuit := euchre.highestOffsuit(hand, trump)
	if highestOffsuit != nilCard() {
		return highestOffsuit
	}
	// No offsuits, choose highest trump.
	return euchre.highestTrump(hand, trump)
}

func (euchre *Euchre) follow(player *game.Player, trump card.Suit, table [4]card.Card) card.Card {
	// lead is table[0]
	lead := table[0]
	leadSuit := lead.Suit
	if lead == leftBower(trump) {
		leadSuit = leftSuit(lead.Suit)
	}

	// Find me.
	me := -1
	for i := 0; i < len(table); i++ {
		if table[i] == nilCard() {
			me = i
			break
		}
	}

	if me == -1 {
		panic("follow called when table is full")
	}

	currentWinner := currentWinningCard(trump, table)
	partnerWinning := false
	if me == 2 || me == 3 {
		partner := me - 2
		partnerWinning = (currentWinner == table[partner])
	}

	// Try to beat card.
	mustFollowSuit := euchre.countSuit(player, leadSuit, trump) != 0
	if mustFollowSuit {
		contender := currentWinner
		if !partnerWinning {
			// If can beat do it. This will play the highest card of the given suit.
			for _, c := range player.Hand.Cards {
				if c.Suit == leadSuit && cardLessThan(trump, contender, c) {
					contender = c
				}
			}
			if contender != currentWinner {
				return contender
			}
		}
		contender = nilCard()

		// Can't beat or partner is winning choose lowest of that suit.
		for _, c := range player.Hand.Cards {
			if c.Suit == leadSuit && (contender == nilCard() || cardLessThan(trump, c, contender)) {
				contender = c
			}
		}
		if contender == nilCard() {
			panic("countTrump is non-zero so there should always be a card here.")
		}
		return contender
	}

	// Do not have to follow suit, still try to win.
	contender := currentWinner
	if !partnerWinning {
		for _, c := range player.Hand.Cards {
			if cardLessThan(trump, contender, c) {
				contender = c
			}
		}
	}

	if contender == currentWinner {
		contender = nilCard()
		// can't beat and do not have to follow suit choose the lowest valued card there is.
		for _, c := range player.Hand.Cards {
			if contender == nilCard() || cardLessThan(trump, c, contender) {
				contender = c
			}
		}
	}
	return contender
}

// Returns if contender < c, assumes that contender is currently winning so it is either the lead card, or trump.
func cardLessThan(trump card.Suit, contender card.Card, c card.Card) bool {
	// New contender if:

	// Special case left bower.
	if contender == leftBower(trump) && c != rightBower(trump) {
		return false
	}
	if c == leftBower(trump) && contender != rightBower(trump) {
		return true
	}

	// 1 ) Trumped non-trump contender
	trumpedContender := c.Suit == trump && contender.Suit != trump

	// 2) Played higher value card of same suit off contender (handles trump/trump case)
	higherValueThanContender := c.Suit == contender.Suit && c.Value > contender.Value

	return trumpedContender || higherValueThanContender
}

// Returns the current winning card on the table with the given trump.
func currentWinningCard(trump card.Suit, table [4]card.Card) card.Card {
	contender := nilCard()
	for _, c := range table {
		if contender == nilCard() {
			contender = c
			continue
		}
		if c == nilCard() {
			continue
		}

		// Compare c to contender.
		if cardLessThan(trump, contender, c) {
			contender = c
		}
	}
	if contender == nilCard() {
		panic("Should always have a contender.")
	}
	return contender
}

func nilCard() card.Card {
	return card.Card{Value: 0, Suit: card.Diamond}
}

// May return nilCard()
func (euchre *Euchre) highestOffsuit(hand card.Hand, trump card.Suit) card.Card {
	// May return nil
	contender := nilCard()
	for _, c := range hand.Cards {
		if c == leftBower(trump) || c.Suit == trump {
			continue
		}
		if contender.Value < c.Value {
			contender = c
		}
	}
	return contender
}

func (euchre *Euchre) highestTrump(hand card.Hand, trump card.Suit) card.Card {
	if hasRightBower(hand, trump) {
		return rightBower(trump)
	}
	if hasLeftBower(hand, trump) {
		return leftBower(trump)
	}
	contender := nilCard()
	for _, c := range hand.Cards {
		if contender.Value < c.Value {
			contender = c
		}
	}
	return contender

}

func (euchre *Euchre) getDealer() *game.Player {
	var p *game.Player
	for it := euchre.Players.iter(); !it.Done; it.Next() {
		p = it.Val
	}
	return p
}

func (euchre *Euchre) printHands() {
	for it := euchre.Players.iter(); !it.Done; it.Next() {
		fmt.Println(it.Val.Name, ": ", it.Val.Hand)
	}
}

// Decide whether or not to make 'card' trump. If ordered the player
// who called trump is returned, otherwise nil is returned.
func (euchre *Euchre) order(flipped card.Card) *game.Player {
	var orderedBy *game.Player
	for it := euchre.Players.iter(); !it.Done; it.Next() {
		if euchre.shouldCall(it.Val, flipped.Suit) {
			fmt.Println(it.Val.Name, " calls ", flipped)
			orderedBy = it.Val
			break
		} else {
			fmt.Println(it.Val.Name, " passes")
		}
	}
	return orderedBy
}

func (euchre *Euchre) chooseTrump(flipped card.Card) (*game.Player, card.Suit) {
	for it := euchre.Players.iter(); !it.Done; it.Next() {
		p := it.Val
		for _, s := range card.Suits() {
			if flipped.Suit == s {
				continue
			}
			if euchre.shouldCall(p, s) {
				fmt.Println(it.Val.Name, " calls ", s)
				return p, s
			}
		}
		fmt.Println(it.Val.Name, " passes")
	}
	panic("Need to implement skipping deal")
}

func (euchre *Euchre) shouldCall(p *game.Player, s card.Suit) bool {
	numTrump := euchre.countSuit(p, s, s)
	if p == euchre.getDealer() {
		numTrump++
	}
	return numTrump >= 3 || (hasLeftBower(p.Hand, s) && hasRightBower(p.Hand, s))
}

func (euchre *Euchre) countSuit(p *game.Player, s card.Suit, trump card.Suit) int {
	hand := p.Hand
	count := hand.Count(s)

	// Add one for left bower
	hasLeft := hand.Find(leftBower(trump)) != -1
	if s == leftSuit(trump) && hasLeft {
		count++
	}

	return count
}

func leftSuit(s card.Suit) card.Suit {
	switch s {
	case card.Heart:
		return card.Diamond
	case card.Diamond:
		return card.Heart
	case card.Club:
		return card.Spade
	case card.Spade:
		return card.Club
	}
	panic("unkonwn suit")
}

func leftBower(trump card.Suit) card.Card {
	return card.Card{Value: card.Jack, Suit: leftSuit(trump)}
}

func rightBower(trump card.Suit) card.Card {
	return card.Card{Value: card.Jack, Suit: trump}
}

func hasLeftBower(hand card.Hand, s card.Suit) bool {
	return hand.Find(leftBower(s)) != -1
}

func hasRightBower(hand card.Hand, s card.Suit) bool {
	return hand.Find(rightBower(s)) != -1
}

func assignHands(players circularArray, hands []card.Hand) {
	i := 0
	for it := players.iter(); !it.Done; it.Next() {
		it.Val.Hand = hands[i]
		i++
	}
}
