package models

import "encoding/json"

type winOrLose interface {
	Win()
	Lose()
}

type Player struct {
	Name       string
	BuyIn      int
	Hand       Hand
	CurrentBet int
}

func (player *Player) Win() {
	player.BuyIn += player.CurrentBet
}

func (player *Player) Lose() {
	player.BuyIn -= player.CurrentBet
}

func (player *Player) ToJson() string {
	result, _ := json.Marshal(*player)
	return string(result)
}
