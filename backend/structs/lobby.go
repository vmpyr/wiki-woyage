package structs

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Lobby struct {
	LobbyID          string          // lobbyID
	GameID           string          // gameID
	PlayerIDs        map[string]bool // playerIDs -> true if ready / ingame
	AdminPlayerID    string          // playerID of the admin
	Conn             *websocket.Conn
	LobbyStructMutex sync.Mutex
}
