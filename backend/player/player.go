package player

import (
	"errors"
	"log"
	"sync"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

var (
	MAX_PLAYER_NAME_LENGTH int = utils.GetWWEnvVars("WW_MAX_PLAYER_NAME_LENGTH", 20)
	MIN_PLAYER_NAME_LENGTH int = utils.GetWWEnvVars("WW_MIN_PLAYER_NAME_LENGTH", 3)
)

var (
	players          = make(map[string]*st.Player) // playerID -> Player
	takenPlayerNames = make(map[string]string)     // playerName -> playerID
	playerMutex      sync.RWMutex
)

func CreatePlayer(conn *websocket.Conn, playerName string) (*st.Player, error) {
	if len(playerName) < MIN_PLAYER_NAME_LENGTH || len(playerName) > MAX_PLAYER_NAME_LENGTH {
		log.Printf("Player name %s is invalid", playerName)
		return nil, errors.New("player name is invalid")
	}

	playerMutex.Lock()
	defer playerMutex.Unlock()

	if _, ok := takenPlayerNames[playerName]; ok {
		log.Printf("Player name %s already taken", playerName)
		return nil, errors.New("player name already taken")
	}

	playerID := utils.GeneratePlayerID(&players)
	player := &st.Player{
		PlayerID:          playerID,
		PlayerName:        playerName,
		LobbyID:           "",
		GameID:            "",
		Conn:              conn,
		PlayerStructMutex: sync.Mutex{},
	}
	players[playerID] = player
	takenPlayerNames[playerName] = playerID
	log.Printf("Player %s created with ID %s", playerName, playerID)
	return player, nil
}

func DeletePlayer(playerID string) error {
	playerMutex.Lock()
	defer playerMutex.Unlock()

	player, ok := players[playerID]
	if !ok {
		log.Printf("Player %s not found", playerID)
		return errors.New("player not found")
	}

	delete(takenPlayerNames, player.PlayerName)
	delete(players, playerID)
	log.Printf("Player %s deleted", playerID)
	return nil
}

func GetPlayer(playerID string) (*st.Player, error) {
	playerMutex.RLock()
	defer playerMutex.RUnlock()

	player, ok := players[playerID]
	if !ok {
		log.Printf("Player %s not found", playerID)
		return nil, errors.New("player not found")
	}

	return player, nil
}

func GetPlayerName(playerID string) (string, error) {
	player, err := GetPlayer(playerID)
	if err != nil {
		return "", err
	}

	return player.PlayerName, nil
}

func GetPlayerLobbyID(playerID string) (string, error) {
	player, err := GetPlayer(playerID)
	if err != nil {
		return "", err
	}

	return player.LobbyID, nil
}
