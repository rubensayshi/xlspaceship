package ssgame

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pkg/errors"
)

// helper function for a blank board
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

// the base type for our baords to share
type BaseBoard struct {
	hits   CoordsGroup
	misses CoordsGroup
}

// our own board which contains our own placed shaceships
type SelfBoard struct {
	*BaseBoard
	spaceships []*Spaceship
}

// our opponent's board for which we don't know his spaceships, we do know how many there are left alive
type OpponentBoard struct {
	*BaseBoard
	spaceshipsAlive uint8
}

// generate a random board for ourselves with the specified spaceships
//  we retry to create a random board 100 times incase the spaceships didn't fit
func NewRandomSelfBoard(spaceships [][]string) (*SelfBoard, error) {
	for i := 0; i < 100; i++ {
		board, err := newRandomSelfBoard(spaceships)
		if err != nil {
			return nil, err
		}

		if board != nil {
			return board, nil
		}
	}

	return nil, errors.New("Failed to create a random board")
}

// generate a random board for ourselves with the specified spaceships
//  internal function for NewRandomSelfBoard to use
// board can be nil when we failed to place a spaceship
func newRandomSelfBoard(spaceships [][]string) (*SelfBoard, error) {
	board, err := NewBlankSelfBoard()
	if err != nil {
		return nil, err
	}

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

func newBaseBoard() *BaseBoard {
	return &BaseBoard{
		hits:   make(CoordsGroup, 0),
		misses: make(CoordsGroup, 0),
	}
}

func NewSelfBoard() *SelfBoard {
	return &SelfBoard{
		BaseBoard:  newBaseBoard(),
		spaceships: make([]*Spaceship, 0),
	}
}

func NewBlankSelfBoard() (*SelfBoard, error) {
	board := NewSelfBoard()

	err := FillBoardFromPattern(board.BaseBoard, BlankBoardPattern())
	if err != nil {
		return nil, err
	}

	return board, nil
}

func NewOpponentBoard(spaceshipsAlive uint8) *OpponentBoard {
	return &OpponentBoard{
		BaseBoard:       newBaseBoard(),
		spaceshipsAlive: spaceshipsAlive,
	}
}

func NewBlankOpponentBoard(spaceshipsAlive uint8) (*OpponentBoard, error) {
	board := NewOpponentBoard(spaceshipsAlive)

	err := FillBoardFromPattern(board.BaseBoard, BlankBoardPattern())
	if err != nil {
		return nil, err
	}

	return board, nil
}

// fill a board with a pattern
func FillBoardFromPattern(board *BaseBoard, pattern []string) error {
	// sanity check the input
	if len(pattern) != ROWS {
		return errors.New("pattern incorrect amount of rows")
	}

	// sanity check the input
	for _, row := range pattern {
		if len(row) != COLS {
			return errors.New("pattern incorrect amount of cols")
		}

		// @TODO: is there a nicer way to do this with a builtin?
		for _, char := range []byte(row) {
			if char != byte(CoordsBlank) && char != byte(CoordsShip) && char != byte(CoordsHit) && char != byte(CoordsMiss) {
				return errors.New("pattern incorrect symbol for coords")
			}
		}
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
				board.hits = append(board.hits, &Coords{x: int8(x), y: int8(y)})
			case CoordsMiss:
				board.misses = append(board.misses, &Coords{x: int8(x), y: int8(y)})
			}
		}
	}

	return nil
}

func (b *BaseBoard) buildBasePattern() [][]byte {
	// @TODO: considering the board size is constant we could just have a const string to copy for this instead of building the blank state everytime
	pattern := make([][]byte, ROWS)
	for y, _ := range pattern {
		pattern[y] = make([]byte, COLS)

		for x := 0; x < COLS; x++ {
			pattern[y][x] = byte(CoordsBlank)
		}
	}

	return pattern
}

func (b *BaseBoard) applyHitsAndMissesToPattern(pattern [][]byte) {
	// add hits to the pattern (will overwrite spaceship coords)
	for _, hit := range b.hits {
		pattern[hit.y][hit.x] = byte(CoordsHit)
	}

	// add misses to the pattern
	for _, miss := range b.misses {
		pattern[miss.y][miss.x] = byte(CoordsMiss)
	}
}

func (b *BaseBoard) patternToStrings(pattern [][]byte) []string {
	// turn the byte arrays into strings
	res := make([]string, ROWS)
	for y, row := range pattern {
		res[y] = string(row)
	}

	return res
}

