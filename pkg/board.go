package pkg

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

const ROWS = 16
const COLS = 16

var coordsRegex = regexp.MustCompile(`^([0-9a-fA-F])x([0-9a-fA-F])$`)

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

type CoordsState byte

func (c CoordsState) String() string {
	switch c {
	case CoordsBlank:
		return CoordsBlankStr
	case CoordsShip:
		return CoordsShipStr
	case CoordsHit:
		return CoordsHitStr
	case CoordsMiss:
		return CoordsMissStr
	}

	panic("Unreachable")
}

const (
	CoordsBlank CoordsState = '.'
	CoordsShip  CoordsState = '*'
	CoordsHit   CoordsState = 'X'
	CoordsMiss  CoordsState = '-'

	CoordsBlankStr string = "."
	CoordsShipStr  string = "*"
	CoordsHitStr   string = "X"
	CoordsMissStr  string = "-"
)

type Coords struct {
	x uint8
	y uint8
}

func CoordsFromString(coordsStr string) (*Coords, error) {
	matches := coordsRegex.FindStringSubmatch(coordsStr)
	if len(matches) != 3 {
		return nil, errors.New(fmt.Sprintf("Failed to parse Coords [%s]", coordsStr))
	}

	x, err := strconv.ParseInt(matches[1], 16, 8)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse Coords [%s]", coordsStr))
	}
	y, err := strconv.ParseInt(matches[2], 16, 8)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse Coords [%s]", coordsStr))
	}

	coords := &Coords{x: uint8(x), y: uint8(y)}

	return coords, nil
}

func (c Coords) String() string {
	return fmt.Sprintf("%Xx%X", c.x, c.y)
}

type SalvoResult struct {
}

type Board struct {
	spaceships []*Spaceship
	hits       []*Coords
	misses     []*Coords
}

func NewRandomBoard() (*Board, error) {
	// we retry to create a random board 100 times incase the spaceships didn't fit (should never happen with default board size and spaceships)
	for i := 0; i < 100; i++ {
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

			// attemp to add spaceship, if we fail we nil the board so that we keep trying
			err = board.AddSpaceship(spaceship)
			if err != nil {
				board = nil
				break
			}
		}

		if board != nil {
			return board, nil
		}
	}

	return nil, errors.New("Failed to create a random board")
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
		hits:       make([]*Coords, 0),
		misses:     make([]*Coords, 0),
	}

	// parse the input
	for y, row := range pattern {
		for x, char := range []byte(row) {
			coordsState := CoordsState(char)

			switch coordsState {
			case CoordsBlank:
				// - nothing to do
			case CoordsShip:
				// @TODO: not implemented
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
	// @TODO: we should store coords of existing spaceships so this isn't O(N2)
	for _, coords := range spaceship.coords {
		if coords.x+x >= ROWS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, x overflow (%d + %d = %d)", x, coords.x, coords.x+x))
		}
		if coords.y+y >= COLS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, y overflow (%d + %d = %d)", y, coords.y, coords.y+y))
		}

		for _, otherSpaceship := range b.spaceships {
			for _, otherCoord := range otherSpaceship.coords {
				if coords.x+x == otherCoord.x && coords.y+y == otherCoord.y {
					return errors.New(fmt.Sprintf("Failed to add spaceship, coords already contains spaceship (%dx%d)", otherCoord.x, otherCoord.y))
				}
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

func (b *Board) ReceiveSalvo(salvo []*Coords) *SalvoResult {

	return nil
}

type Spaceship struct {
	coords []*Coords
	hits   []*Coords
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
			if char != byte(CoordsBlank) && char != byte(CoordsShip) {
				return nil, errors.New("pattern incorrect symbol for coords")
			}
		}
	}

	spaceship := &Spaceship{
		hits: make([]*Coords, 0),
		dead: false,
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
