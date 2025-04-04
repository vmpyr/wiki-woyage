package game

import (
	"errors"
	"log"
	"sync"
	"wiki-woyage/lobby"
	"wiki-woyage/player"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
)

var (
	games     = make(map[string]*st.Game) // gameID -> Game
	gameMutex = sync.RWMutex{}
)

func StartNewGame(p *st.Player, lobbyID string, gameSettings st.GameSettings) (string, error) {
	gameMutex.Lock()
	gameID := utils.GenerateGameID(&games)
	gameMutex.Unlock()

	lobby, _ := lobby.GetLobby(lobbyID)
	if err := utils.MutexExecLobby(lobby, func(l *st.Lobby) error {
		if l.GameID != "" {
			log.Printf("Lobby %s already has a game %s", lobbyID, l.GameID)
			return errors.New("lobby already has a game")
		}
		if l.AdminPlayerID != p.PlayerID {
			log.Printf("Player %s is not admin of lobby %s", p.PlayerID, lobbyID)
			return errors.New("player is not admin of lobby")
		}
		l.GameID = gameID
		return nil
	}); err != nil {
		return "", err
	}

	game := &st.Game{
		GameID:  gameID,
		LobbyID: lobbyID,
		GamePlayerData: map[string]st.GamePlayerData{p.PlayerID: {
			GameID:           gameID,
			PlayerID:         p.PlayerID,
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
		GameStructMutex: sync.Mutex{},
	}

	gameMutex.Lock()
	games[gameID] = game
	gameMutex.Unlock()

	p.PlayerStructMutex.Lock()
	p.GameID = gameID
	p.PlayerStructMutex.Unlock()

	log.Printf("Game %s created for lobby %s by player %s", gameID, lobbyID, p.PlayerID)
	return gameID, nil
}

func JoinGame(p *st.Player, gameID string) error {
	game, _ := GetGame(gameID)
	if err := utils.MutexExecGame(game, func(g *st.Game) error {
		if _, ok := g.GamePlayerData[p.PlayerID]; ok {
			log.Printf("Player %s already in game %s", p.PlayerID, gameID)
			return errors.New("player already in game")
		}
		if playerLobbyID, _ := player.GetPlayerLobbyID(p.PlayerID); playerLobbyID != g.LobbyID {
			log.Printf("Player %s is not in the same lobby as game %s", p.PlayerID, gameID)
			return errors.New("player not in the same lobby as game")
		}
		if g.Finished {
			log.Printf("Game %s already finished", gameID)
			return errors.New("game already finished")
		}
		g.GamePlayerData[p.PlayerID] = st.GamePlayerData{
			GameID:           gameID,
			PlayerID:         p.PlayerID,
			InGame:           true,
			CurrentArticle:   nil,
			ArticleTree:      []map[string]string{},
			HasFinishedRound: true,
			VotedToEnd:       false,
			Score:            0.0,
		}
		return nil
	}); err != nil {
		if err.Error() == "player already in game" {
			return nil
		}
		return err
	}

	lobby, _ := lobby.GetLobby(game.LobbyID)
	lobby.LobbyStructMutex.Lock()
	lobby.PlayerIDs[p.PlayerID] = true
	lobby.LobbyStructMutex.Unlock()

	p.PlayerStructMutex.Lock()
	p.GameID = gameID
	p.PlayerStructMutex.Unlock()

	log.Printf("Player %s joined game %s", p.PlayerID, gameID)
	return nil
}

func VoteToEndGame(p *st.Player, gameID string) error {
	game, _ := GetGame(gameID)

	game.GameStructMutex.Lock()
	if playerData, ok := game.GamePlayerData[p.PlayerID]; !ok {
		log.Printf("Player %s not in game %s", p.PlayerID, gameID)
		game.GameStructMutex.Unlock()
		return errors.New("player not in game")
	} else {
		playerData.VotedToEnd = true
		game.GamePlayerData[p.PlayerID] = playerData
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

		p.PlayerStructMutex.Lock()
		p.GameID = ""
		p.PlayerStructMutex.Unlock()

		log.Printf("Game %s ended by majority", gameID)
	}
	return nil
}

func LeaveGame(p *st.Player, gameID string, playerDisconnect bool) error {
	if !playerDisconnect {
		if err := utils.MutexExecPlayer(p, func(p *st.Player) error {
			if p.GameID != gameID {
				log.Printf("Player %s not in game %s", p.PlayerID, gameID)
				return errors.New("player not in game")
			}
			p.GameID = ""
			return nil
		}); err != nil {
			return err
		}
	}

	game, _ := GetGame(gameID)
	return utils.MutexExecGame(game, func(g *st.Game) error {
		delete(g.GamePlayerData, p.PlayerID)
		log.Printf("Player %s left game %s", p.PlayerID, gameID)

		if len(g.GamePlayerData) == 0 {
			gameMutex.Lock()
			delete(games, gameID)
			gameMutex.Unlock()
			log.Printf("No player in game %s deleted, hence deleted", gameID)
		}
		return nil
	})
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
