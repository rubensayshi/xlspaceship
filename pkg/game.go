package pkg

import (
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

func NewGame(opponentPlayerID string) *Game {
	selfBoard := &Board{}
	opponentBoard := &Board{}

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

	return game
}

func NewRandomBoard() (*Board, error) {
	board := &Board{}

	for _, spaceshipPattern := range [][]string{
		SpaceshipPatternWinger,
		SpaceshipPatternAngle,
		SpaceshipPatternAClass,
		SpaceshipPatternBClass,
		SpaceshipPatternSClass,
	} {
		spaceship, err := SpaceshipFromPattern(spaceshipPattern)
		if err != nil {
			return nil, err
		}

		err = board.AddSpaceship(spaceship)
		if err != nil {
			return nil, err
		}
	}

	return board, nil
}
