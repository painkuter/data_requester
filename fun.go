package main

import (
	"math/rand"
	"time"
)

func generateDelay() int {
	// generating random int name
	now := time.Now().UnixNano()
	source := rand.NewSource(now)
	randomizer := rand.New(source)
	return randomizer.Intn(400)
}
