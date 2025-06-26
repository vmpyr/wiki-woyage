package main

import (
	"sync"
	"time"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	id            string
	players       PlayerList
	admin         string
	eventHandlers map[string]EventHandler
	mutex         sync.RWMutex

	createdAt  time.Time
	lastActive time.Time
}

func (o *Orchestrator) NewLobby() (*Lobby, error) {
	l := &Lobby{
		id:            GenerateLobbyID(&o.lobbies),
		players:       make(PlayerList),
		admin:         "",
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
	l.eventHandlers[EventAmIAdmin] = l.HandleAmIAdmin
}

func (l *Lobby) HandleEvent(event Event, player *Player) {
	HandleEventGeneric(event, player, l.eventHandlers)
}

func (l *Lobby) HandleAmIAdmin(payload any, player *Player) error {
	isAdmin := false
	if l.admin == player.clientID {
		isAdmin = true
	}
	err := player.SendWSResponse(ResponseAmIAdmin, AmIAdminResponse{
		IsAdmin: isAdmin,
	})
	return err
}
