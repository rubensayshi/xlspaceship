package ssclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sync"

	"github.com/gorilla/mux"

	_ "github.com/rubensayshi/xlspaceship/statik" // registers our static files to serve
)

func Serve(xl *XLSpaceship, port int, wg *sync.WaitGroup) {
	r := mux.NewRouter()

	// add go routing handlers
	AddWhoAmIGameHandler(xl, r)
	AddNewGameHandler(xl, r)
	AddInitGameHandler(xl, r)
	AddGameStatusHandler(xl, r)
	AddReceiveSalvoHandler(xl, r)
	AddFireSalvoHandler(xl, r)

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

func AddWhoAmIGameHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s \n", r.Method, r.RequestURI)

		req := &WhoAmIRequest{}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get game status: %s", xlRes.err)))
			return
		}

		res, ok := xlRes.res.(*WhoAmIResponse)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get game status: invalid response type: %T", xlRes.res)))
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

func AddNewGameHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/protocol/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%xl: %xl \n", r.Method, r.RequestURI)

		req := &NewGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: %s", xlRes.err)))
			return
		}

		res, ok := xlRes.res.(*NewGameResponse)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: invalid response type: %T", xlRes.res)))
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

func AddInitGameHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%xl \n", r.RequestURI)

		req := &InitGameRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to init game: %s", xlRes.err)))
			return
		}

		gameID, ok := xlRes.res.(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to create game: invalid response type: %T", xlRes.res)))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Location", fmt.Sprintf("/xl-spaceship/user/game/%s", gameID))
		w.WriteHeader(http.StatusSeeOther)
		w.Write([]byte(fmt.Sprintf("A new game has been created at xl-spaceship/user/game/%s", gameID)))
	})
}

func AddGameStatusHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%xl \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		req := &GameStatusRequest{GameID: gameID}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get game status: %s", xlRes.err)))
			return
		}

		res, ok := xlRes.res.(*GameStatusResponse)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get game status: invalid response type: %T", xlRes.res)))
			return
		}

		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Game not found")))
		} else {
			resJson, err := json.MarshalIndent(res, "", "    ")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resJson)
		}
	})
}

func AddFireSalvoHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/user/game/{gameID}/fire", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%xl \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		req := &FireSalvoRequest{
			GameID: gameID,
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to fire salvo: %s", xlRes.err)))
			return
		}

		res, ok := xlRes.res.(*SalvoResponse)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to fire salvo: invalid response type: %T", xlRes.res)))
			return
		}

		resJson, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if res.AlreadyFinished {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write(resJson)
	})
}

func AddReceiveSalvoHandler(xl *XLSpaceship, r *mux.Router) {
	r.HandleFunc("/xl-spaceship/protocol/game/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%xl \n", r.RequestURI)

		vars := mux.Vars(r)
		gameID := vars["gameID"]

		req := &ReceiveSalvoRequest{GameID: gameID}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad JSON"))
			return
		}

		xlRes := xl.HandleRequest(req)
		if xlRes.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to receive salvo: %s", xlRes.err)))
			return
		}

		res, ok := xlRes.res.(*SalvoResponse)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to receive salvo: invalid response type: %T", xlRes.res)))
			return
		}

		resJson, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if res.AlreadyFinished {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write(resJson)
	})
}
