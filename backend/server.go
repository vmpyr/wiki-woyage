package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// enforce later
const (
	MAX_PLAYERS_PER_LOBBY = 8
	MAX_LOBBIES           = 100
	MAX_USERNAME_LENGTH   = 20
	MIN_USERNAME_LENGTH   = 3
	LOBBY_TICKER_INTERVAL = 120
	LOBBY_IDLE_TIMEOUT    = 360
	SVELTE_CLIENT_PORT    = 5174
)

var port int
var listeningInterface string

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&listeningInterface, "interface", "127.0.0.1", "Interface to listen on")
	flag.Parse()
}

func main() {
	orchestrator := CreateOrchestrator()
	orchestrator.StartLobbyCleanup(LOBBY_TICKER_INTERVAL*time.Second, LOBBY_IDLE_TIMEOUT*time.Second)

	http.HandleFunc("/ws", orchestrator.ServeWS)
	http.HandleFunc("/api/clientinfo", orchestrator.HandleClientInfo)

	listenOn := fmt.Sprintf("%s:%d", listeningInterface, port)
	fmt.Printf("Server started on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
