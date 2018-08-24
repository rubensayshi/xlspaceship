package ssclient

import (
	"net/http"

	"math/rand"

	"io"

	"fmt"

	"github.com/pkg/errors"
	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
)

type XLRequest struct {
	req     interface{}
	resChan chan *XLResponse
}

type XLResponse struct {
	res interface{}
	err error
}

type XLSpaceship struct {
	Player      *ssgame.Player
	games       map[string]*ssgame.Game
	requester   Requester
	cheat       bool
	reqQueue    chan *XLRequest
	matchIDIncr uint
}

func NewXLSpaceship(playerID string, playerName string, host string, port int) *XLSpaceship {
	s := &XLSpaceship{
		Player: &ssgame.Player{
			PlayerID:     playerID,
			FullName:     playerName,
			ProtocolHost: host,
			ProtocolPort: port,
		},
		games:     make(map[string]*ssgame.Game),
		requester: &HttpRequester{},
		reqQueue:  make(chan *XLRequest, 1),
	}

	// make a seed based on the playerID, that way it's deterministic but different per player
	//  this is purely for easy debugging because this way every time we restart the game all random things will be the same
	var seed int64 = 0
	for _, char := range playerID {
		seed += int64(char)
	}
	rand.Seed(seed)

	return s
}

func (xl *XLSpaceship) EnableCheatMode() {
	xl.cheat = true
}

func (xl *XLSpaceship) NewGameID() string {
	xl.matchIDIncr++
	return fmt.Sprintf("match-%s-%d", xl.Player.PlayerID, xl.matchIDIncr)
}

func (xl *XLSpaceship) Run() {
	for xlReq := range xl.reqQueue {
		switch xlReq.req.(type) {
		case *WhoAmIRequest:
			res, err := xl.WhoAmIRequest(xlReq.req.(*WhoAmIRequest))
			xlReq.resChan <- &XLResponse{res, err}

		case *NewGameRequest:
			res, err := xl.NewGameRequest(xlReq.req.(*NewGameRequest))
			xlReq.resChan <- &XLResponse{res, err}

		case *InitGameRequest:
			res, err := xl.InitNewGameRequest(xlReq.req.(*InitGameRequest))
			xlReq.resChan <- &XLResponse{res, err}

		case *GameStatusRequest:
			res, err := xl.GameStatusRequest(xlReq.req.(*GameStatusRequest))
			xlReq.resChan <- &XLResponse{res, err}

		case *ReceiveSalvoRequest:
			res, err := xl.ReceiveSalvoRequest(xlReq.req.(*ReceiveSalvoRequest))
			xlReq.resChan <- &XLResponse{res, err}

		case *FireSalvoRequest:
			res, err := xl.FireSalvoRequest(xlReq.req.(*FireSalvoRequest))
			xlReq.resChan <- &XLResponse{res, err}

		default:
			panic(fmt.Sprintf("Invalid request type: %T", xlReq.req))
		}

	}
}

func (xl *XLSpaceship) HandleRequest(req interface{}) *XLResponse {
	resChan := make(chan *XLResponse)

	xl.reqQueue <- &XLRequest{
		req:     req,
		resChan: resChan,
	}

	return <-resChan
}

func (xl *XLSpaceship) WhoAmIRequest(req *WhoAmIRequest) (*WhoAmIResponse, error) {
	res := &WhoAmIResponse{
		UserID:   xl.Player.PlayerID,
		FullName: xl.Player.FullName,
		Games:    make([]string, 0, len(xl.games)),
	}

	for gameID, _ := range xl.games {
		res.Games = append(res.Games, gameID)
	}

	return res, nil
}

// handle a NewGameRequest from another player
func (xl *XLSpaceship) NewGameRequest(req *NewGameRequest) (*NewGameResponse, error) {
	opponent := &ssgame.Player{
		PlayerID:     req.UserID,
		FullName:     req.FullName,
		ProtocolHost: req.SpaceshipProtocol.Hostname,
		ProtocolPort: req.SpaceshipProtocol.Port,
	}

	if xl.Player.PlayerID == opponent.PlayerID || xl.Player.FullName == opponent.FullName {
		return nil, errors.Errorf("Failed to create new game: opponent has same user_id or fullname as player")
	}

	game, err := ssgame.CreateNewGame(xl.NewGameID(), opponent, xl.cheat)
	if err != nil {
		return nil, err
	}

	xl.games[game.GameID] = game

	res := NewGameResponseFromGame(xl, game)
	return res, nil
}

