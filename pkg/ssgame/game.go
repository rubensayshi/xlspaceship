package ssgame

import (
	"fmt"
	"math/rand"
)

func init() {
	// give us a "random seed"
	rand.Seed(1)

	// this would be the proper way
	//  but for serious things (eg a casino game, crypto, etc) crypto/rand should be used instead of math/rand
	// rand.Seed(time.Now().UTC().UnixNano())
}

// define the status a game can have as a type
type GameStatus int8

const (
	GameStatusOnGoing GameStatus = 0
	GameStatusDone    GameStatus = 1
)

// the type we use to indicate which player's turn it is or which player won etc
type WhichPlayer int8

func (p WhichPlayer) String() string {
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
	PlayerNone     WhichPlayer = 0
	PlayerSelf     WhichPlayer = 1
	PlayerOpponent WhichPlayer = 2
)

// the type we use to store a player (both self and opponent)
type Player struct {
	PlayerID     string
	FullName     string
	ProtocolHost string
	ProtocolPort int
}

// the type to hold a game between 2 players
type Game struct {
	GameID        string
	Opponent      *Player
	Status        GameStatus
	SelfBoard     *SelfBoard
	OpponentBoard *OpponentBoard
	PlayerTurn    WhichPlayer
	PlayerWon     WhichPlayer
}

// create a new game with a random board for self and a blank board for opponent
func CreateNewGame(gameID string, opponent *Player, cheatToBeFirst bool) (*Game, error) {
	// give ourselves a random board
	selfBoard, err := NewRandomSelfBoard(SpaceshipsSetForBaseGame)
	if err != nil {
		return nil, err
	}

	// give our opponent a blank board
	opponentBoard, err := NewBlankOpponentBoard(uint8(len(SpaceshipsSetForBaseGame)))
	if err != nil {
		return nil, err
	}

	// determine which player get's to go first
	firstPlayer := PlayerSelf
	if !cheatToBeFirst {
		firstPlayer = RandomFirstPlayer()
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

// init a new game that we were challanged to play
func InitNewGame(gameID string, opponent *Player, firstPlayer WhichPlayer) (*Game, error) {
	// give ourselves a random board
	selfBoard, err := NewRandomSelfBoard(SpaceshipsSetForBaseGame)
	if err != nil {
		return nil, err
	}

	// give our opponent a blank board
	opponentBoard, err := NewBlankOpponentBoard(uint8(len(SpaceshipsSetForBaseGame)))
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
