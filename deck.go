package main

import (
	"fmt"
	"math/rand"
)

type card struct {
	digit  digit
	color  color
	column column
}

type digit int
type color int
type column int

const (
	invalid = iota

	red
	green
	blue

	ones
	tens
	hundreds

	any = -1
)

var (
	digits  = []digit{any, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	colors  = []color{any, red, green, blue}
	columns = []column{any, ones, tens, hundreds}
)

func (digit digit) String() string {
	if digit == any {
		return "*"
	} else {
		return fmt.Sprintf("%d", digit)
	}
}

func (color color) String() string {
	switch color {
	case red:
		return "R"
	case green:
		return "G"
	case blue:
		return "B"
	default:
		return "*"
	}
}

func (column column) String() string {
	switch column {
	case ones:
		return "[--*]"
	case tens:
		return "[-*-]"
	case hundreds:
		return "[*--]"
	default:
		return "[***]"
	}
}

func (card card) String() string {
	if card.empty() {
		return "-------"
	}
	return card.color.String() + card.column.String() + card.digit.String()
}

func (card card) empty() bool {
	return card.color == invalid
}

func deck() []card {
	var deck []card
	for _, digit := range digits {
		if digit == any {
			continue
		}
		deck = append(deck, card{digit, any, any})
		for _, color := range colors {
			if color == any {
				continue
			}
			deck = append(deck, card{digit, color, any})
			deck = append(deck, card{digit, color, any})
		}
		for _, column := range columns {
			if column == any {
				continue
			}
			deck = append(deck, card{digit, any, column})
		}
	}
	for _, color := range colors {
		if color == any {
			continue
		}
		for _, column := range columns {
			if column == any {
				continue
			}
			deck = append(deck, card{any, color, column})
		}
	}
	return deck
}

func shuffleAndDeal(nplayers int) ([]card, [][3]card) {
	deck := deck()
	var cards []card
	players := make([][3]card, nplayers)
nextCard:
	for len(deck) > 0 {
		index := rand.Intn(len(deck))
		card := deck[index]
		copy(deck[index:], deck[index+1:])
		deck = deck[:len(deck)-1]
		if card.color == red && card.column == any {
			for i := range players {
				if players[i][0].color != red {
					players[i][0] = card
					continue nextCard
				}
			}
		} else if card.color == blue && card.column == any {
			for i := range players {
				if players[i][1].color != blue {
					players[i][1] = card
					continue nextCard
				}
			}
		} else if card.color == green && card.column == any {
			for i := range players {
				if players[i][2].color != green {
					players[i][2] = card
					continue nextCard
				}
			}
		}
		cards = append(cards, card)
	}
	return cards, players
}
