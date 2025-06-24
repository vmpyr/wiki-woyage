package main

import (
	"log"
	"net/http"
	"sync"
	"time"

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
		log.Println("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	log.Println("New connection request from", username, "for lobby", lobbyID)

	var (
		lobby  *Lobby
		player *Player
		err    error
		ok     bool
	)

	if lobbyID == "" {
		lobby, err = o.NewLobby()
		if err != nil {
			http.Error(w, "Failed to create new lobby", http.StatusInternalServerError)
			return
		}
	} else {
		o.mutex.RLock()
		lobby, ok = o.lobbies[lobbyID]
		o.mutex.RUnlock()
		if !ok {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
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

func (o *Orchestrator) StartLobbyCleanup(interval, idleTimeout time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			o.CleanupLobbies(idleTimeout)
		}
	}()
}

func (o *Orchestrator) CleanupLobbies(idleTimeout time.Duration) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	now := time.Now()
	for lobbyID, lobby := range o.lobbies {
		lobby.mutex.RLock()
		isEmpty := len(lobby.players) == 0
		inactive := now.Sub(lobby.lastActive) > idleTimeout
		lobby.mutex.RUnlock()

		if isEmpty && inactive {
			log.Println("Deleting idle empty lobby:", lobbyID)
			delete(o.lobbies, lobbyID)
		}
	}
}
