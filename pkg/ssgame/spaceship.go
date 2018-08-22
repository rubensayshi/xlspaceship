package ssgame

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

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

func (s *Spaceship) Copy() *Spaceship {
	return &Spaceship{
		pattern: s.pattern,
		coords:  s.coords.Copy(),
		hits:    s.hits.Copy(),
		dead:    s.dead,
	}
}

func (s *Spaceship) CopyWithOffset(x uint8, y uint8) *Spaceship {
	newS := s.Copy()
	newS.Offset(x, y)

	return newS
}

func (s *Spaceship) Offset(x uint8, y uint8) {
	for _, coords := range s.coords {
		coords.x += x
		coords.y += y
	}
}

func (s *Spaceship) String() string {
	return fmt.Sprintf("%s", strings.Join(s.pattern, "\n"))
}
