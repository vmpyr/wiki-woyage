package player

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

func HandleCreatePlayer(conn *websocket.Conn, playerName string) {
	player, err := CreatePlayer(conn, playerName)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type:       "player_created",
		PlayerID:   player.PlayerID,
		PlayerName: player.PlayerName,
	})
}
