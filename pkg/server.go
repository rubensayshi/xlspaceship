package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Serve(s *XLSpaceship) {
	mux := http.NewServeMux()

	AddPingHandler(s, mux)
	AddNewGameHandler(s, mux)
	AddGameStatusHandler(s, mux)
	AddReceiveSalvoHandler(s, mux)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", DEFAULT_PORT), mux); err != nil {
		panic(err)
	}
}

func AddPingHandler(s *XLSpaceship, mux *http.ServeMux) {
	mux.HandleFunc(URI_PREFIX+"/protocol/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}

func AddNewGameHandler(s *XLSpaceship, mux *http.ServeMux) {
	mux.HandleFunc(URI_PREFIX+"/protocol/game/new", func(w http.ResponseWriter, r *http.Request) {
		req := &NewGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		game, err := s.NewGame(req.UserID, req.FullName, req.SpaceshipProtocol.Hostname, req.SpaceshipProtocol.Port)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: %s", err)))
			return
		}

		fmt.Printf("new game:\n %s \n", game)

		res := NewGameResponseFromGame(s, game)

		resJson, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJson)
	})
}

func AddGameStatusHandler(s *XLSpaceship, mux *http.ServeMux) {
	mux.HandleFunc(URI_PREFIX+"/user/game/", func(w http.ResponseWriter, r *http.Request) {
		// @TODO: abstract
		uriChunks := strings.Split(r.RequestURI, "/")
		gameID := uriChunks[len(uriChunks)-1]

		// @TODO: is this possible?
		if gameID == "game" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		game, ok := s.games[gameID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Game not found")))
			return
		}

		fmt.Printf("game:\n %s \n", game)

		res := GameStatusResponseFromGame(s, game)

		resJson, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resJson)
	})
}

func AddReceiveSalvoHandler(s *XLSpaceship, mux *http.ServeMux) {
	mux.HandleFunc(URI_PREFIX+"/protocol/game/", func(w http.ResponseWriter, r *http.Request) {
		req := &ReceiveSalvoRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		// @TODO: abstract
		uriChunks := strings.Split(r.RequestURI, "/")
		gameID := uriChunks[len(uriChunks)-1]

		if gameID == "game" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		game, ok := s.games[gameID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Game not found")))
			return
		}

		fmt.Printf("game:\n %s \n", game)

		// parse salvo into coords
		// @TODO: test for what happens when out of bounds
		salvo := make(CoordsGroup, len(req.Salvo))
		for i, coordsStr := range req.Salvo {
			coords, err := CoordsFromString(coordsStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Coords invalid")))
				return
			}

			salvo[i] = coords
		}

		if game.Status == GameStatusInitializing {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Game hasn't been started yet")))
			return
		}

		if game.Status == GameStatusDone {
			salvoRes := make([]*ShotResult, len(salvo))
			for i, shot := range salvo {
				salvoRes[i] = &ShotResult{
					Coords:     shot,
					ShotStatus: ShotStatusMiss,
				}
			}

			res := ReceiveSalvoResponseFromSalvoResult(salvoRes, s, game)

			resJson, err := json.MarshalIndent(res, "", "    ")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(resJson)
			return
		}

		if game.PlayerTurn != PlayerOpponent {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Not your turn")))
			return
		}

		salvoRes := game.SelfBoard.ReceiveSalvo(salvo)
		game.PlayerTurn = PlayerSelf

		win := true
		for _, spaceship := range game.SelfBoard.spaceships {
			fmt.Printf("spaceship dead? %v \n%s \n", spaceship.dead, spaceship.coords)
			if !spaceship.dead {
				win = false
				break
			}
		}

		if win {
			game.Status = GameStatusDone
			game.PlayerWon = PlayerOpponent
		}

		fmt.Printf("game:\n %s \n", game)

		res := ReceiveSalvoResponseFromSalvoResult(salvoRes, s, game)

		resJson, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resJson)
	})
}
