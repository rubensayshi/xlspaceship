package main

import (
	"flag"

	"os"
	"strconv"

	"fmt"

	"strings"

	"sync"
	"time"

	"github.com/pkg/browser"
	"github.com/rubensayshi/xlspaceship/pkg/ssclient"
)

func maybeGetEnvInt(env string, dflt int) int {
	val := os.Getenv(env)
	if val == "" {
		return dflt
	}

	intval, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(intval)
}

func maybeGetEnvBool(env string, dflt bool) bool {
	val := os.Getenv(env)
	if val == "" {
		return dflt
	}

	switch strings.ToLower(val) {
	case "true":
		return true
	case "1":
		return true
	case "yes":
		return true
	case "false":
		return true
	case "0":
		return true
	case "no":
		return true
	}

	return dflt
}

func maybeGetEnv(env string, dflt string) string {
	val := os.Getenv(env)
	if val == "" {
		return dflt
	}

	return val
}

// define flags for CLI
var fPort = flag.Int("port", maybeGetEnvInt("PORT", 8080), "port to serve the REST API on")
var fPlayerID = flag.String("playerID", maybeGetEnv("PLAYERID", "player-1"), "your player ID")
var fPlayerName = flag.String("playerName", maybeGetEnv("PLAYERNAME", "Player 1"), "your player name")
var fCheat = flag.Bool("cheat", maybeGetEnvBool("CHEAT", false), "enable cheat mode")
var fDontOpenGui = flag.Bool("dontopengui", maybeGetEnvBool("DONTOPENGUI", false), "don't pop open the GUI when the process starts")

func main() {
	fmt.Printf("XLSpaceship starting ... \n")
	flag.Parse()

	s := ssclient.NewXLSpaceship(*fPlayerID, *fPlayerName, "localhost", *fPort)
	if *fCheat {
		s.EnableCheatMode()
	}

	wg := &sync.WaitGroup{}

	ssclient.Serve(s, *fPort, wg)

	guiUrl := fmt.Sprintf("http://localhost:%d/gui/game.html", *fPort)

	if !*fDontOpenGui {
		fmt.Printf("Opening GUI in browser (if it does not open visit: %s", guiUrl)
		go func() {
			<-time.After(time.Millisecond * 100)
			browser.OpenURL(guiUrl)
		}()
	} else {
		fmt.Printf("You can open the GUI in the browser by visiting: %s", guiUrl)
	}

	// if waitgroup finishes then we quit
	wg.Wait()
}
