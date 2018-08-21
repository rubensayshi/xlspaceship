package pkg

import (
	"net/http"

	"fmt"
)

type XLSpaceship struct {
	PlayerID   string
	PlayerName string
	games      map[string]*Game
}

func NewXLSpaceship() *XLSpaceship {
	s := &XLSpaceship{
		PlayerID:   "roobs-1",
		PlayerName: "Roobs",
		games:      make(map[string]*Game),
	}

	return s
}

func (s *XLSpaceship) NewGame(opponentPlayerID string, opponentName string, opponentHost string, opponentPort int) (*Game, error) {
	// @TODO: handle this nicely?
	if !PingOpponent(opponentHost, opponentPort) {
		fmt.Printf("failed to ping opponent \n")
	}

	game, err := NewGame(opponentPlayerID)
	if err != nil {
		return nil, err
	}

	s.games[game.GameID] = game

	return game, nil
}

func PingOpponent(opponentHost string, opponentPort int) bool {
	res, err := http.Get(fmt.Sprintf("http://%s:%d/xl-spaceship/ping", opponentHost, opponentPort))
	if err != nil || res.StatusCode != http.StatusOK {
		return false
	}

	return true
}
