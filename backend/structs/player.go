package structs

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	PlayerID          string
	PlayerName        string
	LobbyID           string
	GameID            string
	LastActive        time.Time
	Conn              *websocket.Conn
	PlayerStructMutex sync.Mutex
}
