package ssclient

import (
	"testing"

	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
	"github.com/stretchr/testify/require"
)

// @TODO: ReceiveSalvo test

func mustCoordsFromString(coordsStr string) *ssgame.Coords {
	coords, err := ssgame.CoordsFromString(coordsStr)
	if err != nil {
		panic(err)
	}

	return coords
}

func TestNewXLSpaceship(t *testing.T) {
	assert := require.New(t)

	xl := NewXLSpaceship("testplayer-1", "Test Player 1", "notlocalhost", 1337)
	assert.NotNil(xl)
	assert.Equal("testplayer-1", xl.Player.PlayerID)
	assert.Equal("notlocalhost", xl.Player.ProtocolHost)
	assert.Equal(1337, xl.Player.ProtocolPort)
}

func TestXLSpaceship_NewGame(t *testing.T) {
	assert := require.New(t)

	xl := NewXLSpaceship("testplayer-1", "Test Player 1", "notlocalhost", 1337)
	assert.NotNil(xl)

	req := &NewGameRequest{
		UserID:   "testplayer-2",
		FullName: "Test Player 2",
		SpaceshipProtocol: SpaceshipProtocol{
			Hostname: "notlocalhost2",
			Port:     6666,
		},
	}

	res, err := xl.NewGame(req)
	assert.NoError(err)
	assert.NotNil(res)

	assert.Equal("testplayer-1", res.UserID)
	assert.Equal("testplayer-1", res.FullName) // @TODO
}

func TestXLSpaceship_InitNewGame(t *testing.T) {
	assert := require.New(t)

	xl := NewXLSpaceship("testplayer-1", "Test Player 1", "notlocalhost", 1337)
	assert.NotNil(xl)

	mockRequester := &MockRequester{}
	xl.requester = mockRequester

	ssProtocol := SpaceshipProtocol{
		Hostname: "notlocalhost2",
		Port:     6666,
	}
	req := &InitGameRequest{
		SpaceshipProtocol: ssProtocol,
	}

	mockRequester.On("NewGame", ssProtocol, NewGameRequest{
		UserID:   "testplayer-1",
		FullName: "testplayer-1",
		SpaceshipProtocol: SpaceshipProtocol{
			Hostname: "notlocalhost",
			Port:     1337,
		},
	}).Return(&NewGameResponse{}, nil)

	res, err := xl.InitNewGame(req)
	assert.NoError(err)
	assert.NotNil(res)

	mockRequester.AssertExpectations(t)
}

func TestXLSpaceship_FireSalvo(t *testing.T) {
	assert := require.New(t)

	xl := NewXLSpaceship("testplayer-1", "Test Player 1", "notlocalhost", 1337)
	assert.NotNil(xl)

	mockRequester := &MockRequester{}
	xl.requester = mockRequester

	ssProtocol := SpaceshipProtocol{
		Hostname: "notlocalhost2",
		Port:     6666,
	}

	newGameRes, err := xl.NewGame(&NewGameRequest{
		UserID:            "testplayer-2",
		SpaceshipProtocol: ssProtocol,
	})
	assert.NoError(err)
	assert.NotNil(newGameRes)

	game := xl.games[newGameRes.GameID]

	mockRequester.On("ReceiveSalvo", ssProtocol, game.GameID, ReceiveSalvoRequest{
		Salvo: []string{"0x0", "1x1"},
	}).Return(&SalvoResponse{}, nil)

	res, err := xl.FireSalvo(game, ssgame.CoordsGroup{
		mustCoordsFromString("0x0"),
		mustCoordsFromString("1x1"),
	})
	assert.NoError(err)
	assert.NotNil(res)

	mockRequester.AssertExpectations(t)
}
