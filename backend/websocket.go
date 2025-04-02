package main

import (
	"log"
	"net/http"

	"wiki-woyage/lobby"
	"wiki-woyage/player"
	"wiki-woyage/structs"

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

		switch msg.Type {
		case "create_player":
			player.HandleCreatePlayer(conn, msg.PlayerName)
		case "create_lobby":
			lobby.HandleCreateLobby(conn, msg.PlayerID)
		case "join_lobby":
			lobby.HandleJoinLobby(conn, msg.PlayerID, msg.LobbyID)
		case "leave_lobby":
			lobby.HandleLeaveLobby(conn, msg.PlayerID, msg.LobbyID)
		case "toggle_ready":
			lobby.HandleToggleReady(conn, msg.PlayerID, msg.LobbyID)
		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}
}
