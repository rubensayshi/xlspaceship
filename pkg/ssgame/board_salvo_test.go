package ssgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func NewBasicTestBoardWithSpaceship(assert *require.Assertions) *Board {
	board, err := BoardFromPattern(BlankBoardPattern())
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

	res = board.ApplyShot(&Coords{0, 0})
	assert.NotNil(res)
	assert.Equal(ShotStatusMiss, res.ShotStatus)
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
