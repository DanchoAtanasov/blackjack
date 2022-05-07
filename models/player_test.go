package models

import "testing"

func TestPlayer(t *testing.T) {
	testWinOrLose := func(t testing.TB, startingBuyIn, bet int, isWin bool) {
		t.Helper()
		player := Player{
			Name:       "Test Player",
			BuyIn:      startingBuyIn,
			CurrentBet: 1,
		}
		var expected int

		if isWin {
			player.Win()
			expected = startingBuyIn + player.CurrentBet
		} else {
			player.Lose()
			expected = startingBuyIn - player.CurrentBet
		}

		if player.BuyIn != expected {
			t.Errorf("player's buying after win is %d but should be %d\n",
				player.BuyIn, expected)
		}
	}

	t.Run("test player win or lose", func(t *testing.T) {
		testSuite := []struct {
			buyIn int
			bet   int
			isWin bool
		}{
			{10, 1, true},
			{20, 2, true},
			{30, 10, false},
			{40, 5, false},
			{50, 1, false},
		}
		for _, test := range testSuite {
			testWinOrLose(t, test.buyIn, test.bet, test.isWin)
		}
	})
}
