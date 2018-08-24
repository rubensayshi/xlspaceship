package ssclient

import (
	"testing"

	"fmt"

	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
	"github.com/stretchr/testify/require"
)

func TestXLSpaceshipFullGameOld(t *testing.T) {
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

	// doneChan is used to know when the goroutines are done
	doneChan := make(chan error)

	// helper channels to switch between player turns
	xl1GoChan := make(chan bool)
	xl2GoChan := make(chan bool)

	// let the handlers run
	go func() {
		xl1.Run()
	}()
	go func() {
		xl2.Run()
	}()

	// xl1 goroutine
	go func() {
		xl := xl1

		// when we're done we send done signal and we close the xl2 channel so that it won't wait for us to signal anymore
		defer func() {
			doneChan <- nil
			close(xl2GoChan)
		}()

		more := true
		nextTurn := func() bool {
			fmt.Printf("player 1 done...\n")
			if !more {
				return true
			}

			// tell xl2 to go, don't block when it's closed
			select {
			case xl2GoChan <- true:
			default:
			}

			// wait until we're allowed to go
			_, _more := <-xl1GoChan
			fmt.Printf("player 1...\n")
			if !_more {
				more = false
			}

			return true
		}

		// wait for when we're allowed to begin our game
		<-xl1GoChan

		_, err := xl.InitNewGameRequest(&InitGameRequest{
			SpaceshipProtocol{
				Hostname: xl2.Player.ProtocolHost,
				Port:     xl2.Player.ProtocolPort,
			},
		})
		assert.NoError(err)

		assert.Equal(1, len(xl.games))
		game := xl.games["match-testplayer-2-1"]
		assert.NotNil(game)

		// we know opponent will cheat to make himself first
		assert.Equal(ssgame.PlayerOpponent, game.PlayerTurn)

		nextTurn()

		assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

		salvo1Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
			GameID: game.GameID,
			Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
		})
		assert.NoError(err)
		assert.Nil(salvo1Res.GameWon)

		nextTurn()

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
	}()

	// xl2 goroutine
	go func() {
		xl := xl2

		// when we're done we send done signal and we close the xl2 channel so that it won't wait for us to signal anymore
		defer func() {
			doneChan <- nil
			close(xl1GoChan)
		}()
		more := true
		nextTurn := func() bool {
			fmt.Printf("player 2 done...\n")
			if !more {
				return true
			}

			// tell xl2 to go, don't block when it's closed
			select {
			case xl1GoChan <- true:
			default:
			}

			// wait until we're allowed to go
			_, _more := <-xl2GoChan
			fmt.Printf("player 2...\n")
			if !_more {
				more = false
			}

			return true
		}

		// wait for when we're allowed to begin our game
		<-xl2GoChan

		assert.Equal(1, len(xl.games))
		game := xl.games["match-testplayer-2-1"]
		assert.NotNil(game)

		assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

		salvo1Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
			GameID: game.GameID,
			Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
		})
		assert.NoError(err)
		assert.Nil(salvo1Res.GameWon)

		nextTurn()

		assert.Equal(ssgame.PlayerSelf, game.PlayerTurn)

		salvo2Res, err := xl.FireSalvoRequest(&FireSalvoRequest{
			GameID: game.GameID,
			Salvo:  []string{"0x1", "0x2", "0x3", "0x4", "0x5"},
		})
		assert.NoError(err)
		assert.Nil(salvo2Res.GameWon)
	}()

	// gogo
	xl1GoChan <- true

	err := <-doneChan
	assert.NoError(err)

	err = <-doneChan
	assert.NoError(err)
}
