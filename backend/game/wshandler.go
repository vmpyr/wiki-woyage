package game

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
)

func HandleStartNewGame(p *st.Player, lobbyID string, gameSettings st.GameSettings) {
	gameID, err := StartNewGame(p, lobbyID, gameSettings)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type:   "game_started",
		GameID: gameID,
	})
}

func HandleJoinGame(p *st.Player, gameID string) {
	err := JoinGame(p, gameID)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type: "game_joined",
	})
}

func HandleVoteToEnd(p *st.Player, gameID string) {
	err := VoteToEndGame(p, gameID)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type:       "voted_to_end_game",
		PlayerName: p.PlayerName,
	})
}
