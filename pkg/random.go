package pkg

import (
	"fmt"
	"math/rand"
)

func RandomGameID() string {
	return fmt.Sprintf("match-%d", rand.Int63())
}

func RandomFirstPlayer() Player {
	if rand.Intn(1) == 0 {
		return PlayerSelf
	} else {
		return PlayerOpponent
	}
}
