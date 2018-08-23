package ssgame

import (
	"fmt"
	"math/rand"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

type GameStatus int8

const (
	GameStatusOnGoing GameStatus = 0
	GameStatusDone    GameStatus = 1
)

type PlayerT int8

func (p PlayerT) String() string {
	switch p {
	case PlayerNone:
		return "none"
	case PlayerSelf:
		return "self"
	case PlayerOpponent:
		return "opponent"
	}

	panic("Unreachable")
}

const (
	PlayerNone     PlayerT = 0
	PlayerSelf     PlayerT = 1
	PlayerOpponent PlayerT = 2
)

type Player struct {
	PlayerID     string
	FullName     string
	ProtocolHost string
	ProtocolPort int
}

type Game struct {
	GameID        string
	Opponent      *Player
	Status        GameStatus
	SelfBoard     *SelfBoard
	OpponentBoard *OpponentBoard
	PlayerTurn    PlayerT
	PlayerWon     PlayerT
}

func CreateNewGame(opponent *Player) (*Game, error) {
	selfBoard, err := NewRandomSelfBoard(SpaceshipsSetForBaseGame)
	if err != nil {
		return nil, err
	}

	opponentBoard, err := NewBlankOpponentBoard()
	if err != nil {
		return nil, err
	}

	firstPlayer := RandomFirstPlayer()

	game := &Game{
		GameID:        RandomGameID(),
		Opponent:      opponent,
		Status:        GameStatusOnGoing,
		SelfBoard:     selfBoard,
		OpponentBoard: opponentBoard,
		PlayerTurn:    firstPlayer,
		PlayerWon:     PlayerNone,
	}

	return game, nil
}

func InitNewGame(gameID string, opponent *Player, firstPlayer PlayerT) (*Game, error) {
	selfBoard, err := NewRandomSelfBoard(SpaceshipsSetForBaseGame)
	if err != nil {
		return nil, err
	}

	opponentBoard, err := NewBlankOpponentBoard()
	if err != nil {
		return nil, err
	}

	game := &Game{
		GameID:        gameID,
		Opponent:      opponent,
		Status:        GameStatusOnGoing,
		SelfBoard:     selfBoard,
		OpponentBoard: opponentBoard,
		PlayerTurn:    firstPlayer,
		PlayerWon:     PlayerNone,
	}

	return game, nil
}

func (g *Game) String() string {
	return fmt.Sprintf(
		"opponent: %s\n"+
			"player-turn: %s\n"+
			"self-board: \n%s\n"+
			"opponent-board: \n%s\n",
		g.Opponent.PlayerID, g.PlayerTurn, g.SelfBoard, g.OpponentBoard)
}
