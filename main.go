package main

import (
	"fmt"
	"strconv"

	// "blackjack/ui"
	"blackjack/models"
	"blackjack/server"
)

func sendPlayerCount(count int, output *server.Server) {
	server.SendData(*output.GetCurrPlayerConn(), strconv.Itoa(count))
}

func readPlayerAction(output *server.Server) string {
	fmt.Println("Hit(H) or Stand(S)")
	var input string
	for {
		// fmt.Scanln(&input)
		input = server.ReadData(*output.GetCurrPlayerConn())
		if input == "H" || input == "S" {
			break
		}
		fmt.Println("Try again")
	}
	return input
}

func takeAction(playerName string, hand *models.Hand, deck *models.Deck, output *server.Server) {
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
			sendPlayerCount(hand.Sum, output)
			input := readPlayerAction(output)
			if input == "S" {
				break
			}
		}

		hand.AddCard(deck.DealCard())
	}
	if playerName != "Dealer" {
		output.ChangePlayer()
	}
}

func clearHands(players []models.Player) {
	for i := range players {
		models.ClearHand(&players[i].Hand)
	}
}

func play(deck *models.Deck, numPlayers int, players []models.Player, output *server.Server) {
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
	output.SendAll(dealerHand.ToJson())

	// Players' turn
	for i := 0; i < numPlayers; i++ {
		fmt.Printf("%s's turn, buy in: %d\n", players[i].Name, players[i].BuyIn)
		// Check for Blackjack
		if players[i].Hand.Sum == 21 {
			fmt.Printf("Hand is %v\n", players[i].Hand.Cards)
			fmt.Println("Blackjack!")
			players[i].Hand.IsBlackjack = true
			fmt.Println("---------------------------------")
			continue
		}

		takeAction(players[i].Name, &players[i].Hand, deck, output)
		fmt.Println("---------------------------------")
	}

	// Dealer's turn
	takeAction("Dealer", &dealerHand, deck, output)
	fmt.Println("---------------------------------")

	for i := 0; i < numPlayers; i++ {
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
			// No change in money
		}
	}

	clearHands(players)
}

func main() {
	fmt.Println("Welcome to Blackjack")
	fmt.Println("Running server:")
	output := server.NewServer()
	go output.Serve()
	output.WaitForPlayers()

	fmt.Println("Getting a fresh deck of cards")
	deck := models.Deck{}
	deck.BuildDeck(6)

	fmt.Println("Shuffling cards...")
	models.ShuffleDeck(deck)
	// models.ShanoShuffleDeck(&deck) // TODO remove

	numRounds := 100
	numPlayers := 6

	var players []models.Player
	for i := 0; i < numPlayers; i++ {
		players = append(players, models.Player{
			Name:       "Player " + strconv.Itoa(i+1),
			BuyIn:      100,
			CurrentBet: 1,
		})
	}

	fmt.Println("Lets play!")
	for round := 0; round < numRounds; round++ {
		fmt.Printf("----------Round %d----------\n", round+1)
		play(&deck, numPlayers, players, &output)
		// TODO reshuffle when deck is low
		deck = *models.ShuffleDeckIfLow(&deck, 150)
	}

	fmt.Println("---------------------------------")
	fmt.Println("Final buy ins: ")
	for i := range players {
		fmt.Printf("%s: %d\n", players[i].Name, players[i].BuyIn)
	}
}
