package ssgame

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Spaceship struct {
	coords CoordsGroup
	hits   CoordsGroup
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
		hits: make(CoordsGroup, 0),
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
				spaceship.coords = append(spaceship.coords, &Coords{x: int8(x), y: int8(y)})
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
		coords: s.coords.Copy(),
		hits:   s.hits.Copy(),
		dead:   s.dead,
	}
}

func (s *Spaceship) CopyWithOffset(x int8, y int8) *Spaceship {
	newS := s.Copy()
	newS.Offset(x, y)

	return newS
}

func (s *Spaceship) Offset(x int8, y int8) {
	for _, coords := range s.coords {
		coords.x += x
		coords.y += y
	}
}

func (s *Spaceship) Rotate(rotate uint16) error {
	// rotate coords
	switch rotate {
	case 90:
		for _, coords := range s.coords {
			coords.x, coords.y = coords.y*-1, coords.x
		}
	case 180:
		for _, coords := range s.coords {
			coords.x, coords.y = coords.x, coords.y*-1
		}
	case 270:
		for _, coords := range s.coords {
			coords.x, coords.y = coords.y, coords.x*-1
		}

	default:
		return errors.Errorf("Unsupported rotate: %d", rotate)
	}

	// offset coords if they went out of bounds
	var offsetX int8
	var offsetY int8
	for _, coords := range s.coords {
		if 0-coords.x > offsetX {
			offsetX = 0 - coords.x
		}
		if 0-coords.y > offsetY {
			offsetY = 0 - coords.y
		}
	}

	for _, coords := range s.coords {
		coords.x += offsetX
		coords.y += offsetY
	}

	return nil
}

func (s *Spaceship) ToPattern() []string {
	var maxX int8 = 0
	var maxY int8 = 0
	for _, coords := range s.coords {
		if coords.x > maxX {
			maxX = coords.x
		}
		if coords.y > maxY {
			maxY = coords.y
		}
	}

	pattern := make([][]byte, maxY+1)
	for y, _ := range pattern {
		pattern[y] = make([]byte, maxX+1)

		var x int8
		for ; x <= maxX; x++ {
			pattern[y][x] = byte(CoordsBlank)
		}
	}

	// add spaceships to the pattern
	for _, coords := range s.coords {
		pattern[coords.y][coords.x] = byte(CoordsShip)
	}

	// turn the byte arrays into strings
	res := make([]string, len(pattern))
	for y, row := range pattern {
		res[y] = strings.TrimRight(string(row), ".")
	}

	return res
}

func (s *Spaceship) String() string {
	return fmt.Sprintf("%s", strings.Join(s.ToPattern(), "\n"))
}
