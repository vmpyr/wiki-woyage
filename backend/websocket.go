package main

import (
	"errors"
	"log"
	"net/http"
	"time"

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

	var msg structs.WebSocketMessage
	if err := conn.ReadJSON(&msg); err != nil || msg.PlayerID == "" {
		log.Println("WebSocket Read Error:", err)
		utils.SendError(conn, "Missing / invalid PlayerID")
		conn.Close()
		return
	}

	p, err := player.GetPlayer(msg.PlayerID)
	if err != nil {
		p, err = player.CreatePlayer(msg.PlayerID)
		if err != nil {
			log.Println("Player creation error:", err)
			utils.SendError(conn, err.Error())
			conn.Close()
			return
		}
	}

	p.PlayerStructMutex.Lock()
	p.Conn = conn
	p.LastActive = time.Now()
	p.PlayerStructMutex.Unlock()

	go HandleMessages(p)
}

func HandleMessages(p *structs.Player) {
	defer HandleDisconnect(p)

	for {
		var msg structs.WebSocketMessage
		if err := p.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket closed unexpectedly:", err)
				break
			}
			log.Println("WebSocket Read Error:", err)
			continue
		}

		if msg.LobbyID != "" && !CheckLobbyExists(msg.LobbyID) {
			utils.SendError(p.Conn, "Lobby not found")
			continue
		} else if msg.GameID != "" && !CheckGameExists(msg.GameID) {
			utils.SendError(p.Conn, "Game not found")
			continue
		}

		switch {
		case msg.Type == "create_player":
			player.HandleAddPlayerName(p, msg.PlayerName)
		case msg.Type == "create_lobby":
			lobby.HandleCreateLobby(p)

		case utils.HasSuffix(msg.Type, "_lobby"):
			if err := CheckEmptyIDRequest(msg.LobbyID, "LobbyID"); err != nil {
				log.Println("LobbyID is empty")
				utils.SendError(p.Conn, "LobbyID is empty")
				continue
			}
			switch msg.Type {
			case "join_lobby":
				lobby.HandleJoinLobby(p, msg.LobbyID)
			case "leave_lobby":
				lobby.HandleLeaveLobby(p, msg.LobbyID)
			case "toggle_ready_lobby":
				lobby.HandleToggleReady(p, msg.LobbyID)
			default:
				log.Println("Unknown lobby message type:", msg.Type)
			}

		case utils.HasSuffix(msg.Type, "_game"):
			if err := CheckEmptyIDRequest(msg.GameID, "GameID"); err != nil {
				log.Println("GameID is empty")
				utils.SendError(p.Conn, "GameID is empty")
				continue
			}
			switch msg.Type {
			case "start_game":
				game.HandleStartNewGame(p, msg.LobbyID, msg.GameSettings)
			case "join_game":
				game.HandleJoinGame(p, msg.GameID)
			case "vote_to_end_game":
				game.HandleVoteToEnd(p, msg.GameID)
			default:
				log.Println("Unknown game message type:", msg.Type)
			}

		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}
}

func HandleDisconnect(p *structs.Player) {
	log.Printf("Player %s disconnected, starting grace period...", p.PlayerID)
	time.Sleep(300 * time.Second)

	if time.Since(p.LastActive) > 300*time.Second {
		log.Printf("Player %s did not reconnect, removing...", p.PlayerID)
		player.DeletePlayer(p)
		lobby.LeaveLobby(p, p.LobbyID, true)
		game.LeaveGame(p, p.GameID, true)
	}
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

func CheckEmptyIDRequest(xID string, nameID string) error {
	if xID == "" {
		return errors.New(nameID + " is empty")
	}
	return nil
}
