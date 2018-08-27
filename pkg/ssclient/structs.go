package ssclient

import (
	"github.com/pkg/errors"
	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
)

type SpaceshipProtocol struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type WhoAmIRequest struct {
}

type WhoAmIResponse struct {
	UserID   string   `json:"user_id"`
	FullName string   `json:"full_name"`
	Games    []string `json:"games"`
}

type NewGameRequest struct {
	UserID            string            `json:"user_id"`
	FullName          string            `json:"full_name"`
	SpaceshipProtocol SpaceshipProtocol `json:"spaceship_protocol"`
}

type NewGameResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	GameID   string `json:"game_id"`
	Starting string `json:"starting"`
}

func NewGameResponseFromGame(s *XLSpaceship, game *ssgame.Game) *NewGameResponse {
	res := &NewGameResponse{}

	res.UserID = s.Player.PlayerID
	res.FullName = s.Player.FullName
	res.GameID = game.GameID

	if game.PlayerTurn == ssgame.PlayerSelf {
		res.Starting = s.Player.PlayerID
	} else {
		res.Starting = game.Opponent.PlayerID
	}

	return res
}

type InitGameRequest struct {
	SpaceshipProtocol SpaceshipProtocol `json:"spaceship_protocol"`
}

type GameStatusRequest struct {
	GameID string `json:"game_id"`
}

type GamePlayerTurnResponse struct {
	PlayerTurn string `json:"player_turn"`
}

type GameWonResponse struct {
	Won string `json:"won"`
}

type GameStatusResponse struct {
	GameID   string                   `json:"game_id"`
	Self     GameStatusResponsePlayer `json:"self"`
	Opponent GameStatusResponsePlayer `json:"opponent"`
	Game     interface{}              `json:"game"`
}

type GameStatusResponsePlayer struct {
	UserID string   `json:"user_id"`
	Board  []string `json:"board"`
	Shots  int      `json:"shots"`
}

func GameStatusResponseFromGame(s *XLSpaceship, game *ssgame.Game) *GameStatusResponse {
	res := &GameStatusResponse{
		GameID: game.GameID,
	}

	res.Self = GameStatusResponsePlayer{
		UserID: s.Player.PlayerID,
		Board:  game.SelfBoard.ToPattern(),
		Shots:  game.SelfBoard.CountShipsAlive(),
	}

	res.Opponent = GameStatusResponsePlayer{
		UserID: game.Opponent.PlayerID,
		Board:  game.OpponentBoard.ToPattern(),
		Shots:  game.OpponentBoard.CountShipsAlive(),
	}

	if game.Status == ssgame.GameStatusDone {
		won := s.Player.PlayerID
		if game.PlayerWon == ssgame.PlayerOpponent {
			won = game.Opponent.PlayerID
		}

		res.Game = GameWonResponse{
			Won: won,
		}
	} else {
		playerTurn := s.Player.PlayerID
		if game.PlayerTurn == ssgame.PlayerOpponent {
			playerTurn = game.Opponent.PlayerID
		}

		res.Game = GamePlayerTurnResponse{
			PlayerTurn: playerTurn,
		}
	}

	return res
}

type FireSalvoRequest struct {
	GameID string   `json:"-"`
	Salvo  []string `json:"salvo"`
}

type ReceiveSalvoRequest struct {
	GameID string   `json:"-"`
	Salvo  []string `json:"salvo"`
}

// @TODO: should make custom JSON marshall/unmarshall for the "game" field instead of the hacky way we do now
type SalvoResponse struct {
	Salvo           map[string]string       `json:"salvo"`
	Game            map[string]string       `json:"game"`
	GameWon         *GameWonResponse        `json:"-"`
	GamePlayerTurn  *GamePlayerTurnResponse `json:"-"`
	AlreadyFinished bool                    `json:"-"`
}

func SalvoResponseFromSalvoResult(salvoResult []*ssgame.ShotResult, xl *XLSpaceship, game *ssgame.Game) *SalvoResponse {
	res := &SalvoResponse{
		Salvo: make(map[string]string, len(salvoResult)),
	}

	for _, shotResult := range salvoResult {
		res.Salvo[shotResult.Coords.String()] = shotResult.ShotStatus.String()
	}

	if game.Status == ssgame.GameStatusDone {
		res.GameWon = &GameWonResponse{Won: xl.Player.PlayerID}
		if game.PlayerWon != ssgame.PlayerSelf {
			res.GameWon.Won = game.Opponent.PlayerID
		}
		res.Game = map[string]string{"won": res.GameWon.Won}
	} else {
		res.GamePlayerTurn = &GamePlayerTurnResponse{PlayerTurn: xl.Player.PlayerID}
		if game.PlayerTurn != ssgame.PlayerSelf {
			res.GamePlayerTurn.PlayerTurn = game.Opponent.PlayerID
		}

		res.Game = map[string]string{"player_turn": res.GamePlayerTurn.PlayerTurn}
	}

	return res
}

func (r *SalvoResponse) Normalize() error {
	if playerIdWon, won := r.Game["won"]; won {
		r.GameWon = &GameWonResponse{
			Won: playerIdWon,
		}
	} else if playerIdTurn, turn := r.Game["player_turn"]; turn {
		r.GamePlayerTurn = &GamePlayerTurnResponse{
			PlayerTurn: playerIdTurn,
		}
	} else {
		return errors.Errorf("SalvoResponse should either contain 'won' or 'player_turn'")
	}

	return nil
}
