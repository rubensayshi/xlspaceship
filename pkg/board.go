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
		return CoordBlankStr
	case CoordShip:
		return CoordShipStr
	case CoordHit:
		return CoordHitStr
	case CoordMiss:
		return CoordMissStr
	}

	panic("Unreachable")
}

const (
	CoordBlank CoordState = '.'
	CoordShip  CoordState = '*'
	CoordHit   CoordState = 'X'
	CoordMiss  CoordState = '-'

	CoordBlankStr string = "."
	CoordShipStr  string = "*"
	CoordHitStr   string = "X"
	CoordMissStr  string = "-"
)

type Coord struct {
	x uint8
	y uint8
}

func (c Coord) String() string {
	return fmt.Sprintf("%Xx%X", c.x, c.y)
}

type Board struct {
	spaceships []*Spaceship
	hits       []*Coord
	misses     []*Coord
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

	board := &Board{
		spaceships: make([]*Spaceship, 0),
		hits:       make([]*Coord, 0),
		misses:     make([]*Coord, 0),
	}

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
				board.hits = append(board.hits, &Coord{x: uint8(x), y: uint8(y)})
			case CoordMiss:
				board.misses = append(board.misses, &Coord{x: uint8(x), y: uint8(y)})
			}
		}
	}

	return board, nil
}

func (b *Board) ToPattern() []string {
	// @TODO: considering the board size is constant we could just have a const string to copy for this instead of building the blank state everytime
	pattern := make([][]byte, ROWS)
	for y, _ := range pattern {
		pattern[y] = make([]byte, COLS)

		for x := 0; x < COLS; x++ {
			pattern[y][x] = byte(CoordBlank)
		}
	}

	for _, hit := range b.hits {
		pattern[hit.y][hit.x] = byte(CoordHit)
	}

	for _, miss := range b.misses {
		pattern[miss.y][miss.x] = byte(CoordMiss)
	}

	res := make([]string, ROWS)
	for y, row := range pattern {
		res[y] = string(row)
	}

	return res
}

func (b *Board) AddSpaceship(spaceship *Spaceship) error {
	N := 10000
	// we'll attempt to add the spaceship on random locations until we succeed to reach N
	// @TODO: this could be heavily optimized as we know we don't have to try adding a spaceship of 3 high on Y > 15 - 3
	for i := 0; i < N; i++ {
		x := rand.Intn(COLS)
		y := rand.Intn(ROWS)

		err := b.AddSpaceshipOnCoords(spaceship, uint8(x), uint8(y))
		// we're done if no err
		if err == nil {
			return nil
		}
	}

	return errors.New("Failed to add spaceship, seems impossible")
}

func (b *Board) AddSpaceshipOnCoords(spaceship *Spaceship, x uint8, y uint8) error {
	// @TODO: we should store coords of existing spaceships so this isn't O(N2)
	for _, coord := range spaceship.coords {
		if coord.x+x >= ROWS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, x overflow (%d + %d = %d)", x, coord.x, coord.x+x))
		}
		if coord.y+y >= COLS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, y overflow (%d + %d = %d)", y, coord.y, coord.y+y))
		}

		for _, otherSpaceship := range b.spaceships {
			for _, otherCoord := range otherSpaceship.coords {
				if coord.x+x == otherCoord.x && coord.y+y == otherCoord.y {
					return errors.New(fmt.Sprintf("Failed to add spaceship, coord already contains spaceship (%dx%d)", otherCoord.x, otherCoord.y))
				}
			}
		}
	}

	// offset the spaceship coords with the coords it's placed on
	for _, coord := range spaceship.coords {
		coord.x += x
		coord.y += y
	}

	// add spaceship to board
	b.spaceships = append(b.spaceships, spaceship)

	return nil
}

type Spaceship struct {
	coords []*Coord
	hits   []*Coord
	dead   bool
}

// @TODO: should sanitize any padding
func SpaceshipFromPattern(pattern []string) (*Spaceship, error) {
	// sanity check the input
	if len(pattern) > ROWS {
		return nil, errors.New("pattern too many rows")
	}

	// sanity check the input
	for _, row := range pattern {
		if len(row) > COLS {
			return nil, errors.New("pattern too many cols")
		}

		// @TODO: is there a nicer way to do this with a builtin?
		for _, char := range []byte(row) {
			if char != byte(CoordBlank) && char != byte(CoordShip) {
				return nil, errors.New("pattern incorrect symbol for coord")
			}
		}
	}

	spaceship := &Spaceship{
		hits: make([]*Coord, 0),
		dead: false,
	}

	// parse the input
	for y, row := range pattern {
		for x, char := range []byte(row) {
			coordState := CoordState(char)

			switch coordState {
			case CoordBlank:
				// - nothing to do
			case CoordShip:
				spaceship.coords = append(spaceship.coords, &Coord{x: uint8(x), y: uint8(y)})
			}
		}
	}

	if len(spaceship.coords) == 0 {
		return nil, errors.New("blank spaceship")
	}

	return spaceship, nil
}
