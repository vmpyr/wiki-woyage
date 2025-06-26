package main

import "time"

type GameSettings struct {
	gameType    string
	maxTime     int // minutes
	totalRounds int
	roundTime   int // minutes
}

type Game struct {
	settings  GameSettings
	scores    map[string]int // clientID -> score
	startTime time.Time
	roundData map[string]RoundData // roundID -> roundData
}

type RoundData struct {
	roundID      string
	scores       map[string]int      // clientID -> score
	startArticle map[string]string   // clientID -> article slug
	endArticle   map[string]string   // clientID -> article slugs
	articleTree  map[string][]string // clientID -> article slug list
	startTime    time.Time
	endTime      time.Time
}

const (
	TimeAttack = "time_attack"
	MinSteps   = "min_steps"
)
