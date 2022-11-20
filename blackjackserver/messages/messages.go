package messages

import (
	"blackjack/models"
	"encoding/json"
)

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

type handMessageFunc func(models.Player) string

func playerHandMessage(player models.Player) string {
	messageBytes, _ := json.Marshal(player)
	return makeJsonMessage("PlayerHand", string(messageBytes))
}

func DecodePlayerHandMessage(msg string) (models.Hand, error) {
	var message Message
	err := json.Unmarshal([]byte(msg), &message)
	var hand models.Hand
	err = json.Unmarshal([]byte(message.Message), &hand)
	return hand, err
}

func dealerHandMessage(player models.Player) string {
	return makeJsonMessage("DealerHand", player.Hand.ToJson())
}

type listPlayersFunc func([]*models.Player) string

func listPlayers(players []*models.Player) string {
	playerMessageBytes, _ := json.Marshal(players)
	return makeJsonMessage("ListPlayers", string(playerMessageBytes))
}

// TODO: Add decode for dealer hand message

var (
	START_MSG        string          = gameStart()       // {"type":"Game","message":"Start"}
	OVER_MSG         string          = gameOver()        // {"type":"Game","message":"Over"}
	HIT_MSG          string          = hit()             // {"type":"PlayerAction","message":"Hit"}
	STAND_MSG        string          = stand()           // {"type":"PlayerAction","message":"Stand"}
	BUST_MSG         string          = bust()            // {"type":"HandState","message":"Bust"}
	BLACKJACK_MSG    string          = blackjack()       // {"type":"HandState","message":"Blackjack"}
	LIST_PLAYERS_MSG listPlayersFunc = listPlayers       // {"type":"ListPlayers","message":"[]"}
	PLAYER_HAND_MSG  handMessageFunc = playerHandMessage // {"type":"PlayerHand","message":""}
	DEALER_HAND_MSG  handMessageFunc = dealerHandMessage // {"type":"DealerHand","message":""}

	// dealer hand
)
