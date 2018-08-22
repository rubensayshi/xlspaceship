package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve(s *XLSpaceship, port int) {

	r := mux.NewRouter()

	AddPingHandler(s, r)
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

func AddPingHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/protocol/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}

func AddNewGameHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/protocol/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		req := &NewGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		fmt.Printf("%s: %s \n", r.RequestURI, req)

		game, err := s.NewGame(req.UserID, req.FullName, req.SpaceshipProtocol.Hostname, req.SpaceshipProtocol.Port)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: %s", err)))
			return
		}

		fmt.Printf("new game:\n%s \n", game)

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

func AddInitGameHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/user/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		req := &InitGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		game, err := s.InitNewGame(req.SpaceshipProtocol.Hostname, req.SpaceshipProtocol.Port)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed to init game: %s", err)))
			return
		}

		fmt.Printf("init game:\n%s \n", game)

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Location", fmt.Sprintf("/xl-spaceship/user/game/%s", game.GameID))
		w.WriteHeader(http.StatusSeeOther)
		w.Write([]byte(fmt.Sprintf("A new game has been created at xl-spaceship/user/game/%s", game.GameID)))
	})
}

func AddGameStatusHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/user/game/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

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

func AddFireSalvoHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/user/game/{gameID}/fire", func(w http.ResponseWriter, r *http.Request) {
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

		// @TODO
		if game.Status == GameStatusDone {
			panic("done")
		}

		// @TODO: should we remove this check and rely on the response?
		if game.PlayerTurn != PlayerSelf {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Not your turn")))
			return
		}

		salvoRes, err := s.FireSalvo(game, salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		game.PlayerTurn = PlayerOpponent

		// @TODO
		win := false
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

func AddReceiveSalvoHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc(URI_PREFIX+"/protocol/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
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
