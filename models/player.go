package models

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
