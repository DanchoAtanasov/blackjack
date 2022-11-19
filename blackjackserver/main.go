package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	settings "blackjack/config"
	"blackjack/messages"
	"blackjack/models"
	"blackjack/server"
)

const DIVIDER string = "---------------------------------"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")

type senderFunc func(net.Conn, string, server.Room)
type readerFunc func(net.Conn, models.Hand) string

func sendPlayer(conn net.Conn, message string, room server.Room) {
	server.SendData(conn, message)
}

func sendDealer(conn net.Conn, message string, room server.Room) {
	room.SendAll(message)
}

func readPlayerAction(conn net.Conn, hand models.Hand) string {
	var input string
	retries := 5
	for {
		input = server.ReadData(conn)
		if input == messages.HIT_MSG || input == messages.STAND_MSG {
			break
		}
		if input == "EOF" {
			return "Out"
		}
		retries -= 1
		if retries == 0 {
			return "Out"
		}
		fmt.Printf("Wrong input %s, Try again\n", input)
	}
	return input
}

func readDealerAction(conn net.Conn, hand models.Hand) string {
	if hand.Sum > 17 {
		return messages.STAND_MSG
	}
	// Dealer must hit a soft 17
	if hand.Sum == 17 && hand.NumAces <= 0 {
		return messages.STAND_MSG
	}
	return messages.HIT_MSG
}

func saveResultToFile(players []models.Player, id string) {
	var outputString string
	for i := range players {
		outputString += fmt.Sprintf("%s: %d\n", players[i].Name, players[i].BuyIn)
	}
	os.Mkdir("results", 0755)
	err := os.WriteFile(fmt.Sprintf("./results/%s.log", id), []byte(outputString), 0666)
	if err != nil {
		fmt.Println(err)
		// TODO catch error
	}
	fmt.Println("Results saved to file: ", id)
}

func playTurn(
	player *models.Player,
	deck *models.Deck,
	conn net.Conn,
	readAction readerFunc,
	sendAction senderFunc,
	log *logrus.Logger,
	room *server.Room,
	disconnectedPlayers *[]string,
) {
	for {
		log.Printf("%s's hand is %v", player.Name, player.Hand.Cards)
		log.Printf("Current count: %d", player.Hand.Sum)

		// Send current hand
		if player.Name == "Dealer" {
			sendAction(conn, messages.DEALER_HAND_MSG(*player), *room)
		} else {
			room.SendAll(messages.PLAYER_HAND_MSG(*player))
		}

		if player.Hand.IsBlackjack {
			room.Log.Info("Blackjack!")
			sendAction(conn, messages.BLACKJACK_MSG, *room)
			break
		}

		if player.Hand.IsBust() {
			log.Info("Over 21, bust")
			sendAction(conn, messages.BUST_MSG, *room)
			break
		}

		// Read action
		input := readAction(conn, player.Hand)
		if input == messages.STAND_MSG {
			break
		} else if input == "Out" {
			fmt.Println("Removing disconnected player")
			room.RemoveDisconnectedPlayer()
			*disconnectedPlayers = append(*disconnectedPlayers, player.Name)
			break
		}

		player.Hand.AddCard(deck.DealCard())
	}
}

func clearHands(players []models.Player) {
	for i := range players {
		models.ClearHand(&players[i].Hand)
	}
}

func removeDisconnectedPlayers(players []models.Player, disconnectedPlayersThisTurn []string) []models.Player {
	var activePlayers []models.Player
	fmt.Printf("disconnectedPlayersThisTurn: %v\n", disconnectedPlayersThisTurn)
	for _, player := range players {
		isActive := true
		for _, disconnectedPlayerName := range disconnectedPlayersThisTurn {
			if player.Name == disconnectedPlayerName {
				isActive = false
			}
		}
		if isActive {
			activePlayers = append(activePlayers, player)
		}
	}
	return activePlayers
}

