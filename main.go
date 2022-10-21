package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	settings "blackjack/configs"
	"blackjack/models"
	"blackjack/server"
)

// TODO: convert messages to json
const (
	STAND         string = "S"
	HIT           string = "H"
	BUST_MSG      string = "Bust"
	BLACKJACK_MSG string = "Blackjack"
	OVER_MSG      string = "Over"
	START_MSG     string = "Start"
)
const DIVIDER string = "---------------------------------"

type senderFunc func(net.Conn, string)
type readerFunc func(net.Conn, models.Hand) string

func sendPlayer(conn net.Conn, message string) {
	server.SendData(conn, message)
}

func sendDealer(conn net.Conn, message string) {
	// Do nothing, the dealer is in this server
}

// func readPlayerDetails(conn net.Conn, hand models.Hand) string {

// }

func readPlayerAction(conn net.Conn, hand models.Hand) string {
	var input string
	for {
		input = server.ReadData(conn)
		if input == HIT || input == STAND {
			break
		}
		fmt.Printf("Wrong input %s, Try again", input)
	}
	return input
}

func readDealerAction(conn net.Conn, hand models.Hand) string {
	if hand.Sum > 17 {
		return STAND
	}
	if hand.Sum == 17 && hand.NumAces <= 0 {
		return STAND
	}
	return HIT
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
) {
	for {
		log.Printf("%s's hand is %v", player.Name, player.Hand.Cards)
		log.Printf("Current count: %d", player.Hand.Sum)

		if player.Hand.IsBust() {
			log.Info("Over 21, bust")
			sendAction(conn, BUST_MSG)
			break
		}

		// Send current count
		sendAction(conn, strconv.Itoa(player.Hand.Sum))

		// Read action
		input := readAction(conn, player.Hand)
		if input == STAND {
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

func play(deck *models.Deck, players []models.Player, room *server.Room) {
	dealer := &models.Player{Name: "Dealer"}

	room.Log.Info("Dealing")
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}
	dealer.Hand.AddCard(deck.DealCard())
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}

	room.Log.Printf("Dealer's hand: %v", dealer.Hand.Cards)
	room.SendAll(dealer.Hand.ToJson())

	currConn := *room.GetCurrPlayerConn()
	// Players' turn
	for i := range players {
		currPlayer := &players[i]
		room.Log.Printf("%s's turn, buy in: %d", currPlayer.Name, currPlayer.BuyIn)
		currConn = *room.GetCurrPlayerConn()

		// Check for Blackjack
		if currPlayer.Hand.IsBlackjack {
			room.Log.Printf("Hand is %v", currPlayer.Hand.Cards)
			room.Log.Info("Blackjack!")
			sendPlayer(currConn, BLACKJACK_MSG)
		} else {
			room.Log.Printf("Hit(%s) or Stand(%s)", HIT, STAND)
			playTurn(currPlayer, deck, currConn, readPlayerAction, sendPlayer, room.Log)
		}

		room.ChangePlayer()
		room.Log.Info(DIVIDER)
	}

	// Dealer's turn
	playTurn(dealer, deck, currConn, readDealerAction, sendDealer, room.Log)
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
}

type PlayerDetails struct {
	Name  string
	BuyIn int
}

func fetchPlayerDetails(token string) PlayerDetails {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
			CurrentBet: settings.CurrBet, // TODO: include this in the player details
		})
	}

	room.Log.Info("Lets play!")
	room.SendAll(START_MSG)

	for round := 0; round < settings.NumRoundsPerGame; round++ {
		room.Log.Printf("----------Round %d----------", round+1)
		play(&deck, players, room)
		deck = *models.ShuffleDeckIfLow(&deck, 150)
	}

	room.Log.Info(DIVIDER)
	room.Log.Info("Final buy ins: ")
	for i := range players {
		room.Log.Printf("%s: %d", players[i].Name, players[i].BuyIn)
	}

	go saveResultToFile(players, room.Id)
	room.SendAll(OVER_MSG)
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
