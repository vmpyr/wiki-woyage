package structs

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	PlayerID          string
	PlayerName        string
	LobbyID           string
	GameID            string
	Conn              *websocket.Conn
	PlayerStructMutex sync.Mutex
}
