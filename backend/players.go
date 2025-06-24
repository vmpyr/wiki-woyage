package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerList map[*Player]bool
type Player struct {
	connection *websocket.Conn
	username   string
	lobby      *Lobby

	send  chan []byte
	done  chan bool
	close sync.Once
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
		send:       make(chan []byte, 256),
		done:       make(chan bool),
		close:      sync.Once{},
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
	l.lastActive = time.Now()
	player.lobby = l
	return nil
}

func (l *Lobby) RemovePlayerFromLobby(player *Player) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.players, player)
	l.lastActive = time.Now()
	player.lobby = nil
}

func (p *Player) Run() {
	go p.ReadMessages()
	go p.WriteMessages()
	<-p.done
	p.Cleanup()
}

func (p *Player) ReadMessages() {
	defer p.SignalDisconnect()
	for {
		_, message, err := p.connection.ReadMessage()
		if err != nil {
			log.Println("Error reading message from", p.username, ":", err)
			return
		}
		// TODO: Handle incoming messages
		log.Println("Received message from", p.username, ":", string(message))
	}
}

func (p *Player) WriteMessages() {
	defer p.SignalDisconnect()
	for {
		select {
		case message, ok := <-p.send:
			if !ok {
				log.Println("Send channel closed for", p.username)
				return
			}
			err := p.connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error writing message to", p.username, ":", err)
				return
			}
		case <-p.done:
			log.Println("Done channel closed for", p.username)
			return
		}
	}
}

func (p *Player) SignalDisconnect() {
	p.close.Do(func() {
		close(p.done)
	})
}

func (p *Player) Cleanup() {
	log.Println("Cleaning up player", p.username)
	p.connection.Close()
	close(p.send)
	p.lobby.RemovePlayerFromLobby(p)
}
