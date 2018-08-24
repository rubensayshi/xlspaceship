package ssclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rubensayshi/xlspaceship/pkg/ssgame"

	"sync"

	_ "github.com/rubensayshi/xlspaceship/statik" // registers our static files to serve
)

func Serve(s *XLSpaceship, port int, wg *sync.WaitGroup) {
	r := mux.NewRouter()

	// add go routing handlers
	AddWhoAmIGameHandler(s, r)
	AddNewGameHandler(s, r)
	AddInitGameHandler(s, r)
	AddGameStatusHandler(s, r)
	AddReceiveSalvoHandler(s, r)
	AddFireSalvoHandler(s, r)

	// add static file handler
	ServeAddStaticHandler(r)

	// start serving
	wg.Add(1)
	go func() {
		fmt.Printf("Serve REST API on :%d \n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
			panic(err)
		}

		wg.Done()
	}()
}

func AddWhoAmIGameHandler(s *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s \n", r.Method, r.RequestURI)

		res := &WhoAmIResponse{
			UserID:   s.Player.PlayerID,
			FullName: s.Player.FullName,
			Games:    make([]string, 0, len(s.games)),
		}

		for gameID, _ := range s.games {
			res.Games = append(res.Games, gameID)
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
	r.HandleFunc("/xl-spaceship/user/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
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
		salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Coords invalid")))
			return
		}

		// fire off the salvo
		res, alreadyFinished, err := s.FireSalvo(game, salvo)
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
		if alreadyFinished {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
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
		salvo, err := ssgame.CoordsGroupFromSalvoStrings(req.Salvo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Coords invalid")))
			return
		}

		// check if it's the newOpponent's turn, otherwise he's not allowed to fire
		if game.PlayerTurn != ssgame.PlayerOpponent {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Not your turn")))
			return
		}

		// process the incoming salvo
		res, alreadyFinished, err := s.ReceiveSalvo(game, salvo)
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
		if alreadyFinished {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write(resJson)
	})
}
