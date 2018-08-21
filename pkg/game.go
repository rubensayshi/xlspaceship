package pkg

import (
	"fmt"
	"math/rand"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

type GameStatus uint8

const (
	GameStatusInitializing GameStatus = 0
	GameStatusOnGoing      GameStatus = 1
	GameStatusDone         GameStatus = 2
)

type Player uint8

const (
	PlayerNone     Player = 0
	PlayerSelf     Player = 1
	PlayerOpponent Player = 2
)

type Game struct {
	GameID           string
	OpponentPlayerID string
	Status           GameStatus
	SelfBoard        *Board
	OpponentBoard    *Board
	PlayerTurn       Player
	PlayerWon        Player
}

func NewGame(opponentPlayerID string) (*Game, error) {
	selfBoard, err := NewRandomBoard()
	if err != nil {
		return nil, err
	}

	opponentBoard, err := BoardFromPattern(BlankBoardPattern())
	if err != nil {
		return nil, err
	}

	firstPlayer := RandomFirstPlayer()

	game := &Game{
		GameID:           RandomGameID(),
		OpponentPlayerID: opponentPlayerID,
		Status:           GameStatusInitializing,
		SelfBoard:        selfBoard,
		OpponentBoard:    opponentBoard,
		PlayerTurn:       firstPlayer,
		PlayerWon:        PlayerNone,
	}

	// start the game
	game.Status = GameStatusOnGoing

	return game, nil
}

func (g *Game) String() string {
	return fmt.Sprintf(
		"opponent: %s\n"+
			"self-board: \n%s\n"+
			"opponent-board: \n%s\n",
		g.OpponentPlayerID, g.SelfBoard, g.OpponentBoard)
}
