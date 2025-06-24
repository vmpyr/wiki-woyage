package main

import "github.com/gorilla/websocket"

type PlayerList map[*Player]bool
type Player struct {
	connection *websocket.Conn
	username   string
	lobby      *Lobby
}

func (l *Lobby) CreatePlayer(conn *websocket.Conn, username string) (*Player, error) {
	if ok := CheckUniqueUsername(username, &l.players); !ok {
		conn.Close()
		return nil, ErrUsernameTaken
	}

	player := &Player{
		connection: conn,
		username:   username,
		lobby:      l,
	}
	l.AddPlayerToLobby(player)
	return player, nil
}

func (l *Lobby) AddPlayerToLobby(player *Player) error {
	if player.lobby != nil {
		player.lobby.RemovePlayerFromLobby(player)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.players[player] = true
	player.lobby = l
	return nil
}

func (l *Lobby) RemovePlayerFromLobby(player *Player) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.players, player)
	player.lobby = nil
}
