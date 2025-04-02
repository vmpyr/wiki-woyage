package lobby

import (
	"errors"
	"log"
	"sync"
	"wiki-woyage/player"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

var (
	MAX_PLAYERS_PER_LOBBY int = utils.GetWWEnvVars("WW_MAX_PLAYERS_PER_LOBBY", 8)
	MAX_LOBBIES           int = utils.GetWWEnvVars("WW_MAX_LOBBIES", 100)
)

var (
	lobbies    = make(map[string]*st.Lobby) // lobbyID -> Lobby
	lobbyMutex sync.RWMutex
)

func CreateLobby(conn *websocket.Conn, playerID string) (*st.Lobby, error) {
	lobbyMutex.Lock()
	if len(lobbies) >= MAX_LOBBIES {
		log.Println("Max lobbies reached")
		return nil, errors.New("max lobbies reached")
	}
	lobbyID := utils.GenerateLobbyID(&lobbies)
	lobbyMutex.Unlock()

	player, _ := player.GetPlayer(playerID)
	err := utils.MutexExecPlayer(player, func(p *st.Player) error {
		if p.LobbyID != "" {
			log.Printf("Player %s already in lobby %s", playerID, p.LobbyID)
			return errors.New("player already in another lobby")
		}
		p.LobbyID = lobbyID
		return nil
	})
	if err != nil {
		return nil, err
	}

	lobby := &st.Lobby{
		LobbyID:          lobbyID,
		GameID:           "",
		PlayerIDs:        map[string]bool{playerID: false},
		AdminPlayerID:    playerID,
		Conn:             conn,
		LobbyStructMutex: sync.Mutex{},
	}

	lobbyMutex.Lock()
	lobbies[lobbyID] = lobby
	lobbyMutex.Unlock()

	log.Printf("Lobby %s created by %s", lobbyID, playerID)
	return lobby, nil
}

func JoinLobby(playerID, lobbyID string) error {
	lobby, _ := GetLobby(lobbyID)
	err := utils.MutexExecLobby(lobby, func(l *st.Lobby) error {
		if _, ok := l.PlayerIDs[playerID]; ok {
			log.Printf("Player %s already in lobby %s", playerID, lobbyID)
			return errors.New("player already in lobby")
		}
		if len(l.PlayerIDs) >= MAX_PLAYERS_PER_LOBBY {
			log.Printf("Lobby %s is full", lobbyID)
			return errors.New("lobby is full")
		}
		l.PlayerIDs[playerID] = false
		log.Printf("Player %s joined lobby %s", playerID, lobbyID)
		return nil
	})
	if err != nil {
		if err.Error() == "player already in lobby" {
			return nil
		}
		return err
	}

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	player.LobbyID = lobbyID
	player.PlayerStructMutex.Unlock()

	log.Printf("Player %s joined lobby %s", playerID, lobbyID)
	return nil
}

func LeaveLobby(playerID, lobbyID string) error {
	player, _ := player.GetPlayer(playerID)
	err := utils.MutexExecPlayer(player, func(p *st.Player) error {
		if p.LobbyID != lobbyID {
			log.Printf("Player %s not in lobby %s", playerID, lobbyID)
			return errors.New("player not in lobby")
		}
		p.LobbyID = ""
		return nil
	})
	if err != nil {
		return err
	}

	lobby, _ := GetLobby(lobbyID)
	return utils.MutexExecLobby(lobby, func(l *st.Lobby) error {
		delete(l.PlayerIDs, playerID)
		log.Printf("Player %s left lobby %s", playerID, lobbyID)

		if l.AdminPlayerID == playerID && len(l.PlayerIDs) > 0 {
			for newAdminID := range l.PlayerIDs {
				l.AdminPlayerID = newAdminID
				break
			}
			log.Printf("Player %s is now admin of lobby %s", l.AdminPlayerID, lobbyID)
		}

		if len(l.PlayerIDs) == 0 {
			lobbyMutex.Lock()
			delete(lobbies, lobbyID)
			lobbyMutex.Unlock()
			log.Printf("No player in lobby %s deleted, hence deleted", lobbyID)
		}
		return nil
	})
}

func ToggleReady(playerID, lobbyID string) (bool, error) {
	lobby, _ := GetLobby(lobbyID)
	lobby.LobbyStructMutex.Lock()
	defer lobby.LobbyStructMutex.Unlock()

	if _, ok := lobby.PlayerIDs[playerID]; !ok {
		log.Println("Player not in lobby")
		return false, errors.New("player not in lobby")
	}

	lobby.PlayerIDs[playerID] = !lobby.PlayerIDs[playerID]
	toggledTo := lobby.PlayerIDs[playerID]
	log.Printf("Player %s toggled to %t in lobby %s", playerID, toggledTo, lobbyID)

	return toggledTo, nil
}

func GetLobby(lobbyID string) (*st.Lobby, error) {
	lobbyMutex.RLock()
	defer lobbyMutex.RUnlock()

	lobby, ok := lobbies[lobbyID]
	if !ok {
		log.Printf("Lobby %s not found", lobbyID)
		return nil, errors.New("lobby not found")
	}

	return lobby, nil
}

func GetLobbyGameID(lobbyID string) (string, error) {
	lobby, err := GetLobby(lobbyID)
	if err != nil {
		return "", err
	}

	return lobby.GameID, nil
}

func GetAdminPlayerID(lobbyID string) (string, error) {
	lobby, err := GetLobby(lobbyID)
	if err != nil {
		return "", err
	}

	return lobby.AdminPlayerID, nil
}
