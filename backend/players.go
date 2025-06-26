package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerList map[string]*Player // clientID -> Player
type Player struct {
	connection *websocket.Conn
	username   string
	clientID   string
	lobby      *Lobby

	send  chan []byte
	done  chan bool
	close sync.Once
}

func (l *Lobby) CreatePlayer(conn *websocket.Conn, username, clientID string, orchestrator *Orchestrator) (*Player, error) {
	if ok := CheckUniqueUsername(username, &l.players); !ok {
		conn.Close()
		return nil, ErrUsernameTaken
	}

	player := &Player{
		connection: conn,
		username:   username,
		clientID:   clientID,
		lobby:      l,
		send:       make(chan []byte, 256),
		done:       make(chan bool),
		close:      sync.Once{},
	}
	l.AddPlayerToLobby(player)
	orchestrator.mutex.Lock()
	orchestrator.players[player.clientID] = player
	orchestrator.mutex.Unlock()
	return player, nil
}

func (l *Lobby) AddPlayerToLobby(player *Player) error {
	if player.lobby != nil {
		player.lobby.RemovePlayerFromLobby(player)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.players[player.clientID] = player
	l.lastActive = time.Now()
	player.lobby = l
	if l.admin == "" {
		l.admin = player.clientID
	}
	return nil
}

func (l *Lobby) RemovePlayerFromLobby(player *Player) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.players, player.clientID)
	l.lastActive = time.Now()
	if l.admin == player.clientID {
		if len(l.players) > 0 {
			for newAdmin := range l.players {
				l.admin = newAdmin
				break
			}
		} else {
			l.admin = ""
		}
	}
	player.lobby = nil
}
