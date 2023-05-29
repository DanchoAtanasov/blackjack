package config

import "time"

// Blackjack game settings
const NumDecksInShoe int = 6
const ShuffleThreshold int = 150

// General game settings
const RoomSize int = 2
const TimeBetweenRounds time.Duration = 5 * time.Second

// Websocket settings
const ReadTimeout time.Duration = 1000 * time.Second

// Run settings
const (
	PlayMode int = iota
	AuditMode
)
const Mode = PlayMode
const AuditLogFile = "./audit/3f987ada-4c95-4acf-a466-1952085d7e5b.log"
