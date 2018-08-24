package ssclient

import (
	"github.com/pkg/errors"
)

type MemRequester struct {
	reqChan chan *XLRequest
}

func (r *MemRequester) NewGame(dest SpaceshipProtocol, req *NewGameRequest) (*NewGameResponse, error) {
	resChan := make(chan *XLResponse)

	r.reqChan <- &XLRequest{
		req:     req,
		resChan: resChan,
	}

	xlRes := <-resChan

	res, ok := xlRes.res.(*NewGameResponse)
	if !ok {
		return nil, errors.Errorf("Failed to request new game: Invalid response type: %T", res)
	}

	return res, nil
}

func (r *MemRequester) ReceiveSalvo(dest SpaceshipProtocol, req *ReceiveSalvoRequest) (*SalvoResponse, error) {
	resChan := make(chan *XLResponse)

	r.reqChan <- &XLRequest{
		req:     req,
		resChan: resChan,
	}

	xlRes := <-resChan

	res, ok := xlRes.res.(*SalvoResponse)
	if !ok {
		return nil, errors.Errorf("Failed to request new game: Invalid response type: %T", res)
	}

	return res, nil
}
