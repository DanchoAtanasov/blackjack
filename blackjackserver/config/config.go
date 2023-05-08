package config

import "time"

const NumDecksInShoe int = 6
const ShuffleThreshold int = 150
const RoomSize int = 2
const NumRoundsPerGame int = 3
const InitialBuyIn int = 100
const CurrBet int = 1
const TimeBetweenRounds time.Duration = 5 * time.Second
const ReadTimeout time.Duration = 1000 * time.Second
