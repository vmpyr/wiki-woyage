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

func HandleJoinLobby(conn *websocket.Conn, lobbyID string, playerID string) {
	err := JoinLobby(lobbyID, playerID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: "lobby_joined",
	})
}

func HandleLeaveLobby(conn *websocket.Conn, lobbyID string, playerID string) {
	err := LeaveLobby(lobbyID, playerID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: "lobby_left",
	})
}
