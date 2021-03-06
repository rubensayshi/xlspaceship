package ssgame

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/require"
)

func TestBoard_ToPatternBoard(t *testing.T) {
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

	board := NewSelfBoard()
	err := FillBoardFromPattern(board.BaseBoard, pattern)
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(8, board.CountHits())
	assert.Equal(8, board.CountMisses())

	assert.Equal(pattern, board.ToPattern())

	assert.Equal(strings.Join(pattern, "\n"), board.String())
}

func TestBoardToPatternWithShipsNoMarks(t *testing.T) {
	assert := require.New(t)

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
		"....*...........",
		"..**............",
	}

	board, err := NewBlankSelfBoard()
	assert.NoError(err)
	assert.NotNil(board)
	assert.Equal(0, board.CountHits())
	assert.Equal(0, board.CountMisses())
	assert.Equal(0, len(board.spaceships))

	assert.Equal(BlankBoardPattern(), board.ToPattern())

	winger, _ := SpaceshipFromPattern(SpaceshipPatternWinger)
	angle, _ := SpaceshipFromPattern(SpaceshipPatternAngle)
	aclass, _ := SpaceshipFromPattern(SpaceshipPatternAClass)
	bclass, _ := SpaceshipFromPattern(SpaceshipPatternBClass)
	sclass, _ := SpaceshipFromPattern(SpaceshipPatternSClass)

	assert.NoError(board.AddSpaceshipOnCoords(winger.CopyWithOffset(1, 0)))
	assert.NoError(board.AddSpaceshipOnCoords(angle.CopyWithOffset(10, 0)))
	assert.NoError(board.AddSpaceshipOnCoords(aclass.CopyWithOffset(1, 6)))
	assert.NoError(board.AddSpaceshipOnCoords(bclass.CopyWithOffset(10, 6)))
	assert.NoError(board.AddSpaceshipOnCoords(sclass.CopyWithOffset(1, 11)))

	assert.Equal(expectedPattern, board.ToPattern())
}

func TestBoardToPatternWithShipsAndMarks(t *testing.T) {
	assert := require.New(t)

	pattern := []string{
		".X-X........--..",
		".X-X........--..",
		"..X.........--..",
		"...X............",
		"....-...........",
		".....-..........",
		"......-.........",
		".......-........",
		"........-.......",
		".........-......",
		"..........X.....",
		"...........-....",
		"............-...",
		".............-..",
		"..............-.",
		"...............-",
	}

	expectedPattern := []string{
		".X-X......*.--..",
		".X-X......*.--..",
		"..X.......*.--..",
		".*.X......***...",
		".*.*-...........",
		".....-..........",
		"..*...-...**....",
		".*.*...-..*.*...",
		".***....-.**....",
		".*.*.....-*.*...",
		"..........X*....",
		"..**.......-....",
		".*..........-...",
		"..**.........-..",
		"....*.........-.",
		"..**...........-",
	}

	board := NewSelfBoard()
	err := FillBoardFromPattern(board.BaseBoard, pattern)
	assert.NoError(err)
	assert.NotNil(board)

	assert.Equal(pattern, board.ToPattern())

	winger, _ := SpaceshipFromPattern(SpaceshipPatternWinger)
	angle, _ := SpaceshipFromPattern(SpaceshipPatternAngle)
	aclass, _ := SpaceshipFromPattern(SpaceshipPatternAClass)
	bclass, _ := SpaceshipFromPattern(SpaceshipPatternBClass)
	sclass, _ := SpaceshipFromPattern(SpaceshipPatternSClass)

	assert.NoError(board.AddSpaceshipOnCoords(winger.CopyWithOffset(1, 0)))
	assert.NoError(board.AddSpaceshipOnCoords(angle.CopyWithOffset(10, 0)))
	assert.NoError(board.AddSpaceshipOnCoords(aclass.CopyWithOffset(1, 6)))
	assert.NoError(board.AddSpaceshipOnCoords(bclass.CopyWithOffset(10, 6)))
	assert.NoError(board.AddSpaceshipOnCoords(sclass.CopyWithOffset(1, 11)))

	assert.Equal(expectedPattern, board.ToPattern())
}
