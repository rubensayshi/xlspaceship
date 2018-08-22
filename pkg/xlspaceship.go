package pkg

import (
	"net/http"

	"bytes"
	"encoding/json"
	"fmt"

	"math/rand"

	"io"

	"github.com/pkg/errors"
)

type XLSpaceship struct {
	Player *Player
	games  map[string]*Game
}

func NewXLSpaceship(playerID string, host string, port int) *XLSpaceship {
	s := &XLSpaceship{
		Player: &Player{
			PlayerID:     playerID,
			FullName:     playerID,
			ProtocolHost: host,
			ProtocolPort: port,
		},
		games: make(map[string]*Game),
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
	opponent := &Player{
		PlayerID:     opponentPlayerID,
		FullName:     opponentName,
		ProtocolHost: opponentHost,
		ProtocolPort: opponentPort,
	}

	game, err := CreateNewGame(opponent)
	if err != nil {
		return nil, err
	}

	s.games[game.GameID] = game

	return game, nil
}

func (s *XLSpaceship) InitNewGame(opponentHost string, opponentPort int) (*Game, error) {
	req := NewGameRequest{
		UserID:            s.Player.PlayerID,
		FullName:          s.Player.FullName,
		SpaceshipProtocol: GameRequestSpaceshipProtocol{s.Player.ProtocolHost, s.Player.ProtocolPort},
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
	if newGameRes.Starting != s.Player.PlayerID {
		firstPlayer = PlayerOpponent
	}

	opponent := &Player{
		PlayerID:     newGameRes.UserID,
		FullName:     newGameRes.FullName,
		ProtocolHost: opponentHost,
		ProtocolPort: opponentPort,
	}

	game, err := InitNewGame(newGameRes.GameID, opponent, firstPlayer)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to init game")
	}

	s.games[game.GameID] = game

	return game, nil
}

func (s *XLSpaceship) FireSalvo(game *Game, salvo CoordsGroup) ([]*ShotResult, error) {
	req := ReceiveSalvoRequest{
		Salvo: make([]string, len(salvo)),
	}

	for i, salvo := range salvo {
		req.Salvo[i] = salvo.String()
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fire salvo")
	}

	res, err := Put(fmt.Sprintf("http://%s:%d/xl-spaceship/protocol/game/%s", game.Opponent.ProtocolHost, game.Opponent.ProtocolPort, game.GameID), "application/json", bytes.NewBuffer(reqJson))
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "Failed to fire salvo")
	}
	defer res.Body.Close()

	newGameRes := &ReceiveSalvoResponse{}
	err = json.NewDecoder(res.Body).Decode(newGameRes)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fire salvo")
	}

	shotsRes := make([]*ShotResult, 0, len(newGameRes.Salvo))
	for coordsStr, shotResStr := range newGameRes.Salvo {
		coords, err := CoordsFromString(coordsStr)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to fire salvo")
		}

		shotStatus, err := ShotStatusFromString(shotResStr)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to fire salvo")
		}

		shotsRes = append(shotsRes, &ShotResult{coords, shotStatus})
	}

	game.PlayerTurn = PlayerOpponent

	return shotsRes, nil
}

func Put(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}
