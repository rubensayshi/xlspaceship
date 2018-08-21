package pkg

import (
	"net/http"

	"fmt"
)

type XLSpaceship struct {
	PlayerID   string
	PlayerName string
	games      []*Game
}

func NewXLSpaceship() *XLSpaceship {
	s := &XLSpaceship{
		PlayerID:   "roobs-1",
		PlayerName: "Roobs",
	}

	return s
}

func (s *XLSpaceship) NewGame(opponentPlayerID string, opponentName string, opponentHost string, opponentPort int) *Game {
	// @TODO: handle this nicely?
	if !PingOpponent(opponentHost, opponentPort) {
		fmt.Printf("failed to ping opponent \n")
	}

	game := NewGame(opponentPlayerID)

	s.games = append(s.games, game)

	return game
}

func PingOpponent(opponentHost string, opponentPort int) bool {
	res, err := http.Get(fmt.Sprintf("http://%s:%d/xl-spaceship/ping", opponentHost, opponentPort))
	if err != nil || res.StatusCode != http.StatusOK {
		return false
	}

	return true
}
