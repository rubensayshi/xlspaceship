package ssgame

import (
	"testing"

	"fmt"

	"math/rand"

	"strings"

	"github.com/stretchr/testify/require"
)

func TestNewRandomBoard(t *testing.T) {
	assert := require.New(t)

	_, err := NewRandomBoard(SpaceshipsSetForBaseGame)
	assert.NoError(err)
}

func TestShotStatusFromString(t *testing.T) {
	assert := require.New(t)

	ss, err := ShotStatusFromString("hit")
	assert.NoError(err)
	assert.Equal(ShotStatusHit, ss)

	ss, err = ShotStatusFromString("miss")
	assert.NoError(err)
	assert.Equal(ShotStatusMiss, ss)

	ss, err = ShotStatusFromString("kill")
	assert.NoError(err)
	assert.Equal(ShotStatusKill, ss)
}

func TestShotStatusFromStringInvalid(t *testing.T) {
	assert := require.New(t)

	_, err := ShotStatusFromString("oops")
	assert.Error(err)

	_, err = ShotStatusFromString("hitmiss")
	assert.Error(err)

	_, err = ShotStatusFromString("")
	assert.Error(err)
}

func TestShotStatus_String(t *testing.T) {
	assert := require.New(t)

	assert.Equal("miss", fmt.Sprintf("%s", ShotStatusMiss))
	assert.Equal("hit", fmt.Sprintf("%s", ShotStatusHit))
	assert.Equal("kill", fmt.Sprintf("%s", ShotStatusKill))
}

func TestCoordsState_String(t *testing.T) {
	assert := require.New(t)

	assert.Equal(".", fmt.Sprintf("%s", CoordsBlank))
	assert.Equal("X", fmt.Sprintf("%s", CoordsHit))
	assert.Equal("-", fmt.Sprintf("%s", CoordsMiss))
	assert.Equal("*", fmt.Sprintf("%s", CoordsShip))
}

func TestCoordsGroupFromSalvoStrings(t *testing.T) {
	assert := require.New(t)

	cg, err := CoordsGroupFromSalvoStrings([]string{"0x0", "AxA"})
	assert.NoError(err)
	assert.Equal("0x0", cg[0].String())
	assert.Equal("AxA", cg[1].String())
}

func TestCoordsGroupFromSalvoStringsInvalid(t *testing.T) {
	assert := require.New(t)

	_, err := CoordsGroupFromSalvoStrings([]string{"Gx0", "AxA"})
	assert.Error(err)
}

func TestNewRandomBoardTooMany1(t *testing.T) {
	assert := require.New(t)

	// fresh seed so the "randomness" is consistent
	rand.Seed(1)

	ManySpaceships := [][]string{
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
	}

	// first board should fail with this seed
	board, err := newRandomBoard(ManySpaceships)
	assert.NoError(err)
	assert.Nil(board)

	// second board should also fail with this seed
	board, err = newRandomBoard(ManySpaceships)
	assert.NoError(err)
	assert.Nil(board)

	// third board should pass with this seed
	board, err = newRandomBoard(ManySpaceships)
	assert.NoError(err)
	assert.NotNil(board)
}

func TestNewRandomBoardTooMany2(t *testing.T) {
	assert := require.New(t)

	// fresh seed so the "randomness" is consistent
	rand.Seed(1)

	ManySpaceships := [][]string{
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
		SpaceshipPatternSClass,
	}

	board, err := NewRandomBoard(ManySpaceships)
	assert.NoError(err)
	assert.NotNil(board)
}

func TestSpaceship_String(t *testing.T) {
	assert := require.New(t)

	ss, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)

	assert.Equal(strings.Join(SpaceshipPatternSClass, "\n"), fmt.Sprintf("%s", ss))

}
