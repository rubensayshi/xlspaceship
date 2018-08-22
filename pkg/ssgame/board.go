package ssgame

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

func BlankBoardPattern() []string {
	return []string{
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
	}
}

type Board struct {
	spaceships []*Spaceship
	hits       CoordsGroup
	misses     CoordsGroup
}

func NewRandomBoard(spaceships [][]string) (*Board, error) {
	// we retry to create a random board 100 times incase the spaceships didn't fit (should never happen with default board size and spaceships)
	for i := 0; i < 100; i++ {
		board, err := newRandomBoard(spaceships)
		if err != nil {
			return nil, err
		}

		if board != nil {
			return board, nil
		}
	}

	return nil, errors.New("Failed to create a random board")
}

func newRandomBoard(spaceships [][]string) (*Board, error) {
	board := &Board{}

	for _, spaceshipPattern := range spaceships {
		spaceship, err := SpaceshipFromPattern(spaceshipPattern)
		if err != nil {
			return nil, err
		}

		// attempt to add spaceship, if we fail we nil the board so that we keep trying
		err = board.AddSpaceship(spaceship)
		if err != nil {
			board = nil
			break
		}
	}

	return board, nil
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
			if char != byte(CoordsBlank) && char != byte(CoordsShip) && char != byte(CoordsHit) && char != byte(CoordsMiss) {
				return nil, errors.New("pattern incorrect symbol for coords")
			}
		}
	}

	board := &Board{
		spaceships: make([]*Spaceship, 0),
		hits:       make(CoordsGroup, 0),
		misses:     make(CoordsGroup, 0),
	}

	// parse the input
	for y, row := range pattern {
		for x, char := range []byte(row) {
			coordsState := CoordsState(char)

			switch coordsState {
			case CoordsBlank:
				// - nothing to do
			case CoordsShip:
				// @TODO: not implemented, the dream is to store them and try and match the patterns to our known spaceships
			case CoordsHit:
				board.hits = append(board.hits, &Coords{x: uint8(x), y: uint8(y)})
			case CoordsMiss:
				board.misses = append(board.misses, &Coords{x: uint8(x), y: uint8(y)})
			}
		}
	}

	return board, nil
}

func (b *Board) String() string {
	return fmt.Sprintf("%s", strings.Join(b.ToPattern(), "\n"))
}

func (b *Board) AllShipsDead() bool {
	for _, spaceship := range b.spaceships {
		if !spaceship.dead {
			return false
		}
	}

	return true
}

func (b *Board) ToPattern() []string {
	// @TODO: considering the board size is constant we could just have a const string to copy for this instead of building the blank state everytime
	pattern := make([][]byte, ROWS)
	for y, _ := range pattern {
		pattern[y] = make([]byte, COLS)

		for x := 0; x < COLS; x++ {
			pattern[y][x] = byte(CoordsBlank)
		}
	}

	// add spaceships to the pattern
	for _, spaceship := range b.spaceships {
		for _, coords := range spaceship.coords {
			pattern[coords.y][coords.x] = byte(CoordsShip)
		}
	}

	// add hits to the pattern (will overwrite spaceship coords)
	for _, hit := range b.hits {
		pattern[hit.y][hit.x] = byte(CoordsHit)
	}

	// add misses to the pattern
	for _, miss := range b.misses {
		pattern[miss.y][miss.x] = byte(CoordsMiss)
	}

	// turn the byte arrays into strings
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
	// @TODO: we should store coords of existing spaceships so we don't have to loop over them
	for _, coords := range spaceship.coords {
		offsetCoords := &Coords{
			x: coords.x + x,
			y: coords.y + y,
		}

		// check spaceship stays within bounds
		if offsetCoords.x >= COLS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, y overflow (%d + %d = %d)", y, coords.y, coords.y+y))
		}
		if offsetCoords.y >= ROWS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, x overflow (%d + %d = %d)", x, coords.x, coords.x+x))
		}

		// check spaceship doesn't overlap with other spaceships
		for _, otherSpaceship := range b.spaceships {
			if otherSpaceship.coords.Contains(offsetCoords) {
				return errors.New(fmt.Sprintf("Failed to add spaceship, coords already contains spaceship (%dx%d)", coords.x, coords.y))
			}
		}
	}

	// offset the spaceship coords with the coords it's placed on
	for _, coords := range spaceship.coords {
		coords.x += x
		coords.y += y
	}

	// add spaceship to board
	b.spaceships = append(b.spaceships, spaceship)

	return nil
}

func (b *Board) ReceiveSalvo(salvo CoordsGroup) []*ShotResult {
	res := make([]*ShotResult, len(salvo))
	for i, shot := range salvo {
		res[i] = b.ApplyShot(shot)
	}

	return res
}

func (b *Board) ApplyShot(shot *Coords) *ShotResult {
	// @TODO: same as with AddSpacehipOnCoords it would be a good optimization to store the coords of the ships so we don't have to loop over them
	status := ShotStatusMiss

	for _, spaceship := range b.spaceships {
		// check if it's FRESH hit
		if spaceship.coords.Contains(shot) && !spaceship.hits.Contains(shot) {
			status = ShotStatusHit

			// add the coords as a hit
			spaceship.hits = append(spaceship.hits, shot)
			b.hits = append(b.hits, shot)

			fmt.Printf("")

			// if we've hit all the coords then it's a kill
			if len(spaceship.hits) == len(spaceship.coords) {
				spaceship.dead = true
				status = ShotStatusKill
			}

			// break, can't have more than 1 hit
			break
		}
	}

	if status == ShotStatusMiss {
		b.misses = append(b.misses, shot)
	}

	res := &ShotResult{
		shot,
		status,
	}

	return res
}

type Spaceship struct {
	pattern []string
	coords  CoordsGroup
	hits    CoordsGroup
	dead    bool
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
			if char != byte(CoordsBlank) && char != byte(CoordsShip) {
				return nil, errors.New("pattern incorrect symbol for coords")
			}
		}
	}

	spaceship := &Spaceship{
		pattern: pattern,
		hits:    make(CoordsGroup, 0),
		dead:    false,
	}

	// parse the input
	for y, row := range pattern {
		for x, char := range []byte(row) {
			coordsState := CoordsState(char)

			switch coordsState {
			case CoordsBlank:
				// - nothing to do
			case CoordsShip:
				spaceship.coords = append(spaceship.coords, &Coords{x: uint8(x), y: uint8(y)})
			}
		}
	}

	if len(spaceship.coords) == 0 {
		return nil, errors.New("blank spaceship")
	}

	return spaceship, nil
}

func (s *Spaceship) String() string {
	return fmt.Sprintf("%s", strings.Join(s.pattern, "\n"))
}
