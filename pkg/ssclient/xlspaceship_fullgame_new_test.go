package ssclient

import (
	"testing"

	"fmt"

	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
	"github.com/stretchr/testify/require"
)

func TestXLSpaceshipFullGameNew(t *testing.T) {
	assert := require.New(t)

	xl1 := NewXLSpaceship("testplayer-1", "Test Player 1", "notlocalhost", 1337)
	assert.NotNil(xl1)
	xl2 := NewXLSpaceship("testplayer-2", "Test Player 2", "notlocalhost", 1338)
	assert.NotNil(xl2)

	xl1.EnableCheatMode()
	xl2.EnableCheatMode()

	reqChan1 := make(chan *XLRequest, 1)
	reqChan2 := make(chan *XLRequest, 1)

	memRequester1 := &MemRequester{reqChan1}
	memRequester2 := &MemRequester{reqChan2}

	xl1.reqQueue = reqChan1
	xl2.reqQueue = reqChan2
	xl1.requester = memRequester2
	xl2.requester = memRequester1

	// let the handlers run
	go func() {
		xl1.Run()
	}()
	go func() {
		xl2.Run()
	}()

	gameID := "match-testplayer-2-1"

	turns := []func(xl *XLSpaceship){
		// player 1
		func(xl *XLSpaceship) {
			_, err := xl.InitNewGameRequest(&InitGameRequest{
				SpaceshipProtocol{
					Hostname: xl2.Player.ProtocolHost,
					Port:     xl2.Player.ProtocolPort,
				},
			})
			assert.NoError(err)

			assert.Equal(1, len(xl.games))
			game := xl.games[gameID]
			assert.NotNil(game)

			// we know opponent will cheat to make himself first
			assert.Equal(ssgame.PlayerOpponent, game.PlayerTurn)
		},
		// player 2
		func(xl *XLSpaceship) {
			game := xl.games[gameID]

			assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

			salvo1Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
				GameID: game.GameID,
				Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
			})
			assert.NoError(err)
			assert.Nil(salvo1Res.GameWon)
		},
		// player 1
		func(xl *XLSpaceship) {
			game := xl.games[gameID]

			assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

			salvo1Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
				GameID: game.GameID,
				Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
			})
			assert.NoError(err)
			assert.Nil(salvo1Res.GameWon)
		},
		// player 2
		func(xl *XLSpaceship) {
			game := xl.games[gameID]

			assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

			salvo1Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
				GameID: game.GameID,
				Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
			})
			assert.NoError(err)
			assert.Nil(salvo1Res.GameWon)
		},
		// player 1
		func(xl *XLSpaceship) {
			game := xl.games[gameID]

			assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

			salvo3 := make([]string, 0, 16*16)
			for x := 0; x < 16; x++ {
				for y := 0; y < 16; y++ {
					salvo3 = append(salvo3, fmt.Sprintf("%Xx%X", x, y))
				}
			}
			salvo3Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
				GameID: game.GameID,
				Salvo:  salvo3,
			})
			assert.NoError(err)
			assert.NotNil(salvo3Res.GameWon)
			assert.Equal(xl.Player.PlayerID, salvo3Res.GameWon.Won)
		},
	}

	xl1Turn := true

	for _, turn := range turns {
		xl := xl1
		if !xl1Turn {
			xl = xl2
		}

		turn(xl)

		xl1Turn = !xl1Turn
	}
}