// send a NewGameRequest to another player
func (xl *XLSpaceship) InitNewGameRequest(req *InitGameRequest) (string, error) {
	newGameReq := &NewGameRequest{
		UserID:            xl.Player.PlayerID,
		FullName:          xl.Player.FullName,
		SpaceshipProtocol: SpaceshipProtocol{xl.Player.ProtocolHost, xl.Player.ProtocolPort},
	}

	newGameRes, err := xl.requester.NewGame(req.SpaceshipProtocol, newGameReq)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to init new game")
	}

	firstPlayer := ssgame.PlayerSelf
	if newGameRes.Starting != xl.Player.PlayerID {
		firstPlayer = ssgame.PlayerOpponent
	}

	opponent := &ssgame.Player{
		PlayerID:     newGameRes.UserID,
		FullName:     newGameRes.FullName,
		ProtocolHost: req.SpaceshipProtocol.Hostname,
		ProtocolPort: req.SpaceshipProtocol.Port,
	}

	game, err := ssgame.InitNewGame(newGameRes.GameID, opponent, firstPlayer)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to init new game")
	}

	xl.games[game.GameID] = game

	return game.GameID, nil
}

// retrieve the GameStatusResponse for a game
func (xl *XLSpaceship) GameStatusRequest(req *GameStatusRequest) (*GameStatusResponse, error) {
	game, ok := xl.games[req.GameID]
	if !ok {
		return nil, nil
	}

	res := GameStatusResponseFromGame(xl, game)

	return res, nil
}

func (xl *XLSpaceship) gameStatus(gameID string) (*GameStatusResponse, bool) {
	game, ok := xl.games[gameID]
	if !ok {
		return nil, false
	}

	res := GameStatusResponseFromGame(xl, game)

	return res, true
}

// receive a salvo from another player
func (xl *XLSpaceship) ReceiveSalvoRequest(req *ReceiveSalvoRequest) (*SalvoResponse, error) {
	// check if game exists
	game, ok := xl.games[req.GameID]
	if !ok {
		return nil, errors.Errorf("Game not found")
	}

	// parse salvo into coords
	salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
	if err != nil {
		return nil, errors.Errorf("Coords invalid")
	}

	// check if it's the opponent's turn, otherwise he's not allowed to fire
	// @TODO: alreadyFinished check should go before this
	if game.PlayerTurn != ssgame.PlayerOpponent {
		return nil, errors.Errorf("Not your turn")
	}

	// process the incoming salvo
	res, alreadyFinished, err := xl.receiveSalvo(game, salvo)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to receive salvo")
	}

	res.AlreadyFinished = alreadyFinished

	return res, nil
}

// receive a salvo from another player
func (xl *XLSpaceship) receiveSalvo(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, bool, error) {
	// check that we're not cheating
	if !xl.cheat && len(salvo) > game.OpponentBoard.CountShipsAlive() {
		return nil, false, errors.Errorf("More shots than ships alive (%d)", game.OpponentBoard.CountShipsAlive())
	}

	// if the game is already done then we create a mock response with misses
	if game.Status == ssgame.GameStatusDone {
		res, err := xl.ReceiveSalvoGameFinished(game, salvo)
		if err != nil {
			return nil, false, errors.Wrapf(err, "Failed to fire salvo")
		}

		return res, true, nil
	}

	salvoRes := game.SelfBoard.ReceiveSalvo(salvo)
	game.PlayerTurn = ssgame.PlayerSelf

	if game.SelfBoard.AllShipsDead() {
		game.Status = ssgame.GameStatusDone
		game.PlayerWon = ssgame.PlayerOpponent
	}

	return SalvoResponseFromSalvoResult(salvoRes, xl, game), false, nil
}

