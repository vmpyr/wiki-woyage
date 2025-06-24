package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Orchestrator struct {
	lobbies      LobbyList
	eventHandler map[string]EventHandler
	mutex        sync.RWMutex
}

func CreateOrchestrator() *Orchestrator {
	o := &Orchestrator{
		lobbies:      make(LobbyList),
		eventHandler: make(map[string]EventHandler),
	}
	return o
}

func (o *Orchestrator) ServeWS(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	lobbyID := r.URL.Query().Get("lobbyID")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	log.Println("New connection request from", username, "for lobby", lobbyID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var (
		lobby  *Lobby
		player *Player
	)

	if lobbyID == "" {
		lobby, err = o.NewLobby()
		if err != nil {
			log.Println("Failed to create new lobby:", err)
			http.Error(w, "Failed to create new lobby", http.StatusInternalServerError)
			return
		}
		log.Println("Created new lobby with ID:", lobby.id)
	} else {
		o.mutex.RLock()
		lobby, ok := o.lobbies[lobbyID]
		o.mutex.RUnlock()
		if !ok {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}
		log.Println("Joining existing lobby with ID:", lobby.id)
	}

	player, err = lobby.CreatePlayer(conn, username)
	if err != nil {
		log.Println("Failed to create player:", err)
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}
	err = lobby.AddPlayerToLobby(player)
	if err != nil {
		log.Println("Failed to add player to lobby:", err)
		http.Error(w, "Failed to add player to lobby", http.StatusInternalServerError)
		return
	}

	log.Println("Player", player.username, "joined lobby", lobby.id)

	go player.Run()
}
