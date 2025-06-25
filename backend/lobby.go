package main

import (
	"log"
	"sync"
	"time"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	id            string
	players       PlayerList
	eventHandlers map[string]EventHandler
	mutex         sync.RWMutex

	createdAt  time.Time
	lastActive time.Time
}

func (o *Orchestrator) NewLobby() (*Lobby, error) {
	l := &Lobby{
		id:            GenerateLobbyID(&o.lobbies),
		players:       make(PlayerList),
		eventHandlers: make(map[string]EventHandler),
		mutex:         sync.RWMutex{},
		createdAt:     time.Now(),
		lastActive:    time.Now(),
	}
	l.SetupEventHandlers()
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.lobbies[l.id] = l
	return l, nil
}

func (o *Orchestrator) DeleteLobby(lobbyID string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if lobby, ok := o.lobbies[lobbyID]; ok {
		for _, player := range lobby.players {
			player.SignalDisconnect()
			player.connection.Close()
			close(player.send)
		}
		delete(o.lobbies, lobbyID)
	}
}

func (l *Lobby) SetupEventHandlers() {
	// events to be added
}

func (l *Lobby) HandleEvent(event Event, player *Player) {
	if handler, ok := l.eventHandlers[event.Type]; ok {
		if err := handler(event, player); err != nil {
			log.Println("Error handling event for", player.username, ":", err)
		}
	} else {
		log.Println("No handler for event type", event.Type, "from", player.username)
	}
}
