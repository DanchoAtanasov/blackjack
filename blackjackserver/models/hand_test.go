package models

import (
	"testing"
)

func TestHand(t *testing.T) {
	t.Run("test add card", func(t *testing.T) {
		hand := Hand{}
		if len(hand.Cards) != 0 {
			t.Error("new hand should be empty")
		}

		card := Card{
			ValueStr: "2",
			Suit:     "Spades",
			value:    2,
		}
		hand.AddCard(card)

		if len(hand.Cards) != 1 {
			t.Error("hand size should be 1")
		}
		if hand.Sum != card.value {
			t.Errorf("hand sum is %d but should be %d\n", hand.Sum, card.value)
		}

		aceCard := Card{
			ValueStr: "A",
			Suit:     "Spades",
			value:    11,
		}
		hand.AddCard(aceCard)

		if len(hand.Cards) != 2 {
			t.Error("hand size should be 2")
		}
		twoCardSum := card.value + aceCard.value
		if hand.Sum != twoCardSum {
			t.Errorf("hand sum is %d but should be %d\n", hand.Sum, twoCardSum)
		}
		if hand.NumAces != 1 {
			t.Errorf("hand num aces is %d but should be %d\n", hand.NumAces, 1)
		}
	})

	t.Run("test hand bust", func(t *testing.T) {
		checkBust := func(t testing.TB, hand *Hand, expected bool) {
			t.Helper()
			got := hand.IsBust()
			if got != expected {
				t.Errorf("isBust is %t but should be %t\n", got, expected)
			}
		}
		hand := Hand{}
		checkBust(t, &hand, false)
		hand.AddCard(Card{
			ValueStr: "K",
			Suit:     "Spades",
			value:    10,
		})
		checkBust(t, &hand, false)
		hand.AddCard(Card{
			ValueStr: "10",
			Suit:     "Spades",
			value:    10,
		})
		checkBust(t, &hand, false)
		hand.AddCard(Card{
			ValueStr: "Q",
			Suit:     "Spades",
			value:    10,
		})
		checkBust(t, &hand, true)
	})

	t.Run("test increase sum", func(t *testing.T) {
		hand := Hand{}
		if hand.Sum != 0 {
			t.Errorf("empty hand sum is %d should be 0\n", hand.Sum)
		}

		hand.increaseSum(10)
		if hand.Sum != 10 {
			t.Errorf("empty hand sum is %d should be 10\n", hand.Sum)
		}

		hand.NumAces = 1
		hand.increaseSum(14)
		if hand.Sum != 14 {
			t.Errorf("sum is %d should be %d\n", hand.Sum, 14)
		}
		if hand.NumAces != 0 {
			t.Errorf("numAces is %d should be %d\n", hand.NumAces, 0)
		}
	})

	t.Run("test clear hand", func(t *testing.T) {
		hand := Hand{}
		hand.AddCard(Card{
			ValueStr: "2",
			Suit:     "Spades",
			value:    2,
		})

		ClearHand(&hand)
		if len(hand.Cards) != 0 {
			t.Errorf("cards in hand is %d should be %d\n", len(hand.Cards), 0)
		}
		if hand.Sum != 0 {
			t.Errorf("sum is %d should be %d\n", hand.Sum, 0)
		}
	})

	t.Run("test get winner", func(t *testing.T) {
		handA := Hand{}
		handA.IsBlackjack = true
		handB := Hand{}
		result := GetWinner(handA, handB)
		expected := 2
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 15}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		expected = 1
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 13}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		expected = -1
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 22}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		expected = -1
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 15}
		handB = Hand{Sum: 22}
		result = GetWinner(handA, handB)
		expected = 1
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 17}
		handB = Hand{Sum: 17}
		result = GetWinner(handA, handB)
		expected = 0
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}

		handA = Hand{Sum: 21, IsBlackjack: false}
		handB = Hand{Sum: 21}
		result = GetWinner(handA, handB)
		expected = 0
		if result != expected {
			t.Errorf("GetWinner is %d but should be %d\n", result, expected)
		}
	})

	t.Run("test remove card", func(t *testing.T) {
		hand := Hand{}
		hand.AddCard(Card{
			ValueStr: "2",
			Suit:     "Spades",
			value:    2,
		})
		secondCard := Card{
			ValueStr: "3",
			Suit:     "Spades",
			value:    3,
		}
		hand.AddCard(secondCard)

		removedCard := hand.RemoveCard()

		if len(hand.Cards) != 1 {
			t.Errorf("cards in hand is %d should be %d\n", len(hand.Cards), 1)
		}
		if removedCard != secondCard {
			t.Error("Removed card isn't the last added one.")
		}
	})

	t.Run("test remove card check ace count", func(t *testing.T) {
		hand := Hand{}
		aceCard := Card{
			ValueStr: "A",
			Suit:     "Spades",
			value:    11,
		}
		hand.AddCard(aceCard)
		hand.AddCard(aceCard)

		if hand.NumAces != 1 {
			// Should be 1 even though we add 2 aces as 11 + 11 = 22 which is bust
			// so we do 22 - 10 = 12 sum and 1 NumAces
			t.Errorf("Num aces is %d should be %d\n", hand.NumAces, 1)
		}

		removedCard := hand.RemoveCard()

		if len(hand.Cards) != 1 {
			t.Errorf("cards in hand is %d should be %d\n", len(hand.Cards), 1)
		}
		if removedCard != aceCard {
			t.Error("Removed card isn't the last added one.")
		}
		if hand.NumAces != 1 {
			t.Errorf("Num aces is %d should be %d\n", hand.NumAces, 1)
		}
	})

	t.Run("test remove card from empty hand", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected RemoveCard to panic, but it didn't")
			}
		}()
		hand := Hand{}
		hand.RemoveCard()
	})
}
