package messages

import "encoding/json"

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func makeJsonMessage(type_ string, message_ string) string {
	newMessage := Message{
		Type:    type_,
		Message: message_,
	}
	messageBytes, _ := json.Marshal(newMessage)
	return string(messageBytes)

}

func makeActionMessage(action string) string {
	return makeJsonMessage("PlayerAction", action)
}

func stand() string {
	return makeActionMessage("Stand")
}

func hit() string {
	return makeActionMessage("Hit")
}

func makeHandStateMessage(state string) string {
	return makeJsonMessage("HandState", state)
}

func bust() string {
	return makeHandStateMessage("Bust")
}

func blackjack() string {
	return makeHandStateMessage("Blackjack")
}

func makeGameMessage(message string) string {
	return makeJsonMessage("Game", message)
}

func gameOver() string {
	return makeGameMessage("Over")
}

func gameStart() string {
	return makeGameMessage("Start")
}

func makeHandMessage(hand string) {

}

// TODO: convert messages to json
var (
	STAND_MSG     string = stand()
	HIT_MSG       string = hit()
	BUST_MSG      string = bust()
	BLACKJACK_MSG string = blackjack()
	OVER_MSG      string = gameOver()
	START_MSG     string = gameStart()
	// dealer hand
	// player hand
)
