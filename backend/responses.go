package main

const (
	ResponseLobbyJoined = "lobby_joined"
	ResponsePlayerList  = "player_list"
	ResponseNewPlayer   = "new_player"
	ResponseAmIAdmin    = "am_i_admin"
)

type Response struct {
	Type     string `json:"type"`
	Response any    `json:"response"`
}

type LobbyJoinedResponse struct {
	LobbyID  string `json:"lobbyID"`
	Username string `json:"username"`
}

type PlayerListResponse struct {
	Players []string `json:"players"`
}

type NewPlayerResponse struct {
	Username string `json:"username"`
}

type AmIAdminResponse struct {
	IsAdmin bool `json:"isAdmin"`
}
