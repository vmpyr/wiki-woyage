package structs

import (
	"sync"
)

type Game struct {
	LobbyID         string
	PlayerData      map[string]PlayerData // playerID -> PlayerData
	RoundData       RoundData
	Settings        GameSettings
	VotesToEnd      map[string]bool // playerID -> vote
	Finished        bool
	GameStructMutex sync.Mutex
}

type PlayerData struct {
	PlayerID         string
	InGame           bool
	CurrentArticle   map[string]string   // { "title": "Article Title", "slug": "article-slug" }
	ArticleTree      []map[string]string // [{ "title": "Article Title", "slug": "article-slug" }, ...]
	HasFinishedRound bool
	Score            float64
}

type RoundData struct {
	RoundNumber        int
	TimeElapsed        int
	OriginArticle      map[string]string // { "title": "Article Title", "slug": "article-slug" }
	DestinationArticle map[string]string // { "title": "Article Title", "slug": "article-slug" }
	Finished           bool              // For checking if new round should be started
}

type GameSettings struct {
	GameType            string // "fastest" "min steps"
	ScoringType         string // "points" "first only" "time" "steps"
	RoundWaitTime       int
	MaxGameDuration     int
	AllowCountries      bool
	AllowFindInArticles bool
}
