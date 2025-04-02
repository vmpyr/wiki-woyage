package lobby

import (
	"errors"
	"log"
	"sync"
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

func CreateLobby(conn *websocket.Conn, playerID string) (*st.Lobby, error) {
	lobbyMutex.Lock()
	defer lobbyMutex.Unlock()

	if len(lobbies) >= MAX_LOBBIES {
		log.Println("Max lobbies reached")
		return nil, errors.New("max lobbies reached")
	}
	if _, ok := lobbies[playerID]; ok {
		log.Println("Player already in lobby")
		return nil, errors.New("player already in a lobby")
	}

	lobbyID := utils.GenerateLobbyID(&lobbies)
	lobby := &st.Lobby{
		LobbyID:          lobbyID,
		GameID:           "",
		PlayerIDs:        make(map[string]bool),
		Conn:             conn,
		LobbyStructMutex: sync.Mutex{},
	}
	lobby.PlayerIDs[playerID] = false
	lobbies[lobbyID] = lobby

	log.Printf("Lobby %s created by %s", lobbyID, playerID)
	return lobby, nil
}

func JoinLobby(playerID, lobbyID string) error {
	lobby, err := GetLobby(lobbyID)
	if err != nil {
		return err
	}

	lobby.LobbyStructMutex.Lock()
	defer lobby.LobbyStructMutex.Unlock()

	if _, ok := lobby.PlayerIDs[playerID]; ok {
		log.Println("Player already in lobby")
		return nil
	}
	if len(lobby.PlayerIDs) >= MAX_PLAYERS_PER_LOBBY {
		log.Printf("Lobby %s is full", lobbyID)
		return errors.New("lobby is full")
	}

	lobby.PlayerIDs[playerID] = false
	log.Printf("Player %s joined lobby %s", playerID, lobbyID)
	return nil
}

func LeaveLobby(playerID, lobbyID string) error {
	lobby, err := GetLobby(lobbyID)
	if err != nil {
		return err
	}

	lobby.LobbyStructMutex.Lock()
	defer lobby.LobbyStructMutex.Unlock()

	delete(lobby.PlayerIDs, playerID)
	log.Printf("Player %s left lobby %s", playerID, lobbyID)

	if len(lobby.PlayerIDs) == 0 {
		delete(lobbies, lobbyID)
		log.Printf("Lobby %s deleted", lobbyID)
	}
	return nil
}

func ToggleReady(playerID, lobbyID string) (bool, error) {
	lobby, err := GetLobby(lobbyID)
	if err != nil {
		return false, err
	}

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
