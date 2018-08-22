package ssgame

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
