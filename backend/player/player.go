package player

import (
	"errors"
	"log"
	"sync"
	"time"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
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

func CreatePlayer(playerID string) (*st.Player, error) {
	p := &st.Player{
		PlayerID:          playerID,
		PlayerName:        playerID,
		LobbyID:           "",
		GameID:            "",
		LastActive:        time.Now(),
		Conn:              nil,
		PlayerStructMutex: sync.Mutex{},
	}

	playerMutex.Lock()
	defer playerMutex.Unlock()
	players[playerID] = p
	takenPlayerNames[playerID] = playerID
	log.Printf("Player created with ID %s", playerID)
	return p, nil
}

func AddPlayerName(p *st.Player, playerName string) error {
	if len(playerName) < MIN_PLAYER_NAME_LENGTH || len(playerName) > MAX_PLAYER_NAME_LENGTH {
		log.Printf("Player name %s is invalid", playerName)
		return errors.New("player name is invalid")
	}

	playerMutex.Lock()
	if _, ok := takenPlayerNames[playerName]; ok {
		log.Printf("Player name %s already taken", playerName)
		return errors.New("player name already taken")
	}
	delete(takenPlayerNames, p.PlayerName)
	takenPlayerNames[playerName] = p.PlayerID
	playerMutex.Unlock()

	p.PlayerStructMutex.Lock()
	p.PlayerName = playerName
	p.PlayerStructMutex.Unlock()
	log.Printf("Player %s changed name to %s", p.PlayerID, playerName)
	return nil
}

func DeletePlayer(p *st.Player) error {
	playerMutex.Lock()
	defer playerMutex.Unlock()

	delete(takenPlayerNames, p.PlayerName)
	delete(players, p.PlayerID)
	log.Printf("Player %s deleted", p.PlayerID)
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
