package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func auction(state *state, in *bufio.Reader) {
	droppedOut := make([]bool, len(state.players))
	nDroppedOut := 0
	currentBid := -1
	currentWinner := -1
	currentPlayer := state.currentPlayer
	for {
		if nDroppedOut >= len(state.players) {
			break
		}
		if droppedOut[currentPlayer] {
			currentPlayer = (currentPlayer + 1) % len(state.players)
			continue
		}
		if currentPlayer == currentWinner {
			state.players[currentPlayer].cash -= currentBid
			state.currentPlayer = currentPlayer
			wonAuction(state, in)
			break
		}
		if currentBid >= state.players[currentPlayer].cash {
			droppedOut[currentPlayer] = true
			nDroppedOut++
			currentPlayer = (currentPlayer + 1) % len(state.players)
			continue
		}
		if currentBid >= 0 {
			fmt.Printf("[Current bid: Player %d: $%d] ", currentWinner, currentBid)
		}
		fmt.Printf("Player %d bid (minimum $%d) or drop:\n", currentPlayer, currentBid+1)
		line, err := in.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}
		line = strings.TrimSpace(line)
		if line == "drop" {
			droppedOut[currentPlayer] = true
			nDroppedOut++
			currentPlayer = (currentPlayer + 1) % len(state.players)
			continue
		}
		bid := 0
		if n, err := fmt.Sscanf(line, "%d", &bid); n != 1 || err != nil {
			continue
		}
		if bid <= currentBid || bid > state.players[currentPlayer].cash {
			continue
		}
		currentBid = bid
		currentWinner = currentPlayer
		currentPlayer = (currentPlayer + 1) % len(state.players)
	}
}

func wonAuction(state *state, in *bufio.Reader) {
	for {
		fmt.Printf("Player %d: choose card to take (1:%s 2:%s 3:%s):\n", state.currentPlayer, state.lot[0], state.lot[1], state.lot[2])
		line, err := in.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}
		line = strings.TrimSpace(line)
		n := -1
		switch line {
		case "1":
			n = 0
		case "2":
			n = 1
		case "3":
			n = 2
		default:
			continue
		}
		state.players[state.currentPlayer].cards = append(state.players[state.currentPlayer].cards, ownedcard{true, false, state.lot[n]})
		state.lot[n] = card{}
		break
	}
	for i := range state.lot {
		if state.lot[i].empty() {
			continue
		}
		for {
			fmt.Printf("Player %d: kill or leave card %s:\n", state.currentPlayer, state.lot[i])
			line, err := in.ReadString('\n')
			if err != nil {
				os.Exit(0)
			}
			line = strings.TrimSpace(line)
			switch line {
			case "kill", "k":
				state.lot[i] = card{}
			case "leave", "l":
			default:
				continue
			}
			break
		}
	}
	state.fillLot()
}

func payout(state *state, in *bufio.Reader) {
	placePayouts := state.placePayouts()
	revealPayouts := state.revealPayouts()
	for i := range state.players {
		fmt.Printf("%s", state)
		currentPlayer := (len(state.players) + state.currentPlayer - i) % len(state.players)
		nPlaced := 0
		nRevealed := 0
		unrevealed := []*ownedcard{}
		for j, card := range state.players[currentPlayer].cards {
			if !card.revealed {
				unrevealed = append(unrevealed, &state.players[currentPlayer].cards[j])
			}
		}
		for j := range state.players[currentPlayer].grid {
			for k, card := range state.players[currentPlayer].grid[j] {
				if !card.card.empty() && !card.revealed {
					unrevealed = append(unrevealed, &state.players[currentPlayer].grid[j][k])
				}
			}
		}
		for {
			fmt.Printf("Player %d (place payout: %d, reveal payout: %d), place or reveal or done [", currentPlayer, placePayouts[currentPlayer], revealPayouts[currentPlayer])
			for i, card := range state.players[currentPlayer].cards {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%d:%s", i+1, card)
			}
			fmt.Printf("]:\n")
			line, err := in.ReadString('\n')
			if err != nil {
				os.Exit(0)
			}
			line = strings.TrimSpace(line)
			if line == "done" {
				break
			}
			var i, row, col int
			if n, err := fmt.Sscanf(line, "place %d at %d %d", &i, &row, &col); n == 3 && err == nil {
				if i < 1 || i > len(state.players[currentPlayer].cards) || !state.players[currentPlayer].grid.valid(state.players[currentPlayer].cards[i-1], row, col) {
					continue
				}
				nPlaced++
				state.players[currentPlayer].grid[row][col] = state.players[currentPlayer].cards[i-1]
				copy(state.players[currentPlayer].cards[i-1:], state.players[currentPlayer].cards[i:])
				state.players[currentPlayer].cards = state.players[currentPlayer].cards[:len(state.players[currentPlayer].cards)-1]
				fmt.Printf("%s", state.players[currentPlayer].grid)
				continue
			}
			var c rune
			if n, err := fmt.Sscanf(line, "reveal %c", &c); n == 1 && err == nil {
				var color color
				switch c {
				case 'r', 'R':
					color = red
				case 'g', 'G':
					color = green
				case 'b', 'B':
					color = blue
				default:
					continue
				}
				for _, c := range unrevealed {
					if !c.revealed && c.card.color == color {
						c.revealed = true
						nRevealed++
						break
					}
				}
				continue
			}
		}
		state.players[currentPlayer].cash += nPlaced*placePayouts[currentPlayer] + nRevealed*revealPayouts[currentPlayer]
	}
}

func playGame(state *state) {
	in := bufio.NewReader(os.Stdin)
	for !state.completed() {
		fmt.Printf("%s", state)
		auction(state, in)
		payout(state, in)
	}
}
