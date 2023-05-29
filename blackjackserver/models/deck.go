package models

import (
	"math"
	"math/rand"
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
	randSource    rand.Source
	randGenerator *rand.Rand
	size          int
	cards         []Card
}

func GetNewDeck(size int, seed int64) Deck {
	deck := Deck{}
	deck.randSource = rand.NewSource(seed)
	deck.randGenerator = rand.New(deck.randSource)
	deck.size = size
	deck.BuildDeck()

	return deck
}

func GetNewShuffledDeck(size int, seed int64) Deck {
	deck := GetNewDeck(size, seed)
	ShuffleDeck(deck)
	return deck
}

// Initialize deck with card values
func (deck *Deck) BuildDeck() {
	deck.cards = make([]Card, 0, deck.size)
	for i := 0; i < deck.size; i++ {
		for _, suit := range suites {
			for _, val := range cardValues {
				deck.cards = append(deck.cards, Card{val, suit, cardValuesMap[val]})
			}
		}
	}
}

func (deck *Deck) DealCard() Card {
	var dealtCard Card
	dealtCard, deck.cards = deck.cards[0], deck.cards[1:]
	return dealtCard
}

func ShuffleDeck(deck Deck) {
	for i := range deck.cards {
		j := deck.randGenerator.Intn(i + 1)
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	}
}

func ShuffleDeckIfLow(deck *Deck, threshold int) {
	if len(deck.cards) > threshold {
		return
	}
	deck.BuildDeck()
	ShuffleDeck(*deck)
}

// Used for testing specific behaviour
func ShanoShuffleDeck(deck *Deck) {
	deck.cards[0] = Card{ValueStr: "A", Suit: "Spades", value: 11}
	deck.cards[1] = Card{ValueStr: "A", Suit: "Clubs", value: 11}
	deck.cards[2] = Card{ValueStr: "K", Suit: "Spades", value: 10}
}
