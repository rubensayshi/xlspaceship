package ssclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rubensayshi/xlspaceship/pkg/ssgame"
)

func Serve(s *XLSpaceship, port int) {
	r := mux.NewRouter()

	AddNewGameHandler(s, r)
	AddInitGameHandler(s, r)
	AddGameStatusHandler(s, r)
	AddReceiveSalvoHandler(s, r)
	AddFireSalvoHandler(s, r)

	fmt.Printf("Serve :%d \n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		panic(err)
	}
}

func AddNewGameHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/protocol/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s \n", r.Method, r.RequestURI)

		req := &NewGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		res, err := s.NewGame(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: %s", err)))
			return
		}

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

func AddInitGameHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		req := &InitGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		gameID, err := s.InitNewGame(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed to init game: %s", err)))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Location", fmt.Sprintf("/xl-spaceship/user/game/%s", gameID))
		w.WriteHeader(http.StatusSeeOther)
		w.Write([]byte(fmt.Sprintf("A new game has been created at xl-spaceship/user/game/%s", gameID)))
	})
}

func AddGameStatusHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		res, ok := s.GameStatus(gameID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Game not found")))
			return
		}

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

func AddFireSalvoHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/{gameID}/fire", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		req := &ReceiveSalvoRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		// check if the game exists
		game, ok := s.games[gameID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Game not found")))
			return
		}

		// parse salvo into coords
		// @TODO: test for what happens when out of bounds
		salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Coords invalid")))
			return
		}

		// @TODO
		if game.Status == ssgame.GameStatusDone {
			panic("done")
		}

		// fire off the salvo
		res, err := s.FireSalvo(game, salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}

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

func AddReceiveSalvoHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/protocol/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		req := &ReceiveSalvoRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		// check if game exists
		game, ok := s.games[gameID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Game not found")))
			return
		}

		// parse salvo into coords
		// @TODO: test for what happens when out of bounds
		salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Coords invalid")))
			return
		}

		// if the game is already done then we create a mock response with misses
		if game.Status == ssgame.GameStatusDone {
			salvoRes := make([]*ssgame.ShotResult, len(salvo))
			for i, shot := range salvo {
				salvoRes[i] = &ssgame.ShotResult{
					Coords:     shot,
					ShotStatus: ssgame.ShotStatusMiss,
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

		// check if it's the opponent's turn, otherwise he's not allowed to fire
		if game.PlayerTurn != ssgame.PlayerOpponent {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Not your turn")))
			return
		}

		// process the incoming salvo
		res, err := s.ReceiveSalvo(game, salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed to receive salvo")))
		}

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
