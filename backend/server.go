package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int
var listeningInterface string

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&listeningInterface, "interface", "127.0.0.1", "Interface to listen on")
	flag.Parse()
}

func main() {
	http.HandleFunc("/ws", HandleWebSocket)

	listenOn := fmt.Sprintf("%s:%d", listeningInterface, port)
	fmt.Printf("Server started on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
