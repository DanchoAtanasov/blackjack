package models

import (
	"math"
	"math/rand"
	"time"
)

var cardValues = [...]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var suites = [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}
var cardValuesMap = make(map[string]int)

func init() {
	for idx, val := range cardValues {
		if val == "A" {
			cardValuesMap[val] = 11
		} else {
			cardValuesMap[val] = int(math.Min(float64(idx+2), 10))
		}
	}
}

type deckInterface interface {
	BuildDeck()
	DealCard() Card
}

type Deck struct {
	cards []Card
}

func (deck *Deck) BuildDeck() {
	for _, suit := range suites {
		for _, val := range cardValues {
			deck.cards = append(deck.cards, Card{val, suit, cardValuesMap[val]})
		}
	}
}

func (deck *Deck) DealCard() Card {
	dealtCard := deck.cards[0]
	deck.cards = deck.cards[1:]
	return dealtCard
}

func ShuffleDeck(deck Deck) {
	source := rand.NewSource(time.Now().UnixNano())
	randGenerator := rand.New(source)
	for i := range deck.cards {
		j := randGenerator.Intn(i + 1)
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	}
}

func ShanoShuffleDeck(deck *Deck) {
	deck.cards[0] = Card{valueStr: "A", suit: "Spades", value: 11}
	deck.cards[1] = Card{valueStr: "A", suit: "Clubs", value: 11}
	deck.cards[2] = Card{valueStr: "K", suit: "Spades", value: 10}
}
