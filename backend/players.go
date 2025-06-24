package main

import (
	"encoding/json"
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

func (p *Player) SendEvent(eventType string, payload any) error {
	msg := map[string]any{
		"type":    eventType,
		"payload": payload,
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to marshal send event message for", p.username, ":", err)
		return err
	}
	select {
	case p.send <- jsonMsg:
		log.Println("Sent message to", p.username, ":", string(jsonMsg))
	case <-p.done:
		log.Println("Done channel closed for", p.username, ", cannot send message")
	}
	return nil
}

func (p *Player) Run(orchestrator *Orchestrator) {
	go p.ReadMessages(orchestrator)
	go p.WriteMessages()
	<-p.done
	p.Cleanup()
}

func (p *Player) ReadMessages(orchestrator *Orchestrator) {
	defer p.SignalDisconnect()
	for {
		_, message, err := p.connection.ReadMessage()
		if err != nil {
			log.Println("Error reading message from", p.username, ":", err)
			return
		}

		log.Println("Received message from", p.username, ":", string(message))

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println("Bad message from", p.username, ":", err)
			continue
		}

		switch event.Handler {
		case "lobby", "":
			p.lobby.HandleEvent(event, p)
			p.lobby.lastActive = time.Now()
		case "orchestrator":
			orchestrator.HandleEvent(event, p)
		default:
			log.Println("Unknown handler for event from", p.username, ":", event.Handler)
		}
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
