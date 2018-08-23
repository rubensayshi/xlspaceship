package ssclient

import (
	"net/http"

	"math/rand"

	"io"

	"github.com/pkg/errors"
	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
)

type XLSpaceship struct {
	Player    *ssgame.Player
	games     map[string]*ssgame.Game
	requester Requester
	cheat     bool
}

func NewXLSpaceship(playerID string, playerName string, host string, port int) *XLSpaceship {
	s := &XLSpaceship{
		Player: &ssgame.Player{
			PlayerID:     playerID,
			FullName:     playerName,
			ProtocolHost: host,
			ProtocolPort: port,
		},
		games:     make(map[string]*ssgame.Game),
		requester: &HttpRequester{},
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

func (s *XLSpaceship) EnableCheatMode() {
	s.cheat = true
}

func (s *XLSpaceship) NewGame(req *NewGameRequest) (*NewGameResponse, error) {
	opponent := &ssgame.Player{
		PlayerID:     req.UserID,
		FullName:     req.FullName,
		ProtocolHost: req.SpaceshipProtocol.Hostname,
		ProtocolPort: req.SpaceshipProtocol.Port,
	}

	game, err := ssgame.CreateNewGame(opponent)
	if err != nil {
		return nil, err
	}

	s.games[game.GameID] = game

	res := NewGameResponseFromGame(s, game)
	return res, nil
}

func (s *XLSpaceship) InitNewGame(req *InitGameRequest) (string, error) {
	newGameReq := &NewGameRequest{
		UserID:            s.Player.PlayerID,
		FullName:          s.Player.FullName,
		SpaceshipProtocol: SpaceshipProtocol{s.Player.ProtocolHost, s.Player.ProtocolPort},
	}

	newGameRes, err := s.requester.NewGame(req.SpaceshipProtocol, newGameReq)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to init new game")
	}

	firstPlayer := ssgame.PlayerSelf
	if newGameRes.Starting != s.Player.PlayerID {
		firstPlayer = ssgame.PlayerOpponent
	}

	opponent := &ssgame.Player{
		PlayerID:     newGameRes.UserID,
		FullName:     newGameRes.FullName,
		ProtocolHost: req.SpaceshipProtocol.Hostname,
		ProtocolPort: req.SpaceshipProtocol.Port,
	}

	game, err := ssgame.InitNewGame(newGameRes.GameID, opponent, firstPlayer)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to init new game")
	}

	s.games[game.GameID] = game

	return game.GameID, nil
}

func (s *XLSpaceship) GameStatus(gameID string) (*GameStatusResponse, bool) {
	game, ok := s.games[gameID]
	if !ok {
		return nil, false
	}

	res := GameStatusResponseFromGame(s, game)

	return res, true
}

func (s *XLSpaceship) FireSalvo(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, error) {
	// check that we're not cheating
	if !s.cheat && len(salvo) > game.SelfBoard.CountShipsAlive() {
		return nil, errors.Errorf("More shots than ships alive (%d)", game.SelfBoard.CountShipsAlive())
	}

	req := &ReceiveSalvoRequest{
		Salvo: make([]string, len(salvo)),
	}

	for i, salvo := range salvo {
		req.Salvo[i] = salvo.String()
	}

	res, err := s.requester.ReceiveSalvo(SpaceshipProtocol{
		Hostname: game.Opponent.ProtocolHost,
		Port:     game.Opponent.ProtocolPort,
	}, game.GameID, req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fire salvo")
	}

	salvoRes := make([]*ssgame.ShotResult, 0, len(res.Salvo))
	for coordsStr, shotResStr := range res.Salvo {
		coords, err := ssgame.CoordsFromString(coordsStr)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to fire salvo")
		}

		shotStatus, err := ssgame.ShotStatusFromString(shotResStr)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to fire salvo")
		}

		salvoRes = append(salvoRes, &ssgame.ShotResult{coords, shotStatus})

		game.OpponentBoard.ApplyShotStatus(coords, shotStatus)
	}

	game.PlayerTurn = ssgame.PlayerOpponent

	// @TODO
	win := false
	if win {
		game.Status = ssgame.GameStatusDone
		game.PlayerWon = ssgame.PlayerOpponent
	}

	// @TODO: mark result on our end

	return ReceiveSalvoResponseFromSalvoResult(salvoRes, s, game), nil
}

func (s *XLSpaceship) ReceiveSalvo(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, error) {
	// @TODO: we need to assert that the amount of shots in the salvo match the rules (1 per ship alive), but for that we need track when we kill enemy ships

	salvoRes := game.SelfBoard.ReceiveSalvo(salvo)
	game.PlayerTurn = ssgame.PlayerSelf

	if game.SelfBoard.AllShipsDead() {
		game.Status = ssgame.GameStatusDone
		game.PlayerWon = ssgame.PlayerOpponent
	}

	return ReceiveSalvoResponseFromSalvoResult(salvoRes, s, game), nil
}

// Helper function to do PUT requests because http builtin only has helpers for GET and POST >_>
func Put(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}
