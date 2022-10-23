package messages

import (
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
}
