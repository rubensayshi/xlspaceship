package pkg

const (
	DEFAULT_PORT = "3001"
	URI_PREFIX   = "/xl-spaceship"
)

type GameRequestSpaceshipProtocol struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type NewGameRequest struct {
	UserID            string                       `json:"user_id"`
	FullName          string                       `json:"full_name"`
	SpaceshipProtocol GameRequestSpaceshipProtocol `json:"spaceship_protocol"`
}

type NewGameResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	GameID   string `json:"game_id"`
	Starting string `json:"starting"`
}

func NewGameResponseFromGame(s *XLSpaceship, game *Game) *NewGameResponse {
	res := &NewGameResponse{}

	res.UserID = s.Player.PlayerID
	res.FullName = s.Player.FullName
	res.GameID = game.GameID

	if game.PlayerTurn == PlayerSelf {
		res.Starting = s.Player.PlayerID
	} else {
		res.Starting = game.Opponent.PlayerID
	}

	return res
}

type InitGameRequest struct {
	SpaceshipProtocol GameRequestSpaceshipProtocol `json:"spaceship_protocol"`
}

type GameStatusResponse struct {
	Self     GameStatusResponsePlayer `json:"self"`
	Opponent GameStatusResponsePlayer `json:"opponent"`
	Game     GameStatusResponseGame   `json:"game"`
}

type GameStatusResponsePlayer struct {
	UserID string   `json:"user_id"`
	Board  []string `json:"board"`
}

type GameStatusResponseGame struct {
	PlayerTurn string `json:"player_turn"`
}

func GameStatusResponseFromGame(s *XLSpaceship, game *Game) *GameStatusResponse {
	res := &GameStatusResponse{}

	res.Self = GameStatusResponsePlayer{
		Board: game.SelfBoard.ToPattern(),
	}

	res.Opponent = GameStatusResponsePlayer{
		Board: game.OpponentBoard.ToPattern(),
	}

	return res
}

type ReceiveSalvoRequest struct {
	Salvo []string `json:"salvo"`
}

type SalvoResponse struct {
	Salvo map[string]string `json:"salvo"`
	Game  interface{}       `json:"game"` // ewww interface, but alternative is having multiple structs for this response
}

type ReceiveSalvoResponseGame struct {
	PlayerTurn string `json:"player_turn"`
}

type ReceiveSalvoWonResponseGame struct {
	Won string `json:"won"`
}

func ReceiveSalvoResponseFromSalvoResult(salvoResult []*ShotResult, s *XLSpaceship, game *Game) *SalvoResponse {
	res := &SalvoResponse{
		Salvo: make(map[string]string, len(salvoResult)),
	}

	for _, shotResult := range salvoResult {
		res.Salvo[shotResult.Coords.String()] = shotResult.ShotStatus.String()
	}

	if game.Status == GameStatusDone {
		gameRes := ReceiveSalvoWonResponseGame{}
		gameRes.Won = game.Opponent.PlayerID

		res.Game = gameRes
	} else {
		gameRes := ReceiveSalvoResponseGame{}
		if game.PlayerTurn == PlayerSelf {
			gameRes.PlayerTurn = s.Player.PlayerID
		} else {
			gameRes.PlayerTurn = game.Opponent.PlayerID
		}

		res.Game = gameRes
	}

	return res
}
