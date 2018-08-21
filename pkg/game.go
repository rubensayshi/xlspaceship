package pkg

import (
	"errors"
	"fmt"
	"math/rand"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

const ROWS = 16
const COLS = 16

type CoordState byte

func (c CoordState) String() string {
	switch c {
	case CoordBlank:
		return "."
	case CoordShip:
		return "*"
	case CoordHit:
		return "X"
	case CoordMiss:
		return "-"
	}

	panic("Unreachable")
}

const (
	CoordBlank CoordState = '.'
	CoordShip  CoordState = '*'
	CoordHit   CoordState = 'X'
	CoordMiss  CoordState = '-'
)

type Coord struct {
	x int8
	y int8
}

func (c Coord) String() string {
	return fmt.Sprintf("%Xx%X", c.x, c.y)
}

type GameStatus int8

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
	hits       []Coord
	misses     []Coord
}

func BoardFromPattern(pattern []string) (*Board, error) {
	// sanity check the input
	if len(pattern) != ROWS {
		return nil, errors.New("pattern incorrect amount of rows")
	}

	// sanity check the input
	for _, row := range pattern {
		if len(row) != COLS {
			return nil, errors.New("pattern incorrect amount of cols")
		}

		// @TODO: is there a nicer way to do this with a builtin?
		for _, char := range []byte(row) {
			if char != byte(CoordBlank) && char != byte(CoordShip) && char != byte(CoordHit) && char != byte(CoordMiss) {
				return nil, errors.New("pattern incorrect symbol for coord")
			}
		}
	}

	board := &Board{}

	// parse the input
	for y, row := range pattern {
		for x, char := range []byte(row) {
			coordState := CoordState(char)

			switch coordState {
			case CoordBlank:
				// - nothing to do
			case CoordShip:
				// @TODO: not implemented
			case CoordHit:
				board.hits = append(board.hits, Coord{x: int8(x), y: int8(y)})
			case CoordMiss:
				board.misses = append(board.misses, Coord{x: int8(x), y: int8(y)})
			}
		}
	}

	return board, nil
}

type Spaceship struct {
	coords []int8
}

func SpaceshipFromPattern(pattern []string) (*Spaceship, error) {
	spaceship := &Spaceship{}

	return spaceship, nil
}