func (b *BaseBoard) ToPattern() []string {
	pattern := b.buildBasePattern()
	b.applyHitsAndMissesToPattern(pattern)

	return b.patternToStrings(pattern)
}

func (b *SelfBoard) applyShipsToPattern(pattern [][]byte) {
	// add spaceships to the pattern
	for _, spaceship := range b.spaceships {
		for _, coords := range spaceship.coords {
			pattern[coords.y][coords.x] = byte(CoordsShip)
		}
	}
}

func (b *SelfBoard) ToPattern() []string {
	pattern := b.buildBasePattern()
	b.applyShipsToPattern(pattern)
	b.applyHitsAndMissesToPattern(pattern)

	return b.patternToStrings(pattern)
}

// attempt to add a spaceship on random locations until we succeed
//  if we reach the max N attempts then just error out
func (b *SelfBoard) AddSpaceship(spaceship *Spaceship) error {
	N := 10000
	// @TODO: this could be heavily optimized as we know we don't have to try adding a spaceship of 3 high on Y > 15 - 3
	for i := 0; i < N; i++ {
		// randomize x, y offset and rotation
		x := rand.Intn(COLS)
		y := rand.Intn(ROWS)
		rotate := rand.Intn(3) * 90

		newSpaceship := spaceship.CopyWithOffset(int8(x), int8(y)).CopyWithRotate(uint16(rotate))

		err := b.AddSpaceshipOnCoords(newSpaceship)
		// we're done if no err
		if err == nil {
			return nil
		}
	}

	return errors.New("Failed to add spaceship, seems impossible")
}

// add a spaceship on specified locations
//  will error when it's out of bound or overlapping with an existing spaceship
func (b *SelfBoard) AddSpaceshipOnCoords(spaceship *Spaceship) error {
	// @TODO: we should store coords of existing spaceships so we don't have to loop over them
	for _, coords := range spaceship.coords {
		// check spaceship stays within bounds
		if coords.x >= COLS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, y overflow (%s)", coords))
		}
		if coords.y >= ROWS {
			return errors.New(fmt.Sprintf("Failed to add spaceship, x overflow (%s)", coords))
		}

		// check spaceship doesn't overlap with other spaceships
		for _, otherSpaceship := range b.spaceships {
			if otherSpaceship.coords.Contains(coords) {
				return errors.New(fmt.Sprintf("Failed to add spaceship, coords already contains spaceship (%s)", coords))
			}
		}
	}

	// add spaceship to board
	b.spaceships = append(b.spaceships, spaceship)

	return nil
}

// receive a salvo onto our board
func (b *SelfBoard) ReceiveSalvo(salvo CoordsGroup) []*ShotResult {
	res := make([]*ShotResult, len(salvo))
	for i, shot := range salvo {
		res[i] = b.ApplyShot(shot)
	}

	return res
}

// apply a shot to our board
func (b *SelfBoard) ApplyShot(shot *Coords) *ShotResult {
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

	// store the miss, but only if it wasn't already stored as hit or miss
	if status == ShotStatusMiss && !b.hits.Contains(shot) && !b.misses.Contains(shot) {
		b.misses = append(b.misses, shot)
	}

	res := &ShotResult{
		shot,
		status,
	}

	return res
}

// apply one of our shots to opponent's board using the status our opponent told us of the shot
func (b *OpponentBoard) ApplyShotStatus(shot *Coords, status ShotStatus) {
	switch status {
	case ShotStatusMiss:
		b.misses = append(b.misses, shot)
	case ShotStatusHit:
		b.hits = append(b.hits, shot)
	case ShotStatusKill:
		b.hits = append(b.hits, shot)
		b.spaceshipsAlive--
	}
}

func (b *SelfBoard) Spaceships() []*Spaceship {
	return b.spaceships
}

func (b *SelfBoard) CountShipsAlive() int {
	i := 0
	for _, spaceship := range b.spaceships {
		if !spaceship.dead {
			i++
		}
	}

	return i
}

func (b *SelfBoard) AllShipsDead() bool {
	return b.CountShipsAlive() == 0
}

func (b *OpponentBoard) CountShipsAlive() int {
	return int(b.spaceshipsAlive)
}

func (b *OpponentBoard) AllShipsDead() bool {
	return b.spaceshipsAlive == 0
}

func (b *OpponentBoard) String() string {
	return fmt.Sprintf("%s", strings.Join(b.ToPattern(), "\n"))
}

func (b *SelfBoard) String() string {
	return fmt.Sprintf("%s", strings.Join(b.ToPattern(), "\n"))
}
