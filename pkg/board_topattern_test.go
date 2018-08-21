package pkg

import (
	"testing"

	"fmt"

	"strings"

	"github.com/stretchr/testify/require"
)

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

func TestBoardToPatternWithShips(t *testing.T) {
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

	expectedPattern := []string{
		".*.*......*.....",
		".*.*......*.....",
		"..*.......*.....",
		".*.*......***...",
		".*.*............",
		"................",
		"..*.......**....",
		".*.*......*.*...",
		".***......**....",
		".*.*......*.*...",
		"..........**....",
		"..**............",
		".*..............",
		"..**............",
		"...*............",
		".**.............",
	}

	board, err := BoardFromPattern(pattern)
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(0, len(board.hits))
	assert.Equal(0, len(board.misses))
	assert.Equal(0, len(board.spaceships))

	assert.Equal(pattern, board.ToPattern())

	winger, _ := SpaceshipFromPattern(SpaceshipPatternWinger)
	angle, _ := SpaceshipFromPattern(SpaceshipPatternAngle)
	aclass, _ := SpaceshipFromPattern(SpaceshipPatternAClass)
	bclass, _ := SpaceshipFromPattern(SpaceshipPatternBClass)
	sclass, _ := SpaceshipFromPattern(SpaceshipPatternSClass)

	assert.NoError(board.AddSpaceshipOnCoords(winger, 1, 0))
	assert.NoError(board.AddSpaceshipOnCoords(angle, 10, 0))
	assert.NoError(board.AddSpaceshipOnCoords(aclass, 1, 6))
	assert.NoError(board.AddSpaceshipOnCoords(bclass, 10, 6))
	assert.NoError(board.AddSpaceshipOnCoords(sclass, 1, 11))

	fmt.Printf("%s \n", strings.Join(board.ToPattern(), "\n"))
	assert.Equal(expectedPattern, board.ToPattern())
}
