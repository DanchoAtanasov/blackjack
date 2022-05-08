package models

import "testing"

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
		if result != 2 {
			t.Errorf("GetWinner is %d but should be %d\n", result, 2)
		}

		handA = Hand{Sum: 15}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		if result != 1 {
			t.Errorf("GetWinner is %d but should be %d\n", result, 1)
		}

		handA = Hand{Sum: 13}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		if result != -1 {
			t.Errorf("GetWinner is %d but should be %d\n", result, -1)
		}

		handA = Hand{Sum: 22}
		handB = Hand{Sum: 14}
		result = GetWinner(handA, handB)
		if result != -1 {
			t.Errorf("GetWinner is %d but should be %d\n", result, -1)
		}

		handA = Hand{Sum: 15}
		handB = Hand{Sum: 22}
		result = GetWinner(handA, handB)
		if result != 1 {
			t.Errorf("GetWinner is %d but should be %d\n", result, 1)
		}

	})
}