func play(deck *models.Deck, playersPtr *[]models.Player, room *server.Room) {
	dealer := &models.Player{Name: "Dealer"}

	players := *playersPtr
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}
	dealer.Hand.AddCard(deck.DealCard())
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}

	room.Log.Printf("Dealer's hand: %v", dealer.Hand.Cards)
	room.SendAll(messages.LIST_PLAYERS_MSG(players))
	room.SendAll(messages.DEALER_HAND_MSG(*dealer))

	var disconnectedPlayersThisTurn []string
	currConn := *room.GetCurrPlayerConn()
	// Players' turn
	for i := range players {
		currPlayer := &players[i]
		room.Log.Printf("%s's turn, buy in: %d", currPlayer.Name, currPlayer.BuyIn)
		currConn = *room.GetCurrPlayerConn()

		playTurn(currPlayer, deck, currConn, readPlayerAction, sendPlayer, room.Log, room, &disconnectedPlayersThisTurn)

		room.ChangePlayer()
		room.Log.Info(DIVIDER)
	}

	// Dealer's turn
	playTurn(dealer, deck, currConn, readDealerAction, sendDealer, room.Log, room, &disconnectedPlayersThisTurn)
	room.Log.Info(DIVIDER)

	for i := range players {
		currPlayer := &players[i]
		switch models.GetWinner(currPlayer.Hand, dealer.Hand) {
		case 2:
			room.Log.Printf("%s had Blackjack, gets 3x bet", currPlayer.Name)
			currPlayer.Win()
			currPlayer.Win()
		case 1:
			room.Log.Printf("%s wins!", currPlayer.Name)
			currPlayer.Win()
		case -1:
			room.Log.Info("Dealer wins. :(")
			currPlayer.Lose()
		case 0:
			room.Log.Info("Draw")
		}
	}

	clearHands(players)
	*playersPtr = removeDisconnectedPlayers(players, disconnectedPlayersThisTurn)
	time.Sleep(5 * time.Second)
}

type PlayerDetails struct {
	Name    string
	BuyIn   int
	CurrBet int
}

func fetchPlayerDetails(token string) PlayerDetails {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", REDIS_HOST),
		Password: "",
		DB:       0,
	})

	val, err := client.Get(token).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("got from redis %s\n", val)

	var pd PlayerDetails
	err = json.Unmarshal([]byte(val), &pd)
	if err != nil {
		fmt.Println("failed to unmarshal")
	}

	return pd
}

func getPlayerDetails(conn net.Conn) PlayerDetails {
	type Token struct {
		Token string
	}
	var token Token
	msg := server.ReadData(conn)
	fmt.Printf("Received %s", msg)
	err := json.Unmarshal([]byte(msg), &token)
	if err != nil {
		fmt.Println("Failed to unmarshal")
	}
	fmt.Println(token)
	return fetchPlayerDetails(token.Token)
}

func playRoom(room *server.Room, server2 *server.Server) {
	room.Log.Info("Getting a new shuffled deck of cards")
	deck := models.GetNewShuffledDeck(settings.NumDecksInShoe)

	var players []models.Player
	for i := 0; i < settings.RoomSize; i++ {
		room.Log.Info("Getting player details")
		currConn := room.GetCurrPlayerConn()
		room.ChangePlayer()
		pd := getPlayerDetails(*currConn)

		players = append(players, models.Player{
			Name:       pd.Name,
			BuyIn:      pd.BuyIn,
			CurrentBet: pd.CurrBet,
			// NOTE: since Hand is empty when JSON serialised it will be sent a null
			// so it's handled by the frontend. Maybe change the serialization by
			// making a custom serializer or instantiating the hand beforehand, pun intended
		})
	}

	fmt.Println(players)
	room.Log.Info("Lets play!")
	room.SendAll(messages.START_MSG)

	for round := 0; round < settings.NumRoundsPerGame; round++ {
		room.Log.Printf("----------Round %d----------", round+1)
		play(&deck, &players, room)
		deck = *models.ShuffleDeckIfLow(&deck, 150)
	}

	room.Log.Info(DIVIDER)
	room.Log.Info("Final buy ins: ")
	// TODO: Disconnected players winnings are not recorded
	for i := range players {
		room.Log.Printf("%s: %d", players[i].Name, players[i].BuyIn)
	}
	room.SendAll(messages.LIST_PLAYERS_MSG(players))

	go saveResultToFile(players, room.Id)
	room.SendAll(messages.OVER_MSG)
}

func main() {
	fmt.Println("Welcome to Blackjack")
	fmt.Println("Running server:")
	output := server.MakeServer()
	go output.Serve()
	for {
		currRoom := output.WaitForPlayers()
		go playRoom(currRoom, output)
	}
}
