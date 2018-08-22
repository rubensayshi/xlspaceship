package ssgame

import (
	"fmt"
	"regexp"
	"strconv"

	"math/rand"

	"github.com/pkg/errors"
)

func init() {
	// this should be replaced by crypto/rand with proper seeding for secure random numbers
	//  but for this exercise it's much nicer if it's not really random
	rand.Seed(1)
}

const ROWS = 16
const COLS = 16

var coordsRegex = regexp.MustCompile(`^([0-9a-fA-F])[xX]([0-9a-fA-F])$`)

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

type ShotStatus uint8

func ShotStatusFromString(statusStr string) (ShotStatus, error) {
	switch statusStr {
	case ShotStatusMissStr:
		return ShotStatusMiss, nil
	case ShotStatusHitStr:
		return ShotStatusHit, nil
	case ShotStatusKillStr:
		return ShotStatusKill, nil
	default:
		return ShotStatusMiss, errors.New("Invalid ShotStatus string")
	}
}

func (c ShotStatus) String() string {
	switch c {
	case ShotStatusMiss:
		return ShotStatusMissStr
	case ShotStatusHit:
		return ShotStatusHitStr
	case ShotStatusKill:
		return ShotStatusKillStr
	}

	panic("Unreachable")
}

const (
	ShotStatusMiss ShotStatus = 0
	ShotStatusHit  ShotStatus = 1
	ShotStatusKill ShotStatus = 2

	ShotStatusMissStr string = "miss"
	ShotStatusHitStr  string = "hit"
	ShotStatusKillStr string = "kill"
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

type CoordsGroup []*Coords

func CoordsGroupFromSalvoStrings(salvo []string) (CoordsGroup, error) {
	cg := make(CoordsGroup, len(salvo))
	for i, coordsStr := range salvo {
		coords, err := CoordsFromString(coordsStr)
		if err != nil {
			return nil, err
		}

		cg[i] = coords
	}

	return cg, nil
}

func (cg CoordsGroup) Contains(coords *Coords) bool {
	for _, cgCoords := range cg {
		if cgCoords.x == coords.x && cgCoords.y == coords.y {
			return true
		}
	}

	return false
}

func (cg CoordsGroup) Copy() CoordsGroup {
	newCg := make(CoordsGroup, len(cg))
	for i, coords := range cg {
		newCg[i] = &Coords{x: coords.x, y: coords.y}
	}

	return newCg
}

type ShotResult struct {
	Coords     *Coords
	ShotStatus ShotStatus
}
