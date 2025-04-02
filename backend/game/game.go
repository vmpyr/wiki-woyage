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
	gameMutex.Lock()
	gameID := utils.GenerateGameID(&games)
	gameMutex.Unlock()

	lobby, _ := lobby.GetLobby(lobbyID)
	err := utils.MutexExecLobby(lobby, func(l *st.Lobby) error {
		if l.GameID != "" {
			log.Printf("Lobby %s already has a game %s", lobbyID, l.GameID)
			return errors.New("lobby already has a game")
		}
		if l.AdminPlayerID != playerID {
			log.Printf("Player %s is not admin of lobby %s", playerID, lobbyID)
			return errors.New("player is not admin of lobby")
		}
		l.GameID = gameID
		return nil
	})
	if err != nil {
		return nil, err
	}

	game := &st.Game{
		GameID:  gameID,
		LobbyID: lobbyID,
		GamePlayerData: map[string]st.GamePlayerData{playerID: {
			GameID:           gameID,
			PlayerID:         playerID,
			InGame:           true,
			CurrentArticle:   nil,
			ArticleTree:      []map[string]string{},
			HasFinishedRound: true,
			VotedToEnd:       false,
			Score:            0.0,
		}},
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
	err := utils.MutexExecGame(game, func(g *st.Game) error {
		if _, ok := g.GamePlayerData[playerID]; ok {
			log.Printf("Player %s already in game %s", playerID, gameID)
			return errors.New("player already in game")
		}
		if playerLobbyID, _ := player.GetPlayerLobbyID(playerID); playerLobbyID != g.LobbyID {
			log.Printf("Player %s is not in the same lobby as game %s", playerID, gameID)
			return errors.New("player not in the same lobby as game")
		}
		if g.Finished {
			log.Printf("Game %s already finished", gameID)
			return errors.New("game already finished")
		}
		g.GamePlayerData[playerID] = st.GamePlayerData{
			GameID:           gameID,
			PlayerID:         playerID,
			InGame:           true,
			CurrentArticle:   nil,
			ArticleTree:      []map[string]string{},
			HasFinishedRound: true,
			VotedToEnd:       false,
			Score:            0.0,
		}
		return nil
	})
	if err != nil {
		if err.Error() == "player already in game" {
			return nil
		}
		return err
	}

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
