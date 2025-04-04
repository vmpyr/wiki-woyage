package player

import (
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
)

func HandleAddPlayerName(p *st.Player, playerName string) {
	err := AddPlayerName(p, playerName)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type:       "player_created",
		PlayerID:   p.PlayerID,
		PlayerName: p.PlayerName,
	})
}

func HandleDeletePlayer(p *st.Player) {
	err := DeletePlayer(p)
	if err != nil {
		utils.SendError(p.Conn, err.Error())
		return
	}

	utils.SendResponse(p.Conn, st.WebSocketResponse{
		Type: "player_deleted",
	})
}
