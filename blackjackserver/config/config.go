package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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

func getMode() int {
	if GetEnv("MODE", "PLAY") == "AUDIT" {
		fmt.Println("AUDIT MODE Son")
		return AuditMode
	}
	return PlayMode
}

func GetSeed() int64 {
	fmt.Println("Using audit seed")
	seed, err := strconv.ParseInt(GetEnv("SEED", "NONE"), 10, 64)
	if err != nil {
		panic("Seed not read correctly")
	}
	return seed
}

var Mode = getMode()
