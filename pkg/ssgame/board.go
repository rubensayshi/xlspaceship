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

// the base type for our boards to share
type BaseBoard struct {
	grid [][]*GridCell
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
	return &BaseBoard{}
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

	// init the grid with rows
	board.grid = make([][]*GridCell, ROWS)

	// parse the input and add them to the grid
	for y, row := range pattern {
		board.grid[y] = make([]*GridCell, COLS)

		for x, char := range []byte(row) {
			coordsState := CoordsState(char)

			board.grid[y][x] = &GridCell{
				coords: &Coords{x: int8(x), y: int8(y)},
				state:  coordsState,
			}
		}
	}

	return nil
}

func (b *BaseBoard) buildPattern() [][]byte {
	pattern := make([][]byte, ROWS)
	for y, row := range b.grid {
		pattern[y] = make([]byte, COLS)

		for x, cell := range row {
			pattern[y][x] = byte(cell.state)
		}
	}

	return pattern
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
	pattern := b.buildPattern()

	return b.patternToStrings(pattern)
}

// Count the hits, could be stored internally instead of recounting every time, but we don't actually use this outside of tests currently
func (b *BaseBoard) CountHits() int {
	hits := 0

	for _, row := range b.grid {
		for _, cell := range row {
			if cell.state == CoordsHit {
				hits++
			}
		}
	}

	return hits
}

// Count the misses, could be stored internally instead of recounting every time, but we don't actually use this outside of tests currently
func (b *BaseBoard) CountMisses() int {
	misses := 0

	for _, row := range b.grid {
		for _, cell := range row {
			if cell.state == CoordsMiss {
				misses++
			}
		}
	}

	return misses
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

	// add spaceship to grid
	for _, coords := range spaceship.coords {
		if b.grid[coords.y][coords.x].state != CoordsHit {
			b.grid[coords.y][coords.x].state = CoordsShip
		}
		b.grid[coords.y][coords.x].spaceship = spaceship
	}

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
	status := ShotStatusMiss

	// check if shot is within bounds of our grid
	if int(shot.y) < len(b.grid) && int(shot.x) < len(b.grid[shot.y]) {
		cell := b.grid[shot.y][shot.x]

		// check if shot was on a ship (note; previous hits will fail because they're already CoordsHit), this is intended
		if cell.state == CoordsShip {
			cell.state = CoordsHit

			if cell.spaceship.coords.Contains(shot) && !cell.spaceship.hits.Contains(shot) {
				status = ShotStatusHit

				// add the coords as a hit
				cell.spaceship.hits = append(cell.spaceship.hits, shot)

				fmt.Printf("")

				// if we've hit all the coords then it's a kill
				if len(cell.spaceship.hits) == len(cell.spaceship.coords) {
					cell.spaceship.dead = true
					status = ShotStatusKill
				}
			}
		} else if cell.state == CoordsHit {
			// nothing to do, already a hit so leave untouched and return MISS
		} else {
			cell.state = CoordsMiss
		}
	}

	res := &ShotResult{
		shot,
		status,
	}

	return res
}

// apply one of our shots to opponent's board using the status our opponent told us of the shot
func (b *OpponentBoard) ApplyShotStatus(shot *Coords, status ShotStatus) {

	// check if shot is within bounds of our grid
	if int(shot.y) < len(b.grid) && int(shot.x) < len(b.grid[shot.y]) {
		switch status {
		case ShotStatusMiss:
			b.grid[shot.y][shot.x].state = CoordsMiss
		case ShotStatusHit:
			b.grid[shot.y][shot.x].state = CoordsHit
		case ShotStatusKill:
			b.grid[shot.y][shot.x].state = CoordsHit
			b.spaceshipsAlive--
		}
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

func (b *BaseBoard) String() string {
	return fmt.Sprintf("%s", strings.Join(b.ToPattern(), "\n"))
}
