package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
	"wiki-woyage/utils"
)

func GetWWEnvVars(key string, def int) int {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return def
	} else {
		if s, err := strconv.Atoi(val); err == nil {
			return s
		} else {
			return def
		}
	}
}

var (
	MAX_PLAYERS_PER_LOBBY  int = GetWWEnvVars("WW_MAX_PLAYERS_PER_LOBBY", 8)
	MAX_LOBBIES            int = GetWWEnvVars("WW_MAX_LOBBIES", 100)
	MAX_PLAYER_NAME_LENGTH int = GetWWEnvVars("WW_MAX_PLAYER_NAME_LENGTH", 20)
	MIN_PLAYER_NAME_LENGTH int = GetWWEnvVars("WW_MIN_PLAYER_NAME_LENGTH", 3)
)

var (
	lobbies     = make(map[string][]string) // lobbyID -> player names
	playerNames = make(map[string]string)   // playerName -> playerID
	playerIDs   = make(map[string]string)   // playerID -> playerName
	mutex       sync.RWMutex
)

func CreateLobby(playerName string) (string, string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(lobbies) >= MAX_LOBBIES {
		log.Println("Max lobbies reached")
		return "", "", errors.New("max lobbies reached")
	}
	if len(playerName) < MIN_PLAYER_NAME_LENGTH || len(playerName) > MAX_PLAYER_NAME_LENGTH {
		log.Printf("Player name %s is invalid", playerName)
		return "", "", errors.New("player name is invalid")
	}
	if _, ok := playerNames[playerName]; ok {
		log.Printf("Player name %s already in a lobby", playerName)
		return "", "", errors.New("player name already in a lobby")
	}

	lobbyID := utils.GenerateLobbyID(&lobbies)
	lobbies[lobbyID] = []string{playerName}

	playerID := utils.GeneratePlayerID(&playerIDs)
	playerNames[playerName] = playerID
	playerIDs[playerID] = playerName

	log.Printf("Lobby %s created by %s", lobbyID, playerName)
	return lobbyID, playerID, nil
}

func JoinLobby(lobbyID string, playerName string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := playerNames[playerName]; ok {
		log.Println("Player already in a lobby")
		return "", errors.New("player already in a lobby")
	}

	players, ok := lobbies[lobbyID]
	if !ok {
		log.Printf("Lobby %s not found", lobbyID)
		return "", errors.New("lobby not found")
	}
	if len(players) >= MAX_PLAYERS_PER_LOBBY {
		log.Printf("Lobby %s is full", lobbyID)
		return "", errors.New("lobby is full")
	}

	lobbies[lobbyID] = append(players, playerName)
	playerID := utils.GeneratePlayerID(&playerIDs)
	playerNames[playerName] = playerID
	playerIDs[playerID] = playerName

	log.Printf("Player %s joined lobby %s", playerName, lobbyID)
	return playerID, nil
}

func LeaveLobby(lobbyID string, playerID string) error {
	mutex.Lock()
	defer mutex.Unlock()

	players, ok := lobbies[lobbyID]
	if !ok {
		log.Printf("Lobby %s not found", lobbyID)
		return errors.New("lobby not found")
	}
	playerName, ok := playerIDs[playerID]
	if !ok {
		log.Printf("PlayerID %s not found", playerID)
		return errors.New("player not found")
	}

	for i, player := range players {
		if player == playerName {
			lobbies[lobbyID] = append(players[:i], players[i+1:]...)
			delete(playerNames, playerName)
			delete(playerIDs, playerID)
			log.Printf("Player %s left lobby %s", playerName, lobbyID)

			if len(lobbies[lobbyID]) == 0 {
				delete(lobbies, lobbyID)
				log.Printf("Lobby %s deleted", lobbyID)
			}
			break
		}
	}
	return nil
}
