package main

import (
	"sync"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	id      string
	players PlayerList
	mutex   sync.RWMutex
}

func (o *Orchestrator) NewLobby() (*Lobby, error) {
	l := &Lobby{
		id:      GenerateLobbyID(&o.lobbies),
		players: make(PlayerList),
		mutex:   sync.RWMutex{},
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.lobbies[l.id] = l
	return l, nil
}
