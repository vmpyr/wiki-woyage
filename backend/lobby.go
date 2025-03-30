package backend

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

const MaxPlayersPerLobby = 8
const MaxLobbies = 100
const MaxPlayerNameLength = 20
const MinPlayerNameLength = 3

var (
	activeLobbies     = make(map[string][]string)
	activePlayerNames = make(map[string]string)
	activePlayerIDs   = make(map[string]string)
	mutex             sync.RWMutex
)

func generateID(characters string, length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func generateLobbyID() string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		lobbyID := generateID(characters, 5)
		if _, exists := activeLobbies[lobbyID]; !exists {
			return lobbyID
		}
	}
}

func generatePlayerID() string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for {
		playerID := generateID(characters, 25)
		if _, exists := activePlayerIDs[playerID]; !exists {
			return playerID
		}
	}
}

func CreateLobby(playerName string) (string, string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(activeLobbies) >= MaxLobbies {
		log.Println("Max lobbies reached")
		return "", "", errors.New("max lobbies reached")
	}
	if len(playerName) < MinPlayerNameLength || len(playerName) > MaxPlayerNameLength {
		log.Printf("Player name %s is invalid", playerName)
		return "", "", errors.New("player name is invalid")
	}
	if _, exists := activePlayerNames[playerName]; exists {
		log.Printf("Player name %s already in a lobby", playerName)
		return "", "", errors.New("player name already in a lobby")
	}

	lobbyID := generateLobbyID()
	activeLobbies[lobbyID] = []string{playerName}

	playerID := generatePlayerID()
	activePlayerNames[playerName] = playerID
	activePlayerIDs[playerID] = playerName

	log.Printf("Lobby %s created by %s", lobbyID, playerName)
	return lobbyID, playerID, nil
}

func JoinLobby(lobbyID string, playerName string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := activePlayerNames[playerName]; exists {
		log.Println("Player already in a lobby")
		return "", errors.New("player already in a lobby")
	}

	players, exists := activeLobbies[lobbyID]
	if !exists {
		log.Printf("Lobby %s not found", lobbyID)
		return "", errors.New("lobby not found")
	}
	if len(players) >= MaxPlayersPerLobby {
		log.Printf("Lobby %s is full", lobbyID)
		return "", errors.New("lobby is full")
	}

	activeLobbies[lobbyID] = append(players, playerName)
	playerID := generatePlayerID()
	activePlayerNames[playerName] = playerID
	activePlayerIDs[playerID] = playerName

	log.Printf("Player %s joined lobby %s", playerName, lobbyID)
	return playerID, nil
}

func LeaveLobby(lobbyID string, playerID string) error {
	mutex.Lock()
	defer mutex.Unlock()

	players, exists := activeLobbies[lobbyID]
	if !exists {
		log.Printf("Lobby %s not found", lobbyID)
		return errors.New("lobby not found")
	}
	playerName, exists := activePlayerIDs[playerID]
	if !exists {
		log.Printf("PlayerID %s not found", playerID)
		return errors.New("player not found")
	}

	for i, player := range players {
		if player == playerName {
			activeLobbies[lobbyID] = append(players[:i], players[i+1:]...)
			delete(activePlayerNames, playerName)
			delete(activePlayerIDs, playerID)
			log.Printf("Player %s left lobby %s", playerName, lobbyID)

			if len(activeLobbies[lobbyID]) == 0 {
				delete(activeLobbies, lobbyID)
				log.Printf("Lobby %s deleted", lobbyID)
			}
			break
		}
	}
	return nil
}
