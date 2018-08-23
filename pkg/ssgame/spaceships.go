package ssgame

// this file contains the spaceship patterns
//  and the set of spaceships used for a game

var SpaceshipPatternWinger = []string{
	"*.*",
	"*.*",
	".*.",
	"*.*",
	"*.*",
}

var SpaceshipPatternAngle = []string{
	"*",
	"*",
	"*",
	"***",
}

var SpaceshipPatternAClass = []string{
	".*.",
	"*.*",
	"***",
	"*.*",
}

var SpaceshipPatternBClass = []string{
	"**",
	"*.*",
	"**",
	"*.*",
	"**",
}

var SpaceshipPatternSClass = []string{
	".**",
	"*",
	".**",
	"...*",
	".**",
}

var SpaceshipsSetForBaseGame = [][]string{
	SpaceshipPatternWinger,
	SpaceshipPatternAngle,
	SpaceshipPatternAClass,
	SpaceshipPatternBClass,
	SpaceshipPatternSClass,
}
