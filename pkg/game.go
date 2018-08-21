package pkg

import (
	"math/rand"
)

type GameStatus int8

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

const (
	GameStatusInitializing GameStatus = 0
	GameStatusOnGoing      GameStatus = 1
	GameStatusDone         GameStatus = 2
)

type Player int8

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

type Board struct {
	spaceships []*Spaceship
	hits       []int8
	misses     []int8
}

func BoardFromPattern(pattern string) (*Board, error) {
	board := &Board{}

	return board, nil
}

type Spaceship struct {
	coords []int8
}

func SpaceshipFromPattern(pattern string) (*Spaceship, error) {
	spaceship := &Spaceship{}

	return spaceship, nil
}
