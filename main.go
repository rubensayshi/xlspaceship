package main

import (
	"flag"

	"os"
	"strconv"

	"fmt"

	"strings"

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

var fPort = flag.Int("port", maybeGetEnvInt("PORT", 8080), "port to serve the REST API on")
var fPlayerID = flag.String("playerID", maybeGetEnv("PLAYERID", "player-1"), "your player ID")
var fPlayerName = flag.String("playerName", maybeGetEnv("PLAYERNAME", "Player 1"), "your player name")
var fCheat = flag.Bool("cheat", maybeGetEnvBool("CHEAT", false), "enable cheat mode")

func main() {
	fmt.Printf("main \n")
	flag.Parse()

	s := ssclient.NewXLSpaceship(*fPlayerID, *fPlayerName, "localhost", *fPort)
	if *fCheat {
		s.EnableCheatMode()
	}

	ssclient.Serve(s, *fPort)
}
