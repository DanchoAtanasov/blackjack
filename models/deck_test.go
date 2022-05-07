package models

import (
	"testing"
)

func TestDeck(t *testing.T) {
	assertIntEqual := func(t testing.TB, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got %d want %d\n", got, want)
		}
	}

	t.Run("check empty deck size", func(t *testing.T) {
		deck := Deck{}
		numCards := len(deck.cards)
		expected := 0
		assertIntEqual(t, numCards, expected)
	})

	t.Run("check new deck size", func(t *testing.T) {
		numDecksInShoe := 6
		deck := GetNewDeck(numDecksInShoe)
		numCards := len(deck.cards)
		expected := numDecksInShoe * 52

		assertIntEqual(t, numCards, expected)
	})

	t.Run("check two shuffled decks are different", func(t *testing.T) {
		numDecksInShoe := 6
		firstDeck := GetNewShuffledDeck(numDecksInShoe)
		secondDeck := GetNewShuffledDeck(numDecksInShoe)

		for i := 0; i < numDecksInShoe*52; i++ {
			if firstDeck.cards[i] != secondDeck.cards[i] {
				return
			}
		}
		t.Error("decks are the same, should be shuffled")
	})

	t.Run("test deal a card", func(t *testing.T) {
		numDecksInShoe := 6
		deck := GetNewDeck(numDecksInShoe)

		dealtCard := deck.DealCard()
		expectedCard := Card{
			ValueStr: "2",
			Suit:     "Spades",
			value:    2,
		}
		if dealtCard != expectedCard {
			t.Errorf("dealt card is incorrect: got %v, want %v \n", dealtCard, expectedCard)
		}

		dealtCard = deck.DealCard()
		expectedCard = Card{
			ValueStr: "3",
			Suit:     "Spades",
			value:    3,
		}
		if dealtCard != expectedCard {
			t.Errorf("dealt card is incorrect: got %v, want %v \n", dealtCard, expectedCard)
		}
	})

	t.Run("test shuffle deck if low", func(t *testing.T) {
		numDecksInShoe := 6
		deck := GetNewDeck(numDecksInShoe)
		threshold := 200

		deck.cards = deck.cards[:threshold-1]
		newDeck := ShuffleDeckIfLow(&deck, threshold)

		if newDeck == &deck {
			t.Error("deck pointer is the same but should be different")
		}

		for i := 0; i < threshold-1; i++ {
			if deck.cards[i] != newDeck.cards[i] {
				return
			}
		}
		t.Error("deck is the same after shuffling")
	})

	t.Run("test don't shuffle deck if not low", func(t *testing.T) {
		numDecksInShoe := 6
		deck := GetNewDeck(numDecksInShoe)
		threshold := 200

		deck.cards = deck.cards[:threshold+1]
		newDeck := ShuffleDeckIfLow(&deck, threshold)

		if newDeck != &deck {
			t.Error("deck pointer is different but should be the same")
		}

		for i := 0; i < threshold+1; i++ {
			if deck.cards[i] != newDeck.cards[i] {
				t.Error("deck is different after shuffling")
			}
		}
	})
}
