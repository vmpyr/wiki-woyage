package structs

type WebSocketMessage struct {
	Type         string       `json:"type"`
	PlayerName   string       `json:"playerName,omitempty"`
	PlayerID     string       `json:"playerID,omitempty"`
	LobbyID      string       `json:"lobbyID,omitempty"`
	GameID       string       `json:"gameID,omitempty"`
	GameSettings GameSettings `json:"gameSettings,omitempty"`
}

type WebSocketResponse struct {
	Type         string `json:"type"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	PlayerName   string `json:"playerName,omitempty"`
	PlayerID     string `json:"playerID,omitempty"`
	LobbyID      string `json:"lobbyID,omitempty"`
	GameID       string `json:"gameID,omitempty"`
}
