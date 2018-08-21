package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	assert := require.New(t)

	game, err := NewGame("player-1")

	assert.NoError(err)
	assert.Equal("player-1", game.OpponentPlayerID)
	assert.Equal(GameStatusOnGoing, game.Status)
	assert.NotNil(game.SelfBoard)
	assert.NotNil(game.OpponentBoard)
	assert.NotEqual(PlayerNone, game.PlayerTurn)
	assert.Equal(PlayerNone, game.PlayerWon)
}
