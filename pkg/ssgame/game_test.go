package ssgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	assert := require.New(t)

	game, err := CreateNewGame("match-1-1", &Player{
		PlayerID: "player-1",
		FullName: "Player 1",
	}, true)

	assert.NoError(err)
	assert.Equal("player-1", game.Opponent.PlayerID)
	assert.Equal("Player 1", game.Opponent.FullName)
	assert.Equal(GameStatusOnGoing, game.Status)
	assert.NotNil(game.SelfBoard)
	assert.NotNil(game.OpponentBoard)
	assert.Equal(PlayerSelf, game.PlayerTurn)
	assert.Equal(PlayerNone, game.PlayerWon)
}

func TestInitNewGame(t *testing.T) {
	assert := require.New(t)

	game, err := InitNewGame("game-1", &Player{
		PlayerID: "player-1",
		FullName: "Player 1",
	}, PlayerSelf)

	assert.NoError(err)
	assert.Equal("player-1", game.Opponent.PlayerID)
	assert.Equal("Player 1", game.Opponent.FullName)
	assert.Equal(GameStatusOnGoing, game.Status)
	assert.NotNil(game.SelfBoard)
	assert.NotNil(game.OpponentBoard)
	assert.Equal(PlayerSelf, game.PlayerTurn)
	assert.Equal(PlayerNone, game.PlayerWon)
}
