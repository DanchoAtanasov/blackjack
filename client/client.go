package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"client/messages"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var API_SERVER_HOST string = getEnv("API_SERVER_HOST", "localhost")

type Player struct {
	Name  string
	BuyIn int
}
type GameDetails struct {
	GameServer string
	Token      string
}

func sendData(conn net.Conn, msg string) (string, error) {
	err := wsutil.WriteClientText(conn, []byte(msg))
	if err != nil {
		fmt.Printf("Send failed")
		return "", err
	}
	return "OK", err
}

func readData(conn net.Conn) (string, error) {
	msg_bytes, err := wsutil.ReadServerText(conn)
	if err != nil {
		fmt.Printf("Receive failed: %s\n", err)
		return "", err
	}
	msg := string(msg_bytes)
	return msg, err
}

func sendToken(conn net.Conn, gameDetails GameDetails) {
	type Token struct {
		Token string
	}

	token := Token{
		Token: gameDetails.Token,
	}
	tokenJson, _ := json.Marshal(token)

	_, err := sendData(conn, string(tokenJson))
	if err != nil {
		fmt.Printf("Sending token failed: %v\n", err)
	}
}

func waitForStartMessage(conn net.Conn) {
	startMsg, err := readData(conn)
	if err != nil {
		fmt.Printf("Could not receive start message: %s\n", err)
	}
	if startMsg != messages.START {
		fmt.Printf("Wrong start msg received: %s\n", startMsg)
	}
	// TODO: Add retry
}

// playHand contains logic for a player to play 1 hand.
// Logic is simple and it goes like: if below 16 hit, otherwise stand.
func playHand(conn net.Conn, i int) {
	for {
		currentCountString, err := readData(conn)
		if err != nil {
			break
		}

		if currentCountString == messages.BLACKJACK {
			fmt.Printf("[%d] got Blackjack!\n", i)
			break
		}

		if currentCountString == messages.BUST {
			fmt.Printf("[%d] Bust\n", i)
			break
		}

		currentCount, err := strconv.Atoi(currentCountString)
		if err != nil {
			fmt.Printf("[%d] Error converting count. %v\n", i, err)
			break
		}

		fmt.Printf("[%d]Current hand: %d\n", i, currentCount)
		var action string
		if currentCount < 16 {
			action = messages.HIT
		} else {
			action = messages.STAND
		}

		_, err = sendData(conn, action)
		if err != nil {
			break
		}

		if action == messages.STAND {
			break
		}
	}
}

func playRound(conn net.Conn, i int) {
	for {
		dealerHand, err := readData(conn)
		if err != nil {
			fmt.Println("Reading dealer hand failed, exiting")
			break
		}

		fmt.Printf("[%d] Dealer's hand: %s\n", i, dealerHand)
		if dealerHand == messages.OVER {
			fmt.Println("Game is over, ending")
			break
		}
		playHand(conn, i)
	}
}

func play(i int, wg *sync.WaitGroup, player Player, gameDetails GameDetails) {
	defer wg.Done()

	// Start websocket connection to blackjack server
	conn, _, _, err := ws.DefaultDialer.Dial(
		context.Background(),
		fmt.Sprintf("ws://%s/", gameDetails.GameServer),
	)
	if err != nil {
		fmt.Printf("[%d] can not connect: %v\n", i, err)
		return
	}
	defer conn.Close()

	fmt.Printf("[%s] connected\n", player.Name)

	fmt.Printf("[%d] sending token to blackjack server\n", i)
	sendToken(conn, gameDetails)

	fmt.Printf("[%d] Waiting for game to begin\n", i)
	waitForStartMessage(conn)

	fmt.Printf("[%d] Game started\n", i)

	playRound(conn, i)
}

// registerPlayer send the player details to the API server.
// It returns the game server host name and session token.
func registerPlayer(player Player) GameDetails {
	playerJson, _ := json.Marshal(player)
	requestBody := bytes.NewBuffer(playerJson)
	resp, err := http.Post(
		fmt.Sprintf("http://%s:3333/play", API_SERVER_HOST),
		"application/json",
		requestBody,
	)
	if err != nil {
		fmt.Printf("Could not register player: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var gameDetails GameDetails
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not read response body %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(body), &gameDetails)
	if err != nil {
		fmt.Printf("Response body in wrong format: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Registerred successfully, game details %v\n", gameDetails)
	return gameDetails
}

func main() {
	var wg sync.WaitGroup
	numPlayers := 6
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)

		player := Player{
			Name:  fmt.Sprintf("Player %d", i),
			BuyIn: 100 * (i + 1),
		}
		gameDetails := registerPlayer(player)

		go play(i, &wg, player, gameDetails)
	}
	wg.Wait()
}
