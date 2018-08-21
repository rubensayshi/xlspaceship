package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	assert := require.New(t)

	game := NewGame("player-1")

	assert.Equal("player-1", game.OpponentPlayerID)
	assert.Equal(GameStatusOnGoing, game.Status)
	assert.NotNil(game.SelfBoard)
	assert.NotNil(game.OpponentBoard)
	assert.NotEqual(PlayerNone, game.PlayerTurn)
	assert.Equal(PlayerNone, game.PlayerWon)
}

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
