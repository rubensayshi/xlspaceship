package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Serve(s *XLSpaceship) {
	mux := http.NewServeMux()

	mux.HandleFunc(URI_PREFIX+"/protocol/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc(URI_PREFIX+"/protocol/game/new", func(w http.ResponseWriter, r *http.Request) {
		req := &NewGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		game := s.NewGame(req.UserID, req.FullName, req.SpaceshipProtocol.Hostname, req.SpaceshipProtocol.Port)

		fmt.Printf("new game %v \n", game)

		res := NewGameResponseFromGame(s, game)

		resJson, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJson)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", DEFAULT_PORT), mux); err != nil {
		panic(err)
	}
}
