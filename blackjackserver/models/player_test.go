package models

import "testing"

func TestPlayer(t *testing.T) {
	testWinOrLose := func(t testing.TB, startingBuyIn, bet int, isWin bool, isBlackjack bool) {
		t.Helper()
		player := Player{
			Name:    "Test Player",
			BuyIn:   startingBuyIn,
			CurrBet: 1,
		}
		var expected int

		if isBlackjack {
			player.Blackjack()
			expected = startingBuyIn + 2*player.CurrBet
		} else if isWin {
			player.Win()
			expected = startingBuyIn + player.CurrBet
		} else {
			player.Lose()
			expected = startingBuyIn - player.CurrBet
		}

		if player.BuyIn != expected {
			t.Errorf("player's buying after hand is %d but should be %d\n",
				player.BuyIn, expected)
		}
	}

	t.Run("test player win or lose", func(t *testing.T) {
		testSuite := []struct {
			buyIn       int
			bet         int
			isWin       bool
			isBlackjack bool
		}{
			{10, 1, true, false},
			{20, 2, true, false},
			{30, 10, false, false},
			{40, 5, false, false},
			{50, 1, false, false},
			{50, 1, true, true},
		}
		for _, test := range testSuite {
			testWinOrLose(t, test.buyIn, test.bet, test.isWin, test.isBlackjack)
		}
	})
}
