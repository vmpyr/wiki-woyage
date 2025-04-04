package structs

import (
	"sync"
)

type Lobby struct {
	LobbyID          string          // lobbyID
	GameID           string          // gameID
	PlayerIDs        map[string]bool // playerIDs -> true if ready / ingame
	AdminPlayerID    string          // playerID of the admin
	LobbyStructMutex sync.Mutex
}
