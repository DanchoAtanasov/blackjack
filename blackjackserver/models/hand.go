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
	Cards       []Card `json:"cards"`
	Sum         int    `json:"sum"`
	NumAces     int    `json:"-"` // is used to implement soft count
	IsBlackjack bool   `json:"-"`
}

func (hand *Hand) AddCard(card Card) {
	if card.ValueStr == "A" {
		hand.NumAces += 1
	}
	hand.Cards = append(hand.Cards, card)
	hand.increaseSum(card.value) // would remove from NumAces if it's about to bust
	if hand.Sum == 21 && len(hand.Cards) == 2 {
		hand.IsBlackjack = true
	}
}

func (hand *Hand) RemoveCard() Card {
	// Think about improving this
	if len(hand.Cards) <= 0 {
		panic("Can't remove from empty hand")
	}

	var card Card
	hand.Cards, card = hand.Cards[:len(hand.Cards)-1], hand.Cards[len(hand.Cards)-1]
	// If there is an ace in the cards, set NumAces to 1, otherwise the NumAces count is wrong
	// when you split aces
	for i := 0; i < len(hand.Cards); i++ {
		if hand.Cards[i].ValueStr == "A" {
			hand.NumAces = 1
			break
		}
	}

	hand.increaseSum(-card.value)
	if hand.Sum == 21 && len(hand.Cards) == 2 {
		hand.IsBlackjack = true
	}
	return card
}

func (hand *Hand) increaseSum(value int) {
	hand.Sum += value
	// convert soft count to hard count
	if hand.Sum > 21 && hand.NumAces > 0 {
		hand.Sum -= 10
		hand.NumAces -= 1
	}
}

func (hand Hand) IsBust() bool {
	return hand.Sum > 21
}

func (hand Hand) CanSplit() bool {
	if len(hand.Cards) != 2 {
		return false
	}
	return hand.Cards[0].ValueStr == hand.Cards[1].ValueStr
}

func (hand *Hand) ToJson() string {
	result, _ := json.Marshal(*hand)
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
