package pkg

const (
	DEFAULT_PORT = "3001"
	URI_PREFIX   = "/xl-spaceship"
)

type NewGameRequestSpaceshipProtocol struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type NewGameRequest struct {
	UserID            string                          `json:"user_id"`
	FullName          string                          `json:"full_name"`
	SpaceshipProtocol NewGameRequestSpaceshipProtocol `json:"spaceship_protocol"`
}

type NewGameResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	GameID   string `json:"game_id"`
	Starting string `json:"starting"`
}

func NewGameResponseFromGame(s *XLSpaceship, game *Game) *NewGameResponse {
	res := &NewGameResponse{}

	res.UserID = s.PlayerID
	res.FullName = s.PlayerName
	res.GameID = game.GameID

	if game.PlayerTurn == PlayerSelf {
		res.Starting = s.PlayerID
	} else {
		res.Starting = game.OpponentPlayerID
	}

	return res
}
