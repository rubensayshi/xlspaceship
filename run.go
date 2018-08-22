package main

import (
	"flag"

	"os"
	"strconv"

	"fmt"

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

func maybeGetEnv(env string, dflt string) string {
	val := os.Getenv(env)
	if val == "" {
		return dflt
	}

	return val
}

// @TODO: DEFAULT -> SEE DOC
var fPort = flag.Int("port", maybeGetEnvInt("PORT", 8000), "port to serve the REST API on")
var fPlayerID = flag.String("playerID", maybeGetEnv("PLAYERID", "player-1"), "your player ID")

func main() {
	fmt.Printf("main \n")
	flag.Parse()

	s := ssclient.NewXLSpaceship(*fPlayerID, "localhost", *fPort)

	ssclient.Serve(s, *fPort)
}
