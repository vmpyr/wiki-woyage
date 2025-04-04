package structs

import (
	"sync"
)

type Game struct {
	GameID          string
	LobbyID         string
	GamePlayerData  map[string]GamePlayerData // playerID -> PlayerData
	RoundData       RoundData
	Settings        GameSettings
	Finished        bool
	GameStructMutex sync.Mutex
}

type GamePlayerData struct {
	GameID           string
	PlayerID         string
	InGame           bool
	CurrentArticle   map[string]string   // { "title": "Article Title", "slug": "article-slug" }
	ArticleTree      []map[string]string // [{ "title": "Article Title", "slug": "article-slug" }, ...]
	HasFinishedRound bool
	VotedToEnd       bool
	Score            float64
}

type RoundData struct {
	GameID             string
	RoundNumber        int
	TimeElapsed        int
	OriginArticle      map[string]string // { "title": "Article Title", "slug": "article-slug" }
	DestinationArticle map[string]string // { "title": "Article Title", "slug": "article-slug" }
	Finished           bool              // For checking if new round should be started
}

type GameSettings struct {
	GameID              string
	GameType            string // "fastest" "min steps"
	ScoringType         string // "points" "first only" "time" "steps"
	RoundWaitTime       int
	MaxGameDuration     int
	AllowCountries      bool
	AllowFindInArticles bool
}
