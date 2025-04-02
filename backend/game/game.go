package game

import (
	"errors"
	"log"
	"sync"
	"wiki-woyage/lobby"
	"wiki-woyage/player"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"

	"github.com/gorilla/websocket"
)

var (
	games     = make(map[string]*st.Game) // gameID -> Game
	gameMutex = sync.RWMutex{}
)

func StartNewGame(conn *websocket.Conn, playerID string, lobbyID string, gameSettings st.GameSettings) (*st.Game, error) {
	lobby, _ := lobby.GetLobby(lobbyID)
	lobby.LobbyStructMutex.Lock()
	if lobby.GameID != "" {
		log.Printf("Lobby %s already has a game %s", lobbyID, lobby.GameID)
		lobby.LobbyStructMutex.Unlock()
		return nil, errors.New("lobby already has a game")
	}
	if lobby.AdminPlayerID != playerID {
		log.Printf("Player %s is not admin of lobby %s", playerID, lobbyID)
		lobby.LobbyStructMutex.Unlock()
		return nil, errors.New("player is not admin of lobby")
	}
	lobby.LobbyStructMutex.Unlock()

	gameMutex.Lock()
	gameID := utils.GenerateGameID(&games)
	gameMutex.Unlock()

	lobby.LobbyStructMutex.Lock()
	lobby.GameID = gameID
	lobby.LobbyStructMutex.Unlock()

	game := &st.Game{
		GameID:         gameID,
		LobbyID:        lobbyID,
		GamePlayerData: make(map[string]st.GamePlayerData),
		RoundData: st.RoundData{
			GameID:             gameID,
			RoundNumber:        0,
			TimeElapsed:        0,
			OriginArticle:      nil,
			DestinationArticle: nil,
			Finished:           true,
		},
		Settings:        gameSettings,
		Finished:        false,
		Conn:            conn,
		GameStructMutex: sync.Mutex{},
	}

	game.GamePlayerData[playerID] = st.GamePlayerData{
		GameID:           gameID,
		PlayerID:         playerID,
		InGame:           true,
		CurrentArticle:   nil,
		ArticleTree:      []map[string]string{},
		HasFinishedRound: true,
		VotedToEnd:       false,
		Score:            0.0,
	}
	gameMutex.Lock()
	games[gameID] = game
	gameMutex.Unlock()

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	player.GameID = gameID
	player.PlayerStructMutex.Unlock()

	log.Printf("Game %s created for lobby %s by player %s", gameID, lobbyID, playerID)
	return game, nil
}

func JoinGame(playerID, gameID string) error {
	game, _ := GetGame(gameID)

	game.GameStructMutex.Lock()
	if _, ok := game.GamePlayerData[playerID]; ok {
		log.Println("Player already in game")
		game.GameStructMutex.Unlock()
		return nil
	}
	if playerLobbyID, _ := player.GetPlayerLobbyID(playerID); playerLobbyID != game.LobbyID {
		log.Printf("Player %s is not in the same lobby as game %s", playerID, gameID)
		game.GameStructMutex.Unlock()
		return errors.New("player not in the same lobby as game")
	}
	game.GamePlayerData[playerID] = st.GamePlayerData{
		GameID:           gameID,
		PlayerID:         playerID,
		InGame:           true,
		CurrentArticle:   nil,
		ArticleTree:      []map[string]string{},
		HasFinishedRound: true,
		VotedToEnd:       false,
		Score:            0.0,
	}
	game.GameStructMutex.Unlock()

	lobby, _ := lobby.GetLobby(game.LobbyID)
	lobby.LobbyStructMutex.Lock()
	lobby.PlayerIDs[playerID] = true
	lobby.LobbyStructMutex.Unlock()

	player, _ := player.GetPlayer(playerID)
	player.PlayerStructMutex.Lock()
	player.GameID = gameID
	player.PlayerStructMutex.Unlock()

	log.Printf("Player %s joined game %s", playerID, gameID)
	return nil
}

func VoteToEndGame(playerID, gameID string) error {
	game, _ := GetGame(gameID)

	game.GameStructMutex.Lock()
	if playerData, ok := game.GamePlayerData[playerID]; !ok {
		log.Printf("Player %s not in game %s", playerID, gameID)
		game.GameStructMutex.Unlock()
		return errors.New("player not in game")
	} else {
		playerData.VotedToEnd = true
		game.GamePlayerData[playerID] = playerData
	}
	playersInGame := 0
	votesToEnd := 0
	for _, playerData := range game.GamePlayerData {
		if playerData.InGame {
			playersInGame++
			if playerData.VotedToEnd {
				votesToEnd++
			}
		}
	}
	game.GameStructMutex.Unlock()

	if votesToEnd >= playersInGame/2 {
		game.GameStructMutex.Lock()
		game.Finished = true
		game.GameStructMutex.Unlock()

		lobby, _ := lobby.GetLobby(game.LobbyID)
		lobby.LobbyStructMutex.Lock()
		lobby.GameID = ""
		lobby.LobbyStructMutex.Unlock()

		player, _ := player.GetPlayer(playerID)
		player.PlayerStructMutex.Lock()
		player.GameID = ""
		player.PlayerStructMutex.Unlock()

		log.Printf("Game %s ended by majority", gameID)
	}
	return nil
}

func GetGame(gameID string) (*st.Game, error) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()

	game, ok := games[gameID]
	if !ok {
		return nil, errors.New("game not found")
	}

	return game, nil
}
