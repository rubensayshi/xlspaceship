package pkg

import (
	"net/http"

	"bytes"
	"encoding/json"
	"fmt"

	"math/rand"

	"github.com/pkg/errors"
)

type XLSpaceship struct {
	PlayerID   string
	PlayerName string
	games      map[string]*Game
}

func NewXLSpaceship(playerID string) *XLSpaceship {
	s := &XLSpaceship{
		PlayerID:   playerID,
		PlayerName: playerID,
		games:      make(map[string]*Game),
	}

	// make a random seed based on the playerID, that way it's deterministic but different per player
	//  this is purely for easy debugging
	var seed int64 = 0
	for _, char := range playerID {
		seed += int64(char)
	}
	rand.Seed(seed)

	return s
}

func (s *XLSpaceship) NewGame(opponentPlayerID string, opponentName string, opponentHost string, opponentPort int) (*Game, error) {
	game, err := CreateNewGame(opponentPlayerID)
	if err != nil {
		return nil, err
	}

	s.games[game.GameID] = game

	return game, nil
}

func (s *XLSpaceship) InitNewGame(opponentHost string, opponentPort int) (*Game, error) {
	req := NewGameRequest{
		UserID:            s.PlayerID,
		FullName:          s.PlayerName,
		SpaceshipProtocol: GameRequestSpaceshipProtocol{},
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to init game")
	}

	res, err := http.Post(fmt.Sprintf("http://%s:%d/xl-spaceship/protocol/game/new", opponentHost, opponentPort), "application/json", bytes.NewBuffer(reqJson))
	if err != nil || res.StatusCode != http.StatusCreated {
		return nil, errors.Wrapf(err, "Failed to init game")
	}

	defer res.Body.Close()

	newGameRes := &NewGameResponse{}
	err = json.NewDecoder(res.Body).Decode(newGameRes)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to init game")
	}

	firstPlayer := PlayerSelf
	if newGameRes.Starting != s.PlayerID {
		firstPlayer = PlayerOpponent
	}

	game, err := InitNewGame(newGameRes.GameID, newGameRes.UserID, firstPlayer)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to init game")
	}

	s.games[game.GameID] = game

	return game, nil
}
