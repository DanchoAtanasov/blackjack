package models

import "encoding/json"

type winOrLose interface {
	Win()
	Lose()
	Blackjack()
}

type Player struct {
	Name     string
	BuyIn    int
	CurrBet  int
	Hand     Hand
	IsDealer bool `json:"-"`
	Active   bool `json:"-"`
}

func (player *Player) Win() {
	player.BuyIn += player.CurrBet
}

func (player *Player) Lose() {
	player.BuyIn -= player.CurrBet
}

func (player *Player) Blackjack() {
	player.Win()
	player.Win()
}

func (player *Player) ToJson() string {
	result, _ := json.Marshal(*player)
	return string(result)
}
