package main

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Handler string          `json:"handler"` // lobby or orchestrator
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, player *Player) error
