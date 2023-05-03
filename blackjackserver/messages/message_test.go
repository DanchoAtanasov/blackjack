package messages

import (
	"blackjack/models"
	"testing"
)

func TestMessage(t *testing.T) {
	assertStringEqual := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %s want %s\n", got, want)
		}
	}

	t.Run("test game start message", func(t *testing.T) {
		message := gameStart()
		expected := `{"type":"Game","message":"Start"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test game over message", func(t *testing.T) {
		message := gameOver()
		expected := `{"type":"Game","message":"Over"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test hit message", func(t *testing.T) {
		message := hit()
		expected := `{"type":"PlayerAction","message":"Hit"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test stand message", func(t *testing.T) {
		message := stand()
		expected := `{"type":"PlayerAction","message":"Stand"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test leave message", func(t *testing.T) {
		message := leave()
		expected := `{"type":"PlayerAction","message":"Leave"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test bust message", func(t *testing.T) {
		message := bust()
		expected := `{"type":"HandState","message":"Bust"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test blackjack message", func(t *testing.T) {
		message := blackjack()
		expected := `{"type":"HandState","message":"Blackjack"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test player hand message", func(t *testing.T) {
		player := models.Player{}
		player.Hand.AddCard(models.Card{
			ValueStr: "2",
			Suit:     "Spades",
		})
		player.Hand.AddCard(models.Card{
			ValueStr: "3",
			Suit:     "Clubs",
		})
		message := playerHandMessage(player)
		expected := `{"type":"PlayerHand","message":"{\"Name\":\"\",\"BuyIn\":0,\"Hand\":{\"cards\":[{\"ValueStr\":\"2\",\"Suit\":\"Spades\"},{\"ValueStr\":\"3\",\"Suit\":\"Clubs\"}],\"sum\":0},\"CurrentBet\":0}"}`
		assertStringEqual(t, message, expected)
	})

	t.Run("test dealer hand message", func(t *testing.T) {
		player := models.Player{}
		player.Hand.AddCard(models.Card{
			ValueStr: "2",
			Suit:     "Spades",
		})
		message := dealerHandMessage(player)
		expected := `{"type":"DealerHand","message":"{\"cards\":[{\"ValueStr\":\"2\",\"Suit\":\"Spades\"}],\"sum\":0}"}`
		assertStringEqual(t, message, expected)
	})
}
