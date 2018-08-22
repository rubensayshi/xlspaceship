package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/pkg/errors"
)

type Requester interface {
	NewGame(dest SpaceshipProtocol, req *NewGameRequest) (*NewGameResponse, error)
	ReceiveSalvo(dest SpaceshipProtocol, gameID string, req *ReceiveSalvoRequest) (*SalvoResponse, error)
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

func (r *HttpRequester) ReceiveSalvo(dest SpaceshipProtocol, gameID string, req *ReceiveSalvoRequest) (*SalvoResponse, error) {
	fmt.Printf("1")
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}

	fmt.Printf("2")
	res, err := Put(fmt.Sprintf("http://%s:%d/xl-spaceship/protocol/game/%s", dest.Hostname, dest.Port, gameID), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to request receive salvo (http: %d)", res.StatusCode)
	}
	fmt.Printf("3")
	defer res.Body.Close()

	fmt.Printf("4")
	dbg, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("DBG: %s \n", dbg)

	salvoResponse := &SalvoResponse{}
	err = json.NewDecoder(res.Body).Decode(salvoResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to request receive salvo")
	}

	return salvoResponse, nil
}
