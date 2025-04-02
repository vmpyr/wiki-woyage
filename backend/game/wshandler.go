package game

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

func HandleStartNewGame(conn *websocket.Conn, playerID, lobbyID string, gameSettings st.GameSettings) {
	game, err := StartNewGame(conn, playerID, lobbyID, gameSettings)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type:   "game_started",
		GameID: game.GameID,
	})
}

func HandleJoinGame(conn *websocket.Conn, playerID, gameID string) {
	err := JoinGame(playerID, gameID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type: "game_joined",
	})
}

func HandleVoteToEnd(conn *websocket.Conn, playerID, gameID string) {
	err := VoteToEndGame(playerID, gameID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	utils.SendResponse(conn, st.WebSocketResponse{
		Type:     "voted_to_end_game",
		PlayerID: playerID,
	})
}
