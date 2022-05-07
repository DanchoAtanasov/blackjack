package models

import (
	"testing"
)

func TestDeck(t *testing.T) {
	assertCorrectMessage := func(t testing.TB, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got %d want %d\n", got, want)
		}
	}

	t.Run("check empty deck size", func(t *testing.T) {
		deck := Deck{}
		numCards := len(deck.cards)
		expected := 0
		assertCorrectMessage(t, numCards, expected)
	})

	t.Run("check new deck size", func(t *testing.T) {
		numDecksInShoe := 6
		deck := GetNewDeck(numDecksInShoe)
		numCards := len(deck.cards)
		expected := numDecksInShoe * 52

		assertCorrectMessage(t, numCards, expected)
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
}
