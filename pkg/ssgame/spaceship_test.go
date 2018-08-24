package ssgame

import (
	"testing"

	"fmt"

	"strings"

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

func TestSpaceship_ToPatternSimple(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)
	assert.NotNil(spaceship)

	assert.Equal([]string{"***"}, spaceship.ToPattern())
}

func TestSpaceship_ToPatternSClass(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)
	assert.NotNil(spaceship)

	assert.Equal(SpaceshipPatternSClass, spaceship.ToPattern())
}

func TestSpaceship_RotateSimple90(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)
	assert.NotNil(spaceship)

	assert.Equal([]string{"***"}, spaceship.ToPattern())

	spaceship.rotate(90)

	assert.Equal([]string{
		"*",
		"*",
		"*",
	}, spaceship.ToPattern())
}

func TestSpaceship_RotateSclass90(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)
	assert.NotNil(spaceship)

	spaceship.rotate(90)

	assert.Equal([]string{
		"...*",
		"*.*.*",
		"*.*.*",
		".*",
	}, spaceship.ToPattern())

}

func TestSpaceship_RotateSimple180(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)
	assert.NotNil(spaceship)

	assert.Equal([]string{"***"}, spaceship.ToPattern())

	spaceship.rotate(180)

	assert.Equal([]string{
		"***",
	}, spaceship.ToPattern())
}

func TestSpaceship_RotateSclass180(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)

	spaceship.rotate(180)

	assert.Equal([]string{
		".**",
		"...*",
		".**",
		"*",
		".**",
	}, spaceship.ToPattern())

}

func TestSpaceship_RotateSimple270(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern([]string{
		"***",
	})
	assert.NoError(err)
	assert.NotNil(spaceship)

	assert.Equal([]string{"***"}, spaceship.ToPattern())

	spaceship.rotate(270)

	assert.Equal([]string{
		"*",
		"*",
		"*",
	}, spaceship.ToPattern())
}

func TestSpaceship_RotateSclass270(t *testing.T) {
	assert := require.New(t)

	spaceship, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)
	assert.NotNil(spaceship)

	spaceship.rotate(270)

	assert.Equal([]string{
		"...*",
		"*.*.*",
		"*.*.*",
		".*",
	}, spaceship.ToPattern())

}

func TestSpaceship_String(t *testing.T) {
	assert := require.New(t)

	ss, err := SpaceshipFromPattern(SpaceshipPatternSClass)
	assert.NoError(err)

	assert.Equal(strings.Join(SpaceshipPatternSClass, "\n"), fmt.Sprintf("%s", ss))
}
