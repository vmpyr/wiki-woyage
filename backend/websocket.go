package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (p *Player) Run(orchestrator *Orchestrator) {
	go p.ReadMessages(orchestrator)
	go p.WriteMessages()
	<-p.done
	p.Cleanup(orchestrator)
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

func (p *Player) SendWSResponse(responseType string, response any) error {
	msg := Response{
		Type:     responseType,
		Response: response,
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

func (l *Lobby) BroadcastResponse(responseType string, response any) {
	for _, player := range l.players {
		player.SendWSResponse(responseType, response)
	}
}

func (p *Player) SignalDisconnect() {
	p.close.Do(func() {
		close(p.done)
	})
}

func (p *Player) Cleanup(orchestrator *Orchestrator) {
	log.Println("Cleaning up player", p.username)
	p.connection.Close()
	close(p.send)
	p.lobby.RemovePlayerFromLobby(p)
	orchestrator.mutex.Lock()
	delete(orchestrator.players, p.clientID)
	orchestrator.mutex.Unlock()
}
