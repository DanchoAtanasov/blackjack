package messages

import "encoding/json"
import "blackjack/models"

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

func PlayerHandMessage(hand models.Hand) string {
	return makeJsonMessage("PlayerHand", hand.ToJson())
}

func DealerHandMessage(hand models.Hand) string {
	return makeJsonMessage("DealerHand", hand.ToJson())
}

var (
	START_MSG     string = gameStart() // {"type":"Game","message":"Start"}
	OVER_MSG      string = gameOver()  // {"type":"Game","message":"Over"}
	HIT_MSG       string = hit()       // {"type":"PlayerAction","message":"Hit"}
	STAND_MSG     string = stand()     // {"type":"PlayerAction","message":"Stand"}
	BUST_MSG      string = bust()      // {"type":"HandState","message":"Bust"}
	BLACKJACK_MSG string = blackjack() // {"type":"HandState","message":"Blackjack"}
	// dealer hand
	// player hand
)
