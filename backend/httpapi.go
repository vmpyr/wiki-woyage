package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendHTTPResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://localhost:%d", SVELTE_CLIENT_PORT))
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, writeErr := w.Write(jsonBytes)
	if writeErr != nil {
		log.Printf("Failed to write response: %v", writeErr)
	}
}

// TODO: aaaaaah, I should've simply used the websocket connection for this. anywho, change later
func (o *Orchestrator) HandleClientInfo(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientID")
	if clientID == "" {
		log.Println("ClientID is required")
		http.Error(w, "clientID is required", http.StatusBadRequest)
		return
	}

	o.mutex.RLock()
	player, ok := o.players[clientID]
	o.mutex.RUnlock()
	if !ok {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	info := map[string]string{
		"username": player.username,
		"lobbyID":  "",
	}
	if player.lobby != nil {
		info["lobbyID"] = player.lobby.id
	}

	SendHTTPResponse(w, http.StatusOK, info)
}
