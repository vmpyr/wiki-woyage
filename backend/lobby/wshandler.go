package lobby

import (
	"log"
	"wiki-woyage/player"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

func CheckPlayerExists(playerID string) bool {
	if _, err := player.GetPlayer(playerID); err != nil {
		log.Printf("Player %s not found", playerID)
		return false
	}
	return true
}

func HandleCreateLobby(conn *websocket.Conn, playerID string) {
	if !CheckPlayerExists(playerID) {
		utils.SendError(conn, "playerID not found")
		return
	}

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
	if !CheckPlayerExists(playerID) {
		utils.SendError(conn, "playerID not found")
		return
	}

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
	if !CheckPlayerExists(playerID) {
		utils.SendError(conn, "playerID not found")
		return
	}

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
	if !CheckPlayerExists(playerID) {
		utils.SendError(conn, "playerID not found")
		return
	}

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
