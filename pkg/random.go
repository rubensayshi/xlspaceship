package pkg

import (
	"fmt"
	"math/rand"
)

const UUID = 1

var gameID = 0

// Instead of a random gameID we just increment a number and prefix with a UUID
// The UUID should be random, but for debugging it's so much easier if it's not
func RandomGameID() string {
	gameID++
	return fmt.Sprintf("match-%d-%d", UUID, gameID)
}

func RandomFirstPlayer() Player {
	if rand.Intn(1) == 0 {
		return PlayerSelf
	} else {
		return PlayerOpponent
	}
}
