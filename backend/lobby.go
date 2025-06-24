package main

import (
	"sync"
	"time"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	id      string
	players PlayerList
	mutex   sync.RWMutex

	createdAt  time.Time
	lastActive time.Time
}

func (o *Orchestrator) NewLobby() (*Lobby, error) {
	l := &Lobby{
		id:         GenerateLobbyID(&o.lobbies),
		players:    make(PlayerList),
		mutex:      sync.RWMutex{},
		createdAt:  time.Now(),
		lastActive: time.Now(),
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.lobbies[l.id] = l
	return l, nil
}

func (o *Orchestrator) DeleteLobby(lobbyID string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if lobby, ok := o.lobbies[lobbyID]; ok {
		for player := range lobby.players {
			player.SignalDisconnect()
			player.connection.Close()
			close(player.send)
		}
		delete(o.lobbies, lobbyID)
	}
}
