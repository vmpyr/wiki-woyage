package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type       string `json:"type"`
	PlayerName string `json:"playerName,omitempty"`
	LobbyID    string `json:"lobbyID,omitempty"`
	PlayerID   string `json:"playerID,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var connections = make(map[string]map[*websocket.Conn]bool)
var connMutex sync.RWMutex

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			break
		}

		switch msg.Type {
		case "create_lobby":
			HandleCreateLobby(conn, msg.PlayerName)
		case "join_lobby":
			HandleJoinLobby(conn, msg.LobbyID, msg.PlayerName)
		case "leave_lobby":
			HandleLeaveLobby(conn, msg.LobbyID, msg.PlayerID)
		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}
}

func HandleCreateLobby(conn *websocket.Conn, playerName string) {
	lobbyID, playerID, err := CreateLobby(playerName)
	if err != nil {
		SendError(conn, err.Error())
		return
	}

	connMutex.Lock()
	if connections[lobbyID] == nil {
		connections[lobbyID] = make(map[*websocket.Conn]bool)
	}
	connections[lobbyID][conn] = true
	connMutex.Unlock()

	SendResponse(conn, "lobby_created", lobbyID, playerID, playerName)
}

func HandleJoinLobby(conn *websocket.Conn, lobbyID string, playerName string) {
	playerID, err := JoinLobby(lobbyID, playerName)
	if err != nil {
		SendError(conn, err.Error())
		return
	}

	connMutex.Lock()
	if connections[lobbyID] == nil {
		connections[lobbyID] = make(map[*websocket.Conn]bool)
	}
	connections[lobbyID][conn] = true
	connMutex.Unlock()

	SendResponse(conn, "lobby_joined", lobbyID, playerID, playerName)
}

func HandleLeaveLobby(conn *websocket.Conn, lobbyID string, playerID string) {
	err := LeaveLobby(lobbyID, playerID)
	if err != nil {
		SendError(conn, err.Error())
		return
	}

	connMutex.Lock()
	delete(connections[lobbyID], conn)
	if len(connections[lobbyID]) == 0 {
		delete(connections, lobbyID)
	}
	connMutex.Unlock()

	SendResponse(conn, "lobby_left", lobbyID, playerID, "")
}

func SendResponse(conn *websocket.Conn, msgType, lobbyID, playerID, playerName string) {
	response := WebSocketMessage{
		Type:       msgType,
		LobbyID:    lobbyID,
		PlayerID:   playerID,
		PlayerName: playerName,
	}
	log.Println("Sending response:", response)
	conn.WriteJSON(response)
}

func SendError(conn *websocket.Conn, errorMsg string) {
	response := map[string]string{"type": "error", "message": errorMsg}
	log.Println("Sending error response:", errorMsg)
	conn.WriteJSON(response)
}
