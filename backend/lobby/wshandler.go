package lobby

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
)

func HandleCreateLobby(p *st.Player) {
	lobbyID, err := CreateLobby(p)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type:    "lobby_created",
		LobbyID: lobbyID,
	})
}

func HandleJoinLobby(p *st.Player, lobbyID string) {
	err := JoinLobby(p, lobbyID)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type: "lobby_joined",
	})
}

func HandleLeaveLobby(p *st.Player, lobbyID string) {
	err := LeaveLobby(p, lobbyID, false)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type: "lobby_left",
	})
}

func HandleToggleReady(p *st.Player, lobbyID string) {
	toggledTo, err := ToggleReady(p.PlayerID, lobbyID)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	var responseType string
	if toggledTo {
		responseType = "toggle_ready"
	} else {
		responseType = "toggle_not_ready"
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type: responseType,
	})
}
