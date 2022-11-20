package main

import (
	"fmt"
	"os"
	"time"

	settings "blackjack/config"
	"blackjack/messages"
	"blackjack/models"
	"blackjack/server"
)

const DIVIDER string = "---------------------------------"

func readPlayerAction(room *server.Room) string {
	var input string
	retries := 5
	for {
		input = room.ReadCurrPlayer()
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

func readDealerAction(hand models.Hand) string {
	if hand.Sum > 17 {
		return messages.STAND_MSG
	}
	// Dealer must hit a soft 17
	if hand.Sum == 17 && hand.NumAces <= 0 {
		return messages.STAND_MSG
	}
	return messages.HIT_MSG
}

func saveResultToFile(players []*models.Player, id string) {
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
	room *server.Room,
) {
	for {
		room.Log.Printf("%s's hand is %v", player.Name, player.Hand.Cards)
		room.Log.Printf("Current count: %d", player.Hand.Sum)

		// Send current hand
		room.SendAll(messages.HAND_MSG(*player))

		if player.Hand.IsBlackjack {
			if !player.IsDealer {
				room.Log.Info("Blackjack!")
				room.SendCurrPlayer(messages.BLACKJACK_MSG)
				// Maybe send all that it's a blackjack
				// room.SendAll(messages.BLACKJACK_MSG)
			}
			break
		}

		if player.Hand.IsBust() {
			room.Log.Info("Over 21, bust")
			room.SendAll(messages.BUST_MSG)
			break
		}

		// Read action
		var input string
		if player.IsDealer {
			input = readDealerAction(player.Hand)
		} else {
			input = readPlayerAction(room)
		}

		if input == messages.STAND_MSG {
			break
		} else if input == "Out" {
			fmt.Println("Removing disconnected player")
			room.RemoveDisconnectedPlayer()
			player.Active = false
			break
		}

		player.Hand.AddCard(deck.DealCard())
	}
}

func clearHands(players []*models.Player) {
	for i := range players {
		models.ClearHand(&players[i].Hand)
	}
}

func calculateWinners(players []*models.Player, dealer models.Player, room server.Room) {
	for i := range players {
		currPlayer := players[i]
		switch models.GetWinner(currPlayer.Hand, dealer.Hand) {
		case 2:
			room.Log.Printf("%s had Blackjack, gets 3x bet", currPlayer.Name)
			currPlayer.Blackjack()
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
}

func playRound(deck *models.Deck, room *server.Room) {
	dealer := &models.Player{Name: "Dealer", IsDealer: true}

	// Note: players is a slice of pointers as they're created in register and this way they can
	// be updated from the game logic code
	players := room.GetPlayers()

	// Ask who's playing
	room.SendAll(messages.PLAYING_THIS_HAND_MSG)
	room.ReadInMessages()

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

	// Players' turn
	for i := range players {
		currPlayer := players[i]
		room.Log.Printf("%s's turn, buy in: %d, bet: %d", currPlayer.Name, currPlayer.BuyIn, currPlayer.CurrentBet)

		if currPlayer.Active {
			playTurn(currPlayer, deck, room)
		}

		room.ChangePlayer()

		room.Log.Info(DIVIDER)
	}

	// Dealer's turn
	playTurn(dealer, deck, room)
	room.Log.Info(DIVIDER)

	calculateWinners(players, *dealer, *room)

	clearHands(players)
	time.Sleep(settings.TimeBetweenRounds)
}

func playRoom(room *server.Room) {
	room.Log.Info("Getting a new shuffled deck of cards")
	deck := models.GetNewShuffledDeck(settings.NumDecksInShoe)

	room.Log.Info("Lets play!")
	room.SendAll(messages.START_MSG)

	for round := 0; round < settings.NumRoundsPerGame; round++ {
		if room.IsEmpty() {
			room.Log.Println("Room is empty, game over.")
			break
		}

		room.Log.Printf("----------Round %d----------", round+1)

		playRound(&deck, room)

		deck = *models.ShuffleDeckIfLow(&deck, 150)
	}

	room.Log.Info(DIVIDER)
	room.Log.Info("Final buy ins: ")
	// TODO: Disconnected players winnings are not recorded
	players := room.GetPlayers()
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
		go playRoom(currRoom)
	}
}
