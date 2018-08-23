package ssgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoard_AddSpaceshipOnCoordsSimple0x0(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(0, 0))
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsSimple13x15(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(13, 15))
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsSimpleVert15x13(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(15, 13))
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidSimple14x15(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(14, 15))
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidSimpleVert15x14(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship.CopyWithOffset(15, 14))
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlap(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship1, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	spaceship2, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship1.CopyWithOffset(0, 0))
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2.CopyWithOffset(0, 0))
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlap3X0(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship1, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	spaceship2, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship1.CopyWithOffset(3, 0))
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2.CopyWithOffset(3, 0))
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlapWingerB(t *testing.T) {
	assert := require.New(t)

	board, err := NewBlankSelfBoard()
	assert.NoError(err)

	spaceship1, err := SpaceshipFromPattern(SpaceshipPatternWinger)
	assert.NoError(err)

	spaceship2, err := SpaceshipFromPattern(SpaceshipPatternBClass)
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship1.CopyWithOffset(8, 9))
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2.CopyWithOffset(9, 10))
	assert.Error(err)
}
