package pkg

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestSpaceshipFromPatternValidSimple(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)
	assert.NotNil(spaceship)
}

func TestSpaceshipFromPatternValidWinger(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternWinger)
	assert.NoError(err)
	assert.NotNil(spaceship)

	// @TODO: better assert than using fmt.Sprintf
	assert.Equal("[0x0 2x0 0x1 2x1 1x2 0x3 2x3 0x4 2x4]", fmt.Sprintf("%s", spaceship.coords))
}

func TestSpaceshipFromPatternValidAngle(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternAngle)
	assert.NoError(err)
	assert.NotNil(spaceship)
}

func TestSpaceshipFromPatternValidAClass(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternAClass)
	assert.NoError(err)
	assert.NotNil(spaceship)
}

func TestSpaceshipFromPatternValidBClass(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternBClass)
	assert.NoError(err)
	assert.NotNil(spaceship)
}

func TestSpaceshipFromPatternValidSClass(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)
	assert.NotNil(spaceship)
}

func TestSpaceshipFromPatternInvalidBlank(t *testing.T) {
	assert := require.New(t)

	_, err := SpaceshipFromPattern([]string{})
	assert.Error(err)

	_, err = SpaceshipFromPattern([]string{
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

func TestSpaceshipFromPatternInvalidRows(t *testing.T) {
	assert := require.New(t)

	_, err := SpaceshipFromPattern([]string{
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
		"*...............",
	})
	assert.Error(err)
}

func TestSpaceshipFromPatternInvalidCols(t *testing.T) {
	assert := require.New(t)

	_, err := SpaceshipFromPattern([]string{
		"*****************",
	})

	assert.Error(err)
}

func TestSpaceshipFromPatternInvalidChars(t *testing.T) {
	assert := require.New(t)

	_, err := SpaceshipFromPattern([]string{"**A"})
	assert.Error(err)

	_, err = SpaceshipFromPattern([]string{"**1"})
	assert.Error(err)

	_, err = SpaceshipFromPattern([]string{"**€"})
	assert.Error(err)
}
