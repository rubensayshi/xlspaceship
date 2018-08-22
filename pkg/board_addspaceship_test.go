package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoard_AddSpaceshipOnCoordsSimple0x0(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship, 0, 0)
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsSimple13x15(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship, 13, 15)
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsSimpleVert15x13(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship, 15, 13)
	assert.NoError(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidSimple14x15(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship, 14, 15)
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidSimpleVert15x14(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship, err := SpaceshipFromPattern([]string{
		"*",
		"*",
		"*",
	})
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship, 15, 14)
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlap(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
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

	err = board.AddSpaceshipOnCoords(spaceship1, 0, 0)
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2, 0, 0)
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlap3X0(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
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

	err = board.AddSpaceshipOnCoords(spaceship1, 3, 0)
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2, 3, 0)
	assert.Error(err)
}

func TestBoard_AddSpaceshipOnCoordsInvalidOverlapWingerB(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
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

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)

	spaceship1, err := SpaceshipFromPattern(SpaceshipPatternWinger)
	assert.NoError(err)

	spaceship2, err := SpaceshipFromPattern(SpaceshipPatternBClass)
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship1, 8, 9)
	assert.NoError(err)

	err = board.AddSpaceshipOnCoords(spaceship2, 9, 10)
	assert.Error(err)
}
