package main

import (
	"flag"

	"os"
	"strconv"

	"fmt"

	"strings"

	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/rubensayshi/xlspaceship/pkg/ssclient"
)

// helper to get int from `env` var or use default
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

// helper to get bool from `env` var (parsing string to bool) or use default
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

// helper to get string from `env` var or use default
func maybeGetEnv(env string, dflt string) string {
	val := os.Getenv(env)
	if val == "" {
		return dflt
	}

	return val
}

// define flags for CLI, most with env var fallback
var fPort = flag.Int("port", maybeGetEnvInt("PORT", 8080), "port to serve the REST API on")
var fPlayerID = flag.String("playerID", maybeGetEnv("PLAYERID", ""), "your player ID")
var fPlayerName = flag.String("playerName", maybeGetEnv("PLAYERNAME", ""), "your player name")
var fCheat = flag.Bool("cheat", maybeGetEnvBool("CHEAT", false), "enable cheat mode")
var fDontOpenGui = flag.Bool("dontopengui", maybeGetEnvBool("DONTOPENGUI", false), "don't pop open the GUI when the process starts")

func maybePromptPlayerID() {
	for *fPlayerID == "" {
		prompt := promptui.Prompt{
			Label: "Player ID",
		}

		result, err := prompt.Run()
		if err != nil {
			panic(err)
		}

		*fPlayerID = result
	}
}

func maybePromptPlayerName() {
	for *fPlayerName == "" {
		prompt := promptui.Prompt{
			Label: "Player Name",
		}

		result, err := prompt.Run()
		if err != nil {
			panic(err)
		}

		*fPlayerName = result
	}
}

func main() {
	fmt.Printf("XLSpaceship starting ... \n")
	flag.Parse()

	// prompt for player ID and name if they're not provided already
	maybePromptPlayerID()
	maybePromptPlayerName()

	// init the main controller of the game
	s := ssclient.NewXLSpaceship(*fPlayerID, *fPlayerName, "localhost", *fPort)
	// enable cheat mode if configured
	if *fCheat {
		s.EnableCheatMode()
	}

	// create wg that will control when we exit
	wg := &sync.WaitGroup{}

	go func() {
		s.Run()
	}()

	// serve the rest API
	ssclient.Serve(s, *fPort, wg)

	// open or print the gui URL
	guiUrl := fmt.Sprintf("http://localhost:%d/gui/game.html", *fPort)
	if !*fDontOpenGui {
		fmt.Printf("Opening GUI in browser (if it does not open visit: %s\n", guiUrl)
		go func() {
			<-time.After(time.Millisecond * 100)
			browser.OpenURL(guiUrl)
		}()
	} else {
		fmt.Printf("You can open the GUI in the browser by visiting: %s\n", guiUrl)
	}

	// if waitgroup finishes then we quit
	wg.Wait()
}
