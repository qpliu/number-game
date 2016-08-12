package main

import (
	"fmt"
)

type ownedcard struct {
	revealed bool
	starting bool
	card     card
}

type grid [3][3]ownedcard

type player struct {
	cards []ownedcard
	grid  grid
	cash  int
}

type state struct {
	deck          []card
	lot           [3]card
	players       []player
	currentPlayer int
}

func (card ownedcard) String() string {
	if card.revealed {
		return " " + card.card.String() + " "
	} else {
		return "#" + card.card.String() + "#"
	}
}

func (grid grid) String() string {
	str := "+---------------------------+\n"
	for _, row := range grid {
		str += "|"
		for _, c := range row {
			str += c.String()
		}
		str += "|\n"
	}
	str += "+---------------------------+\n"
	return str
}

func (player player) String() string {
	s := fmt.Sprintf("$%d cards: ", player.cash)
	for _, card := range player.cards {
		s += card.String()
	}
	s += "\n" + player.grid.String()
	return s
}

func setup(nplayers int) *state {
	state := &state{}
	deck, initialCards := shuffleAndDeal(nplayers)
	state.deck = deck
	state.players = make([]player, nplayers)
	for i := range state.players {
		state.players[i].cash = 10
		for _, card := range initialCards[i] {
			state.players[i].cards = append(state.players[i].cards, ownedcard{false, true, card})
		}
	}
	state.fillLot()
	return state
}

func (state state) String() string {
	s := fmt.Sprintf("Deck: %d Lot: %s %s %s\n", len(state.deck), state.lot[0], state.lot[1], state.lot[2])
	for i, p := range state.players {
		s += fmt.Sprintf("Player %d: %s", i, p)
	}
	return s
}

func (grid grid) valid(card ownedcard, row, column int) bool {
	// Variation: enforce each column has a different color
	// Variation: enforce that each column needs a starting card
	if row < 0 || column < 0 || row >= 3 || column >= 3 {
		return false
	}
	if !grid[row][column].card.empty() {
		return false
	}
	switch card.card.column {
	case ones:
		if column != 2 {
			return false
		}
	case tens:
		if column != 1 {
			return false
		}
	case hundreds:
		if column != 0 {
			return false
		}
	}
	if card.card.color != any {
		for i := 0; i < 3; i++ {
			if !grid[i][column].card.empty() && grid[i][column].card.color != any && grid[i][column].card.color != card.card.color {
				return false
			}
		}
	}
	grid[row][column] = card
	top := grid[0]
	mid := grid[1]
	bot := grid[2]
	for topn := 0; topn < 1000; topn++ {
		if !gridRowMatches(top, topn) {
			continue
		}
		for midn := 0; midn < 1000; midn++ {
			if !gridRowMatches(mid, midn) {
				continue
			}
			botn := (1000 + 2*midn - topn) % 1000
			if gridRowMatches(bot, botn) {
				return true
			}
		}
	}
	return false
}

func gridRowMatches(gridRow [3]ownedcard, number int) bool {
	return gridDigitMatches(gridRow[0], (number/100)%10) && gridDigitMatches(gridRow[1], (number/10)%10) && gridDigitMatches(gridRow[2], number%10)
}

func gridDigitMatches(gridDigit ownedcard, digit int) bool {
	return gridDigit.card.empty() || gridDigit.card.digit == any || int(gridDigit.card.digit) == digit
}

func (grid grid) completed() bool {
	for _, row := range grid {
		for _, digit := range row {
			if !digit.revealed || digit.card.empty() {
				return false
			}
		}
	}
	return true
}

func (state *state) completed() bool {
	for _, player := range state.players {
		if player.grid.completed() {
			return true
		}
	}
	for _, card := range state.lot {
		if card.empty() {
			return true
		}
	}
	return false
}

func (state *state) placePayouts() []int {
	// Variation: only opponent's empty spots pay out
	payout := 0
	for _, player := range state.players {
		for _, row := range player.grid {
			for _, digit := range row {
				if digit.card.empty() {
					payout++
				}
			}
		}
	}
	payouts := make([]int, len(state.players))
	for i := range payouts {
		payouts[i] = payout
	}
	return payouts
}

func (state *state) revealPayouts() []int {
	// Variation: only opponent's unrevealed pay out
	// Variation: payout multiplier (e.g. each unrevealed pays out 3)
	payout := 0
	for _, player := range state.players {
		for _, row := range player.grid {
			for _, card := range row {
				if !card.card.empty() && !card.revealed {
					payout++
				}
			}
		}
		for _, card := range player.cards {
			if !card.card.empty() && !card.revealed {
				payout++
			}
		}
	}
	payouts := make([]int, len(state.players))
	for i := range payouts {
		payouts[i] = payout
	}
	return payouts
}

func (state *state) fillLot() {
	for i := range state.lot {
		if len(state.deck) == 0 {
			break
		}
		if state.lot[i].empty() {
			state.lot[i] = state.deck[len(state.deck)-1]
			state.deck = state.deck[:len(state.deck)-1]
		}
	}
}
