XL Spaceship
============
This is the XL Spaceship game, writting in Go with a Angular 1.4 GUI.

You can find the latest release on github and download the binaries if you just want to play: https://github.com/rubensayshi/xlspaceship/releases

Below are instructions for developers to develop and build the project, you can find my notes about the project in [NOTES.md](NOTES.md)

As desired there's a `run.sh` and `build.sh` but they're just wrappers around what is descibed below.

Develop
-------
#### Requirements
 - go 1.10+
 - Glide for deps management (https://github.com/Masterminds/glide)
 - nodejs / npm

#### Go Deps Install
```
glide update
```

#### GUI Deps Install
```
npm install -g gulp bower
cd ./gui && npm install && cd ..
```

#### Tests
```
make tests
```

or for coverage report:
```
make coverage
```

#### Build
```
make build
```

will generate binaries in `./bin/` for linux 64/386 and windows 64/386 with all the statics bundled in the binary.


#### Run
```
go run main.go --playerId player-1234 --playerName "Player 1" --port 8080
```

if you want to run without building first, ofcourse to play against "another player" you should run a 2nd instance on a different port as well;
```
go run main.go --playerId player-4321 --playerName "Player 2" --port 8090
```


#### Livereload for Go
Get `gin` (https://github.com/codegangsta/gin) to make it easy to restart the process when you make code changes and run like:
```
DONTOPENGUI=true PLAYERID=player-8080 PLAYERNAME="Player 8080" gin --port 8080 --appPort 8081 run main.go
```

Then run another instance for the 2nd player:
```
DONTOPENGUI=true PLAYERID=player-8080 PLAYERNAME="Player 8080" gin --port 8080 --appPort 8081 run main.go
```

#### Livereload for GUI
By default it will serve statics from filesystem instead of bundling the statics into the binary, so you can actually use the livereload chrome extension if you run the watch:
```
make watch-build-gui
```
