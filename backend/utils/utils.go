package utils

import (
	"math/rand"
	"wiki-woyage/structs"
)

func GenerateID(characters string, length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func GenerateLobbyID(lobbies *map[string][]string) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		lobbyID := GenerateID(characters, 5)
		if _, ok := (*lobbies)[lobbyID]; !ok {
			return lobbyID
		}
	}
}

func GeneratePlayerID(playerIDs *map[string]string) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for {
		playerID := GenerateID(characters, 25)
		if _, ok := (*playerIDs)[playerID]; !ok {
			return playerID
		}
	}
}

func GenerateGameID(games *map[string]*structs.Game) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for {
		gameID := GenerateID(characters, 25)
		if _, ok := (*games)[gameID]; !ok {
			return gameID
		}
	}
}
