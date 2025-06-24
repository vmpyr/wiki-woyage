package main

import (
	"math/rand"
)

func GenerateID(characters string, length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func GenerateLobbyID(lobbies *LobbyList) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		lobbyID := GenerateID(characters, 5)
		if _, ok := (*lobbies)[lobbyID]; !ok {
			return lobbyID
		}
	}
}

func CheckUniqueUsername(username string, players *PlayerList) bool {
	for player := range *players {
		if player.username == username {
			return false
		}
	}
	return true
}
