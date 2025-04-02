package main

import (
	"log"
	"net/http"

	"wiki-woyage/game"
	"wiki-woyage/lobby"
	"wiki-woyage/player"
	"wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	for {
		var msg structs.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			break
		}

		if msg.PlayerID != "" && !CheckPlayerExists(msg.PlayerID) {
			utils.SendError(conn, "Player not found")
			continue
		} else if msg.LobbyID != "" && !CheckLobbyExists(msg.LobbyID) {
			utils.SendError(conn, "Lobby not found")
			continue
		} else if msg.GameID != "" && !CheckGameExists(msg.GameID) {
			utils.SendError(conn, "Game not found")
			continue
		}

		switch msg.Type {
		case "create_player":
			player.HandleCreatePlayer(conn, msg.PlayerName)
		case "create_lobby":
			lobby.HandleCreateLobby(conn, msg.PlayerID)
		case "join_lobby":
			lobby.HandleJoinLobby(conn, msg.PlayerID, msg.LobbyID)
		case "leave_lobby":
			lobby.HandleLeaveLobby(conn, msg.PlayerID, msg.LobbyID)
		case "start_game":
			game.HandleStartNewGame(conn, msg.PlayerID, msg.LobbyID, msg.GameSettings)
		case "join_game":
			game.HandleJoinGame(conn, msg.PlayerID, msg.GameID)
		case "vote_to_end_game":
			game.HandleVoteToEnd(conn, msg.PlayerID, msg.GameID)
		case "toggle_ready":
			lobby.HandleToggleReady(conn, msg.PlayerID, msg.LobbyID)
		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}
}

func CheckPlayerExists(playerID string) bool {
	if _, err := player.GetPlayer(playerID); err != nil {
		log.Printf("Player %s not found", playerID)
		return false
	}
	return true
}

func CheckLobbyExists(lobbyID string) bool {
	if _, err := lobby.GetLobby(lobbyID); err != nil {
		log.Printf("Lobby %s not found", lobbyID)
		return false
	}
	return true
}

func CheckGameExists(gameID string) bool {
	if _, err := game.GetGame(gameID); err != nil {
		log.Printf("Game %s not found", gameID)
		return false
	}
	return true
}
