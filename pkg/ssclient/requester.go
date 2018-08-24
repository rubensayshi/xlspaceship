package ssclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Requester interface {
	NewGame(dest SpaceshipProtocol, req *NewGameRequest) (*NewGameResponse, error)
	ReceiveSalvo(dest SpaceshipProtocol, req *ReceiveSalvoRequest) (*SalvoResponse, error)
}

type HttpRequester struct {
}

func (r *HttpRequester) NewGame(dest SpaceshipProtocol, req *NewGameRequest) (*NewGameResponse, error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request new game")
	}

	res, err := http.Post(fmt.Sprintf("http://%s:%d/xl-spaceship/protocol/game/new", dest.Hostname, dest.Port), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request new game")
	}
	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Failed to request new game (http: %d)", res.StatusCode)
	}
	defer res.Body.Close()

	newGameRes := &NewGameResponse{}
	err = json.NewDecoder(res.Body).Decode(newGameRes)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request new game")
	}

	return newGameRes, nil
}

func (r *HttpRequester) ReceiveSalvo(dest SpaceshipProtocol, req *ReceiveSalvoRequest) (*SalvoResponse, error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}

	res, err := Put(fmt.Sprintf("http://%s:%d/xl-spaceship/protocol/game/%s", dest.Hostname, dest.Port, req.GameID), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to request receive salvo (http: %d)", res.StatusCode)
	}
	defer res.Body.Close()

	salvoResponse := &SalvoResponse{}
	err = json.NewDecoder(res.Body).Decode(salvoResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}

	err = salvoResponse.Normalize()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}

	return salvoResponse, nil
}
