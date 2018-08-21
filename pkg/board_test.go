package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRandomBoard(t *testing.T) {
	assert := require.New(t)

	_, err := NewRandomBoard()
	assert.NoError(err)
}
