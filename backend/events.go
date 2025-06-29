package main

import (
	"encoding/json"
	"log"
	"reflect"
)

type Event struct {
	Type    string          `json:"type"`
	Handler string          `json:"handler"` // lobby or orchestrator
	Payload json.RawMessage `json:"payload"`
}

const (
	EventAmIAdmin   = "am_i_admin"
	EventDisconnect = "disconnect"
)

var EventTypeToPayloadType = map[string]reflect.Type{
	EventAmIAdmin:   reflect.TypeOf(AmIAdminPayload{}),
	EventDisconnect: reflect.TypeOf(DisconnectPayload{}),
}

type EventHandler func(payload any, player *Player) error

type AmIAdminPayload struct{}

type DisconnectPayload struct{}

func HandleEventGeneric(
	event Event,
	player *Player,
	eventHandlers map[string]EventHandler,
) {
	payloadType, ok := EventTypeToPayloadType[event.Type]
	if !ok {
		log.Println("No payload type for event type", event.Type)
		return
	}
	payload := reflect.New(payloadType).Interface()
	if err := json.Unmarshal(event.Payload, payload); err != nil {
		log.Println("Failed to unmarshal payload for", player.username, ":", err)
		return
	}

	if handler, ok := eventHandlers[event.Type]; ok {
		if err := handler(payload, player); err != nil {
			log.Println("Error handling event for", player.username, ":", err)
		}
	} else {
		log.Println("No handler for event type", event.Type, "from", player.username)
	}
}
