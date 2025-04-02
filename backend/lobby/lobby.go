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
	lobbyMutex.Unlock()

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	if player.LobbyID != "" {
		log.Printf("Player %s already in lobby %s", playerID, player.LobbyID)
		player.PlayerStructMutex.Unlock()
		return nil, errors.New("player already in another lobby")
	}
	player.PlayerStructMutex.Unlock()

	lobbyMutex.Lock()
	lobbyID := utils.GenerateLobbyID(&lobbies)
	lobbyMutex.Unlock()

	player.PlayerStructMutex.Lock()
	player.LobbyID = lobbyID
	player.PlayerStructMutex.Unlock()

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

	lobby.LobbyStructMutex.Lock()
	if _, ok := lobby.PlayerIDs[playerID]; ok {
		log.Println("Player already in lobby")
		lobby.LobbyStructMutex.Unlock()
		return nil
	}
	if len(lobby.PlayerIDs) >= MAX_PLAYERS_PER_LOBBY {
		log.Printf("Lobby %s is full", lobbyID)
		lobby.LobbyStructMutex.Unlock()
		return errors.New("lobby is full")
	}
	lobby.PlayerIDs[playerID] = false
	lobby.LobbyStructMutex.Unlock()

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	player.LobbyID = lobbyID
	player.PlayerStructMutex.Unlock()

	log.Printf("Player %s joined lobby %s", playerID, lobbyID)
	return nil
}

func LeaveLobby(playerID, lobbyID string) error {
	lobby, _ := GetLobby(lobbyID)

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	if player.LobbyID != lobbyID {
		log.Printf("Player %s not in lobby %s", playerID, lobbyID)
		player.PlayerStructMutex.Unlock()
		return errors.New("player not in lobby")
	} else {
		player.LobbyID = ""
	}
	player.PlayerStructMutex.Unlock()

	lobby.LobbyStructMutex.Lock()
	delete(lobby.PlayerIDs, playerID)
	log.Printf("Player %s left lobby %s", playerID, lobbyID)

	if lobby.AdminPlayerID == playerID && len(lobby.PlayerIDs) > 0 {
		for newAdminID := range lobby.PlayerIDs {
			lobby.AdminPlayerID = newAdminID
			break
		}
		log.Printf("Player %s is now admin of lobby %s", lobby.AdminPlayerID, lobbyID)
	}

	if len(lobby.PlayerIDs) == 0 {
		delete(lobbies, lobbyID)
		log.Printf("No player in lobby %s deleted, hence deleted", lobbyID)
	}
	lobby.LobbyStructMutex.Unlock()
	return nil
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
