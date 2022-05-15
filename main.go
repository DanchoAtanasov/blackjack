package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"

	settings "blackjack/configs"
	"blackjack/models"
	"blackjack/server"
)

type senderFunc func(net.Conn, string)
type readerFunc func(net.Conn, models.Hand) string

func sendPlayer(conn net.Conn, message string) {
	server.SendData(conn, message)
}

func sendDealer(conn net.Conn, message string) {
	// Do nothing, the dealer is in this server
}

func readPlayerAction(conn net.Conn, hand models.Hand) string {
	fmt.Println("Hit(H) or Stand(S)")
	var input string
	for {
		input = server.ReadData(conn)
		if input == "H" || input == "S" {
			break
		}
		fmt.Println("Try again") // TODO fix this loop
	}
	return input
}

func readDealerAction(conn net.Conn, hand models.Hand) string {
	if hand.Sum > 17 {
		return "S"
	}
	if hand.Sum == 17 && hand.NumAces <= 0 {
		return "S"
	}
	return "H"
}

func saveResultToFile(players []models.Player) {
	var outputString string
	for i := range players {
		outputString += fmt.Sprintf("%s: %d\n", players[i].Name, players[i].BuyIn)
	}
	filename := uuid.New()
	// TODO fix path for docker
	err := os.WriteFile(fmt.Sprintf("./%s.log", filename.String()), []byte(outputString), 0666)
	if err != nil {
		fmt.Println(err)
		// TODO catch error
	}
	fmt.Println("Results saved to file: ", filename.String())
}

func playTurn(
	player *models.Player,
	deck *models.Deck,
	conn net.Conn,
	readAction readerFunc,
	sendAction senderFunc,
) {
	for {
		fmt.Printf("%s's hand is %v\n", player.Name, player.Hand.Cards)
		fmt.Printf("Current count: %d\n", player.Hand.Sum)

		if player.Hand.IsBust() {
			fmt.Println("Over 21, bust")
			sendAction(conn, "Bust")
			break
		}

		// Send current count
		sendAction(conn, strconv.Itoa(player.Hand.Sum))

		// Read action
		input := readAction(conn, player.Hand)
		if input == "S" {
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

	fmt.Println("Dealing")
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}
	dealer.Hand.AddCard(deck.DealCard())
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}

	fmt.Printf("Dealer's hand: %v\n", dealer.Hand.Cards)
	room.SendAll(dealer.Hand.ToJson())

	currConn := *room.GetCurrPlayerConn()
	// Players' turn
	for i := range players {
		currPlayer := &players[i]
		fmt.Printf("%s's turn, buy in: %d\n", currPlayer.Name, currPlayer.BuyIn)
		currConn = *room.GetCurrPlayerConn()

		// Check for Blackjack
		if currPlayer.Hand.IsBlackjack {
			fmt.Printf("Hand is %v\n", currPlayer.Hand.Cards)
			fmt.Println("Blackjack!")
			sendPlayer(currConn, "Blackjack")
		} else {
			playTurn(currPlayer, deck, currConn, readPlayerAction, sendPlayer)
		}

		room.ChangePlayer()
		fmt.Println("---------------------------------")
	}

	// Dealer's turn
	playTurn(dealer, deck, currConn, readDealerAction, sendDealer)
	fmt.Println("---------------------------------")

	for i := range players {
		currPlayer := &players[i]
		switch models.GetWinner(currPlayer.Hand, dealer.Hand) {
		case 2:
			fmt.Printf("%s had Blackjack, gets 3x bet\n", currPlayer.Name)
			currPlayer.Win()
			currPlayer.Win()
		case 1:
			fmt.Printf("%s wins!\n", currPlayer.Name)
			currPlayer.Win()
		case -1:
			fmt.Println("Dealer wins. :(")
			currPlayer.Lose()
		case 0:
			fmt.Println("Draw")
		}
	}

	clearHands(players)
}

func playRoom(room *server.Room) {
	fmt.Println("Getting a new shuffled deck of cards")
	deck := models.GetNewShuffledDeck(settings.NumDecksInShoe)

	var players []models.Player
	for i := 0; i < settings.RoomSize; i++ {
		players = append(players, models.Player{
			Name:       "Player " + strconv.Itoa(i+1),
			BuyIn:      settings.InitialBuyIn,
			CurrentBet: settings.CurrBet,
		})
	}

	fmt.Println("Lets play!")
	for round := 0; round < settings.NumRoundsPerGame; round++ {
		fmt.Printf("----------Round %d----------\n", round+1)
		play(&deck, players, room)
		deck = *models.ShuffleDeckIfLow(&deck, 150)
	}

	fmt.Println("---------------------------------")
	fmt.Println("Final buy ins: ")
	for i := range players {
		fmt.Printf("%s: %d\n", players[i].Name, players[i].BuyIn)
	}
	go saveResultToFile(players)
	room.SendAll("Over")
}

func main() {
	fmt.Println("Welcome to Blackjack")
	fmt.Println("Running server:")
	output := server.MakeServer()
	go output.Serve()
	for {
		currRoom := output.GetLastRoom()
		currRoom.WaitForPlayers()
		go playRoom(currRoom)
	}
}
