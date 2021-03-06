package ssgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func NewBasicTestBoardWithSpaceship(assert *require.Assertions) *SelfBoard {
	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(0, 0))
	assert.NoError(err)

	return board
}

func TestBoard_ApplyShotMiss(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{0, 1})
	assert.NotNil(res)
	assert.Equal(ShotStatusMiss, res.ShotStatus)
}

func TestBoard_ApplyShotHit(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)
}

func TestBoard_ApplyShotHitTwice(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)
	assert.Equal(1, board.CountHits())
	assert.Equal(0, board.CountMisses())

	res = board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusMiss, res.ShotStatus)
	assert.Equal(1, board.CountHits())
	assert.Equal(0, board.CountMisses())
}

func TestBoard_ApplyShotKill(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)

	res = board.ApplyShot(&Coords{1, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)

	res = board.ApplyShot(&Coords{2, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusKill, res.ShotStatus)
}

func TestBoard_ReceiveSalvoKill(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ReceiveSalvo(CoordsGroup{
		{0, 0},
		{1, 0},
		{2, 0},
	})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res[0].ShotStatus)
	assert.Equal(ShotStatusHit, res[1].ShotStatus)
	assert.Equal(ShotStatusKill, res[2].ShotStatus)
}

func TestBoard_ApplyShotKillTwice(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)

	res = board.ApplyShot(&Coords{1, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusHit, res.ShotStatus)

	res = board.ApplyShot(&Coords{2, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusKill, res.ShotStatus)

	res = board.ApplyShot(&Coords{2, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusMiss, res.ShotStatus)
}

func TestBoard_ApplyShotStatus(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankOpponentBoard(1)
	assert.NoError(err)

	board.ApplyShotStatus(&Coords{2, 0}, ShotStatusHit)
	board.ApplyShotStatus(&Coords{3, 0}, ShotStatusMiss)

	assert.Equal([]string{
		"..X-............",
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
	}, board.ToPattern())

	assert.Equal(1, board.CountShipsAlive())
	assert.Equal(false, board.AllShipsDead())

	board.ApplyShotStatus(&Coords{1, 0}, ShotStatusKill)

	assert.Equal([]string{
		".XX-............",
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
	}, board.ToPattern())

	assert.Equal(0, board.CountShipsAlive())
	assert.Equal(true, board.AllShipsDead())
}

func TestBoard_ApplyShotStatusKill(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankOpponentBoard(1)
	assert.NoError(err)

	board.ApplyShotStatus(&Coords{2, 0}, ShotStatusKill)

	assert.Equal([]string{
		"..X.............",
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
	}, board.ToPattern())

	assert.Equal(uint8(0), board.spaceshipsAlive)
	assert.Equal(0, board.CountShipsAlive())
	assert.Equal(true, board.AllShipsDead())
}

func TestBoard_ApplyShotOutOfBounds(t *testing.T) {
	assert := require.New(t)

	board := NewBasicTestBoardWithSpaceship(assert)

	res := board.ApplyShot(&Coords{20, 20})
	assert.NotNil(res)
	assert.Equal(ShotStatusMiss, res.ShotStatus)
}

func TestBoard_ApplyShotStatusOutOfBounds(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankOpponentBoard(1)
	assert.NoError(err)

	board.ApplyShotStatus(&Coords{20, 20}, ShotStatusKill)

	assert.Equal([]string{
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
	}, board.ToPattern())
}
