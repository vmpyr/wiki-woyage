package game

import (
	"errors"
	"sync"
	st "wiki-woyage/structs"
	"wiki-woyage/utils"
)

var (
	games            = make(map[string]*st.Game) // gameID -> Game
	lobbiesWithGames = make(map[string]string)   // lobbyID -> gameID
	gameMutex        = sync.RWMutex{}
)

func GetGame(gameID string) (*st.Game, error) {
	gameMutex.RLock()
	defer gameMutex.RUnlock()

	game, ok := games[gameID]
	if !ok {
		return nil, errors.New("game not found")
	}

	return game, nil
}

func StartNewGame(lobbyID string, settings st.GameSettings) (string, error) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if _, ok := lobbiesWithGames[lobbyID]; ok {
		return "", errors.New("game already ok for this lobby")
	}

	gameID := utils.GenerateGameID(&games)
	// TODO: Better handling of game start with a lobby struct to check "ready" players
	game := &st.Game{
		LobbyID:    lobbyID,
		PlayerData: make(map[string]st.PlayerData),
		RoundData: st.RoundData{
			TimeElapsed:        0,
			RoundNumber:        0,
			OriginArticle:      nil,
			DestinationArticle: nil,
			Finished:           true,
		},
		Settings:   settings,
		VotesToEnd: make(map[string]bool),
		Finished:   false,
	}
	games[gameID] = game
	lobbiesWithGames[lobbyID] = gameID

	return gameID, nil
}

func JoinGame(lobbyID, playerID string) (string, error) {
	gameMutex.RLock()
	gameID, ok := lobbiesWithGames[lobbyID]
	gameMutex.RUnlock()
	if !ok {
		return "", errors.New("no game found for this lobby")
	}

	game, err := GetGame(gameID)
	if err != nil {
		return "", err
	}

	(*game).GameStructMutex.Lock()
	defer (*game).GameStructMutex.Unlock()

	if playerData, ok := (*game).PlayerData[playerID]; ok {
		if !playerData.InGame {
			return "", errors.New("player already in game")
		} else {
			playerData.InGame = true
			(*game).PlayerData[playerID] = playerData
			return gameID, nil
		}
	}

	(*game).PlayerData[playerID] = st.PlayerData{
		PlayerID:         playerID,
		InGame:           true,
		CurrentArticle:   (*game).RoundData.OriginArticle,
		ArticleTree:      []map[string]string{(*game).RoundData.OriginArticle},
		HasFinishedRound: false,
		Score:            0.0,
	}

	return gameID, nil
}

func VoteToEndGame(gameID, playerID string) error {
	game, err := GetGame(gameID)
	if err != nil {
		return err
	}

	(*game).GameStructMutex.Lock()
	defer (*game).GameStructMutex.Unlock()

	if _, ok := (*game).PlayerData[playerID]; !ok {
		return errors.New("player not found in game to vote")
	}

	(*game).VotesToEnd[playerID] = true

	inGameCount := 0
	for _, playerData := range (*game).PlayerData {
		if playerData.InGame {
			inGameCount++
		}
	}
	if len((*game).VotesToEnd) >= inGameCount/2 {
		(*game).Finished = true
	}

	return nil
}
