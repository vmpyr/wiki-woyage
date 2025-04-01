package utils

import (
	"math/rand"
	"os"
	"strconv"
	st "wiki-woyage/structs"

	"github.com/gorilla/websocket"
)

func GenerateID(characters string, length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func GenerateLobbyID(lobbies *map[string]*st.Lobby) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		lobbyID := GenerateID(characters, 5)
		if _, ok := (*lobbies)[lobbyID]; !ok {
			return lobbyID
		}
	}
}

func GeneratePlayerID(players *map[string]*st.Player) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for {
		playerID := GenerateID(characters, 25)
		if _, ok := (*players)[playerID]; !ok {
			return playerID
		}
	}
}

func GenerateGameID(games *map[string]*st.Game) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for {
		gameID := GenerateID(characters, 25)
		if _, ok := (*games)[gameID]; !ok {
			return gameID
		}
	}
}

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

func SendError(conn *websocket.Conn, errorMessage string) {
	response := st.WebSocketResponse{
		Type:         "error",
		ErrorMessage: errorMessage,
	}
	err := conn.WriteJSON(response)
	if err != nil {
		conn.Close()
	}
}

func SendResponse(conn *websocket.Conn, response st.WebSocketResponse) {
	err := conn.WriteJSON(response)
	if err != nil {
		conn.Close()
	}
}
