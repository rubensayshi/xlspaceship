package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoardFromPatternEmpty(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern([]string{
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
	})
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(0, len(board.hits))
	assert.Equal(0, len(board.misses))
	assert.Equal(0, len(board.spaceships))
}

func TestBoardFromPatternWithMarks(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern([]string{
		"X-*.............",
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
	})
	assert.NoError(err)
	assert.NotNil(board)
}

func TestBoardFromPatternInvalidRows(t *testing.T) {
	assert := require.New(t)

	_, err := BoardFromPattern([]string{
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
	})
	assert.Error(err)

	_, err = BoardFromPattern([]string{
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
		"................",
	})
	assert.Error(err)
}

func TestBoardFromPatternInvalidCols(t *testing.T) {
	assert := require.New(t)

	invalidRows := []string{
		".................",
		"...............",
		"..............",
		".............",
		"............",
		"...........",
		"..........",
		".........",
		"........",
		".......",
		"......",
		".....",
		"....",
		"...",
		"..",
		".",
		"",
	}

	for row := 0; row < 16; row++ {
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

		for _, invalidRow := range invalidRows {
			pattern[row] = invalidRow

			_, err := BoardFromPattern(pattern)
			assert.Error(err)
		}
	}
}

func TestBoardFromPatternInvalidChars(t *testing.T) {
	assert := require.New(t)

	invalidRows := []string{
		"A...............",
		"1...............",
		"?...............",
		"€...............", // unicode char, will actually give invalid row length err, but it's fine to test it here for now
	}

	for row := 0; row < 16; row++ {
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

		for _, invalidRow := range invalidRows {
			pattern[row] = invalidRow

			_, err := BoardFromPattern(pattern)
			assert.Error(err)
		}
	}
}

func TestBoardFromPatternHits(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern([]string{
		"X...............",
		".X..............",
		"..X.............",
		"...X............",
		"....X...........",
		".....X..........",
		"......X.........",
		".......X........",
		"........X.......",
		".........X......",
		"..........X.....",
		"...........X....",
		"............X...",
		".............X..",
		"..............X.",
		"...............X",
	})
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(16, len(board.hits))
	assert.Equal(0, len(board.misses))
	assert.Equal(0, len(board.spaceships))
}

func TestBoardFromPatternMisses(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern([]string{
		"-...............",
		".-..............",
		"..-.............",
		"...-............",
		"....-...........",
		".....-..........",
		"......-.........",
		".......-........",
		"........-.......",
		".........-......",
		"..........-.....",
		"...........-....",
		"............-...",
		".............-..",
		"..............-.",
		"...............-",
	})
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(0, len(board.hits))
	assert.Equal(16, len(board.misses))
	assert.Equal(0, len(board.spaceships))
}

func TestBoardFromPatternMixed(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern([]string{
		"X...............",
		".-..............",
		"..X.............",
		"...-............",
		"....X...........",
		".....-..........",
		"......X.........",
		".......-........",
		"........X.......",
		".........-......",
		"..........X.....",
		"...........-....",
		"............X...",
		".............-..",
		"..............X.",
		"...............-",
	})
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(8, len(board.hits))
	assert.Equal(8, len(board.misses))
	assert.Equal(0, len(board.spaceships))
}

func TestBoardToPattern(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
		"X...............",
		".-..............",
		"..X.............",
		"...-............",
		"....X...........",
		".....-..........",
		"......X.........",
		".......-........",
		"........X.......",
		".........-......",
		"..........X.....",
		"...........-....",
		"............X...",
		".............-..",
		"..............X.",
		"...............-",
	}

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(8, len(board.hits))
	assert.Equal(8, len(board.misses))
	assert.Equal(0, len(board.spaceships))

	assert.Equal(pattern, board.ToPattern())
}

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
