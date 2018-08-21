package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoardFromPatternEmpty(t *testing.T) {
	assert := require.New(t)

	board, err := BoardFromPattern(BlankBoardPattern())
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
		pattern := BlankBoardPattern()

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
		pattern := BlankBoardPattern()

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
