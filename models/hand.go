package models

import (
	"encoding/json"
)

type handInterface interface {
	AddCard(Card)
	increaseSum(int)
	IsBust() bool
}

type Hand struct {
	Cards       []Card
	Sum         int
	NumAces     int
	IsBlackjack bool
}

func (hand *Hand) AddCard(card Card) {
	if card.ValueStr == "A" {
		hand.NumAces += 1
	}
	hand.Cards = append(hand.Cards, card)
	hand.increaseSum(card.value)
}

func (hand *Hand) increaseSum(value int) {
	hand.Sum += value
	if hand.Sum > 21 && hand.NumAces > 0 {
		hand.Sum -= 10
		hand.NumAces -= 1
	}
}

func (hand Hand) IsBust() bool {
	return hand.Sum > 21
}

func (hand *Hand) ToJson() string {
	result, _ := json.Marshal(hand.Cards)
	return string(result)
}

func GetWinner(handA Hand, handB Hand) int {
	// 1 if greater, -1 if smaller, 0 if equal, 2 for Blackjack
	if handA.IsBlackjack {
		return 2
	}
	if handA.Sum > 21 { // Player bust
		return -1
	}
	if handB.Sum > 21 { // Dealer bust
		return 1
	}
	if handA.Sum > handB.Sum {
		return 1
	}
	if handA.Sum < handB.Sum {
		return -1
	}
	return 0
}

func ClearHand(hand *Hand) {
	*hand = Hand{}
}
