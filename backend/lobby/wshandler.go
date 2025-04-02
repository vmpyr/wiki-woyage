package lobby

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

func HandleCreateLobby(conn *websocket.Conn, playerID string) {
	lobby, err := CreateLobby(conn, playerID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type:    "lobby_created",
		LobbyID: lobby.LobbyID,
	})
}

func HandleJoinLobby(conn *websocket.Conn, playerID string, lobbyID string) {
	err := JoinLobby(playerID, lobbyID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: "lobby_joined",
	})
}

func HandleLeaveLobby(conn *websocket.Conn, playerID string, lobbyID string) {
	err := LeaveLobby(playerID, lobbyID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: "lobby_left",
	})
}

func HandleToggleReady(conn *websocket.Conn, playerID string, lobbyID string) {
	toggledTo, err := ToggleReady(playerID, lobbyID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	var responseType string
	if toggledTo {
		responseType = "toggle_ready"
	} else {
		responseType = "toggle_not_ready"
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: responseType,
	})
}
