package main

import (
	"fmt"
	"strconv"

	settings "blackjack/configs"
	"blackjack/models"
	"blackjack/server"
)

func sendPlayerCount(count int, room *server.Room) {
	server.SendData(*room.GetCurrPlayerConn(), strconv.Itoa(count))
}

func readPlayerAction(room *server.Room) string {
	fmt.Println("Hit(H) or Stand(S)")
	var input string
	for {
		input = server.ReadData(*room.GetCurrPlayerConn())
		if input == "H" || input == "S" {
			break
		}
		fmt.Println("Try again")
	}
	return input
}

func takeAction(playerName string, hand *models.Hand, deck *models.Deck, room *server.Room) {
	for {
		fmt.Printf("%s's hand is %v\n", playerName, hand.Cards)
		fmt.Printf("Current count: %d\n", hand.Sum)

		if hand.IsBust() {
			fmt.Println("Over 21, bust")
			break
		}

		if playerName == "Dealer" {
			if hand.Sum > 17 {
				break
			} else if hand.Sum == 17 && hand.NumAces <= 0 {
				break
			}
		} else {
			sendPlayerCount(hand.Sum, room)
			input := readPlayerAction(room)
			if input == "S" {
				break
			}
		}

		hand.AddCard(deck.DealCard())
	}
	if playerName != "Dealer" {
		room.ChangePlayer()
	}
}

func clearHands(players []models.Player) {
	for i := range players {
		models.ClearHand(&players[i].Hand)
	}
}

func play(deck *models.Deck, players []models.Player, room *server.Room) {
	dealerHand := models.Hand{}

	fmt.Println("Dealing")
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}
	dealerHand.AddCard(deck.DealCard())
	for i := range players {
		players[i].Hand.AddCard(deck.DealCard())
	}

	fmt.Printf("Dealer's hand: %v\n", dealerHand.Cards)
	room.SendAll(dealerHand.ToJson())

	// Players' turn
	for i := range players {
		fmt.Printf("%s's turn, buy in: %d\n", players[i].Name, players[i].BuyIn)
		// Check for Blackjack
		if players[i].Hand.Sum == 21 {
			fmt.Printf("Hand is %v\n", players[i].Hand.Cards)
			fmt.Println("Blackjack!")
			players[i].Hand.IsBlackjack = true
			fmt.Println("---------------------------------")
			continue
		}

		takeAction(players[i].Name, &players[i].Hand, deck, room)
		fmt.Println("---------------------------------")
	}

	// Dealer's turn
	takeAction("Dealer", &dealerHand, deck, room)
	fmt.Println("---------------------------------")

	for i := range players {
		switch models.GetWinner(players[i].Hand, dealerHand) {
		case 2:
			fmt.Printf("%s had Blackjack, gets 3x bet\n", players[i].Name)
			players[i].Win()
			players[i].Win()
		case 1:
			fmt.Printf("%s wins!\n", players[i].Name)
			players[i].Win()
		case -1:
			fmt.Println("Dealer wins. :(")
			players[i].Lose()
		case 0:
			fmt.Println("Draw")
		}
	}

	clearHands(players)
}

func playRoom(room *server.Room) {
	fmt.Println("Getting a fresh deck of cards and shuffling")
	deck := models.GetNewDeck(6)

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
