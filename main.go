package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	state := setup(3)
	fmt.Printf("%s", state)
	playGame(state)
}
