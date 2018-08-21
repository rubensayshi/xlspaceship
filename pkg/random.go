package pkg

import (
	"fmt"
	"math/rand"
)

func RandomGameID() string {
	return fmt.Sprintf("match-%d", rand.Int63())
}

func RandomFirstPlayer() Player {
	if rand.Int31()%2 == 0 {
		return PlayerSelf
	} else {
		return PlayerOpponent
	}
}
