package ssgame

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

type TestCoordFixture struct {
	input string
	x     int8
	y     int8
}

func TestCoordFromString(t *testing.T) {
	assert := require.New(t)

	fixtures := []TestCoordFixture{
		{"0x0", 0, 0},
	}

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			fixtures = append(fixtures, TestCoordFixture{fmt.Sprintf("%xx%x", x, y), int8(x), int8(y)})
			fixtures = append(fixtures, TestCoordFixture{fmt.Sprintf("%Xx%X", x, y), int8(x), int8(y)})
		}
	}

	for _, fixture := range fixtures {
		coord, err := CoordsFromString(fixture.input)
		assert.NoError(err)
		assert.Equal(fixture.x, coord.x)
		assert.Equal(fixture.y, coord.y)
	}
}

func TestCoordFromStringInvalid(t *testing.T) {
	assert := require.New(t)

	fixtures := []TestCoordFixture{
		{"00x0", 0, 0},
		{"0x00", 0, 0},
		{"01x0", 0, 0},
		{"Gx0", 0, 0},
		{"-1x0", 0, 0},
		{"0xx0", 0, 0},
	}

	for _, fixture := range fixtures {
		_, err := CoordsFromString(fixture.input)
		assert.Error(err)
	}
}
