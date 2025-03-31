package structs

type WebSocketMessage struct {
	Type       string `json:"type"`
	PlayerName string `json:"playerName,omitempty"`
	LobbyID    string `json:"lobbyID,omitempty"`
	PlayerID   string `json:"playerID,omitempty"`
}
