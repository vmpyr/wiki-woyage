package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// do later
// const (
// 	MAX_PLAYERS_PER_LOBBY = 8
// 	MAX_LOBBIES           = 100
// 	MAX_USERNAME_LENGTH   = 20
// 	MIN_USERNAME_LENGTH   = 3
// )

var port int
var listeningInterface string

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&listeningInterface, "interface", "127.0.0.1", "Interface to listen on")
	flag.Parse()
}

func main() {
	orchestrator := CreateOrchestrator()

	http.HandleFunc("/ws", orchestrator.ServeWS)

	listenOn := fmt.Sprintf("%s:%d", listeningInterface, port)
	fmt.Printf("Server started on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
