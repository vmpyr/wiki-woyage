package structs

import (
	"sync"
)

type Game struct {
	LobbyID    string
	PlayerData map[string]PlayerData
	GameState  GameState
	Settings   GameSettings
	VotesToEnd map[string]bool
	GameMutex  sync.Mutex
}

type PlayerData struct {
	PlayerID       string
	InGame         bool
	CurrentArticle map[string]string   // { "title": "Article Title", "slug": "article-slug" }
	ArticleTree    []map[string]string // [{ "title": "Article Title", "slug": "article-slug" }, ...]
	HasFinished    bool
	Score          int
}

type GameState struct {
	TimeElapsed        int
	RoundNumber        int
	OriginArticle      map[string]string // { "title": "Article Title", "slug": "article-slug" }
	DestinationArticle map[string]string // { "title": "Article Title", "slug": "article-slug" }
	Finished           bool              // If not true, don't allow game settings to change
}

type GameSettings struct {
	GameType            string // "fastest" "min steps"
	ScoringType         string // "points" "first only" "time" "steps"
	RoundWaitTime       int
	MaxGameDuration     int
	AllowCountries      bool
	AllowFindInArticles bool
}
