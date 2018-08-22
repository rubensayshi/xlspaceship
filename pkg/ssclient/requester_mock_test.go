package ssclient

import (
	"github.com/stretchr/testify/mock"
)

type MockRequester struct {
	mock.Mock
}

func (r *MockRequester) NewGame(dest SpaceshipProtocol, req *NewGameRequest) (*NewGameResponse, error) {
	args := r.Called(dest, *req)

	return args.Get(0).(*NewGameResponse), args.Error(1)
}

func (r *MockRequester) ReceiveSalvo(dest SpaceshipProtocol, gameID string, req *ReceiveSalvoRequest) (*SalvoResponse, error) {
	args := r.Called(dest, gameID, *req)

	return args.Get(0).(*SalvoResponse), args.Error(1)
}
