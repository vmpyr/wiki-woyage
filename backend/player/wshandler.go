package player

import (
	st "wiki-woyage/structs"

	"github.com/gorilla/websocket"
)

func SendError(conn *websocket.Conn, errorMessage string) {
	response := st.WebSocketResponse{
		Type:         "error",
		ErrorMessage: errorMessage,
	}
	err := conn.WriteJSON(response)
	if err != nil {
		conn.Close()
	}
}

func SendResponse(conn *websocket.Conn, responseType string, playerID string, playerName string) {
	response := st.WebSocketResponse{
		Type:       responseType,
		PlayerID:   playerID,
		PlayerName: playerName,
	}
	err := conn.WriteJSON(response)
	if err != nil {
		conn.Close()
	}
}

func HandleCreatePlayer(conn *websocket.Conn, playerName string) {
	player, err := CreatePlayer(conn, playerName)
	if err != nil {
		SendError(conn, err.Error())
		return
	}

	SendResponse(conn, "player_created", player.PlayerID, player.PlayerName)
}