// build a SalvoResponse for when a game is already finished
func (xl *XLSpaceship) ReceiveSalvoGameFinished(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, error) {
	salvoRes := make([]*ssgame.ShotResult, len(salvo))
	for i, shot := range salvo {
		salvoRes[i] = &ssgame.ShotResult{
			Coords:     shot,
			ShotStatus: ssgame.ShotStatusMiss,
		}
	}

	return SalvoResponseFromSalvoResult(salvoRes, xl, game), nil
}

// receive a salvo from another player
func (xl *XLSpaceship) FireSalvoRequest(req *FireSalvoRequest) (*SalvoResponse, error) {
	// check if game exists
	game, ok := xl.games[req.GameID]
	if !ok {
		return nil, errors.Errorf("Game not found")
	}

	// parse salvo into coords
	salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
	if err != nil {
		return nil, errors.Errorf("Coords invalid")
	}

	// check if it's self's turn, otherwise he's not allowed to fire
	// @TODO: alreadyFinished check should go before this
	if game.PlayerTurn != ssgame.PlayerSelf {
		return nil, errors.Errorf("Not your turn")
	}

	// fire off the salvo
	res, alreadyFinished, err := xl.fireSalvo(game, salvo)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fire salvo")
	}

	res.AlreadyFinished = alreadyFinished

	return res, nil
}

// send a salvo to another player
func (xl *XLSpaceship) fireSalvo(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, bool, error) {
	// check that we're not cheating
	if !xl.cheat && len(salvo) > game.SelfBoard.CountShipsAlive() {
		return nil, false, errors.Errorf("More shots than ships alive (%d)", game.SelfBoard.CountShipsAlive())
	}

	// if the game is already done then we create a mock response with misses
	if game.Status == ssgame.GameStatusDone {
		res, err := xl.FireSalvoGameFinished(game, salvo)
		if err != nil {
			return nil, false, errors.Wrapf(err, "Failed to fire salvo")
		}

		return res, true, nil
	}

	req := &ReceiveSalvoRequest{
		GameID: game.GameID,
		Salvo:  make([]string, len(salvo)),
	}

	for i, salvo := range salvo {
		req.Salvo[i] = salvo.String()
	}

	res, err := xl.requester.ReceiveSalvo(SpaceshipProtocol{
		Hostname: game.Opponent.ProtocolHost,
		Port:     game.Opponent.ProtocolPort,
	}, req)
	if err != nil {
		return nil, false, errors.Wrapf(err, "Failed to fire salvo (req)")
	}

	// mark result on our end
	salvoRes := make([]*ssgame.ShotResult, 0, len(res.Salvo))
	for coordsStr, shotResStr := range res.Salvo {
		coords, err := ssgame.CoordsFromString(coordsStr)
		if err != nil {
			return nil, false, errors.Wrapf(err, "Failed to fire salvo")
		}

		shotStatus, err := ssgame.ShotStatusFromString(shotResStr)
		if err != nil {
			return nil, false, errors.Wrapf(err, "Failed to fire salvo")
		}

		salvoRes = append(salvoRes, &ssgame.ShotResult{coords, shotStatus})

		game.OpponentBoard.ApplyShotStatus(coords, shotStatus)
	}

	game.PlayerTurn = ssgame.PlayerOpponent

	if res.GameWon != nil {
		game.Status = ssgame.GameStatusDone
		game.PlayerWon = ssgame.PlayerSelf
	}

	return SalvoResponseFromSalvoResult(salvoRes, xl, game), false, nil
}

// build a SalvoResponse for when a game is already finished
func (xl *XLSpaceship) FireSalvoGameFinished(game *ssgame.Game, salvo ssgame.CoordsGroup) (*SalvoResponse, error) {
	salvoRes := make([]*ssgame.ShotResult, len(salvo))
	for i, shot := range salvo {
		salvoRes[i] = &ssgame.ShotResult{
			Coords:     shot,
			ShotStatus: ssgame.ShotStatusMiss,
		}
	}

	return SalvoResponseFromSalvoResult(salvoRes, xl, game), nil
}

// Helper function to do PUT requests because http builtin only has helpers for GET and POST >_>
func Put(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}
