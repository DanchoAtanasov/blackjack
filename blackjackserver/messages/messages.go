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

func leave() string {
	return makeActionMessage("Leave")
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

func playingThisHand() string {
	return makeGameMessage("IN")
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

// TODO improve names
type PlayerInResponse struct {
	Playing    bool
	CurrentBet int
}

func playingThisHandResponse() string {
	return makeActionMessage("In")
}

func DecodePlayerInMessage(msg string) (PlayerInResponse, error) {
	var message Message
	err := json.Unmarshal([]byte(msg), &message)
	var response PlayerInResponse
	err = json.Unmarshal([]byte(message.Message), &response)
	return response, err
}

func dealerHandMessage(player models.Player) string {
	return makeJsonMessage("DealerHand", player.Hand.ToJson())
}

func handMessage(player models.Player) string {
	if player.IsDealer {
		return dealerHandMessage(player)
	}
	return playerHandMessage(player)
}

type listPlayersFunc func([]*models.Player) string

func listPlayers(players []*models.Player) string {
	playerMessageBytes, _ := json.Marshal(players)
	return makeJsonMessage("ListPlayers", string(playerMessageBytes))
}

// TODO: Add decode for dealer hand message

var (
	START_MSG                  string          = gameStart()       // {"type":"Game","message":"Start"}
	OVER_MSG                   string          = gameOver()        // {"type":"Game","message":"Over"}
	HIT_MSG                    string          = hit()             // {"type":"PlayerAction","message":"Hit"}
	STAND_MSG                  string          = stand()           // {"type":"PlayerAction","message":"Stand"}
	LEAVE_MSG                  string          = leave()           // {"type":"PlayerAction","message":"Leave"}
	BUST_MSG                   string          = bust()            // {"type":"HandState","message":"Bust"}
	BLACKJACK_MSG              string          = blackjack()       // {"type":"HandState","message":"Blackjack"}
	LIST_PLAYERS_MSG           listPlayersFunc = listPlayers       // {"type":"ListPlayers","message":"[]"}
	DEALER_HAND_MSG            handMessageFunc = dealerHandMessage // {"type":"DealerHand","message":""}
	HAND_MSG                   handMessageFunc = handMessage       // {"type":"DealerHand" if dealer else "PlayerHand" ,"message":""}
	PLAYING_THIS_HAND_MSG      string          = playingThisHand()
	PLAYING_THIS_HAND_RESP_MSG string          = playingThisHandResponse()
	// PLAYER_HAND_MSG  handMessageFunc = playerHandMessage // {"type":"PlayerHand","message":""}
)
