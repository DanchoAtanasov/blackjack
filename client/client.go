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

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Player struct {
	Name  string
	BuyIn int
}
type PlayerRequestResponse struct {
	GameServer string
	Token      string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var hostName string = getEnv("BJ_HOST", "localhost")

func sendData(conn net.Conn, msg string) (string, error) {
	err := wsutil.WriteClientText(conn, []byte(msg))
	if err != nil {
		fmt.Printf("Send failed")
		return "", err
	}
	// fmt.Println("Sent ", msg)
	return "OK", err
}

func readData(conn net.Conn, i int) (string, error) {
	msg_bytes, err := wsutil.ReadServerText(conn)
	if err != nil {
		fmt.Printf("[%d] Receive failed, err %s\n", i, err)
		return "", err
	}
	msg := string(msg_bytes)
	// fmt.Println("Received ", msg)
	return msg, err
}

func play(i int, wg *sync.WaitGroup, player Player, pl PlayerRequestResponse) {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), fmt.Sprintf("ws://%s/", pl.GameServer))
	defer wg.Done()
	if err != nil {
		fmt.Printf("[%d] can not connect: %v\n", i, err)
		return
	}
	// defer conn.Close()

	fmt.Printf("[%d] connected\n", i)
	type DetailRequest struct {
		Token string
	}
	detailRequest := DetailRequest{
		Token: pl.Token,
	}
	fmt.Print("sending details\n", i)
	detailRequestString, err := json.Marshal(detailRequest)
	if err != nil {
		fmt.Printf("[%d] can not marshal details: %v\n", i, err)
	}
	_, err = sendData(conn, string(detailRequestString))
	if err != nil {
		fmt.Printf("[%d] can not send details: %v\n", i, err)
	}

	fmt.Printf("[%d] Waiting for game to begin\n", i)

	startMsg, err := readData(conn, i)
	if startMsg != "Start" {
		fmt.Printf("[%d] Wrong start msg received: %s\n", i, startMsg)
	}
	fmt.Printf("[%d] Game started\n", i)

	// Round loop
	for {
		dealerHand, err := readData(conn, i)
		if err != nil {
			fmt.Println("Reading dealer hand failed, exiting")
			break
		}

		fmt.Printf("[%d] Dealer's hand: %s\n", i, dealerHand)
		if dealerHand == "Over" {
			fmt.Println("Game is over, ending")
			break
		}

		for {
			currentCountString, err := readData(conn, i)
			if err != nil {
				break
			}

			if currentCountString == "Blackjack" {
				fmt.Printf("[%d] got Blackjack!\n", i)
				break
			}

			if currentCountString == "Bust" {
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
				action = "H"
			} else {
				action = "S"
			}

			_, err = sendData(conn, action)
			if err != nil {
				break
			}

			if action == "S" {
				break
			}
		}
	}

	err = conn.Close()
	if err != nil {
		fmt.Printf("[%d] can not close: %v\n", i, err)
	} else {
		fmt.Printf("[%d] closed\n", i)
	}
}

func findServer(player Player) PlayerRequestResponse {
	// TODO: Rename function and clean up
	postBody, _ := json.Marshal(player)
	responseBody := bytes.NewBuffer(postBody)
	// TODO: change hostname to env variable
	resp, err := http.Post("http://localhost:3333/play", "application/json", responseBody)
	if err != nil {
		fmt.Printf("An Error Occured %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var playerRequestResponse PlayerRequestResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(body), &playerRequestResponse)
	if err != nil {
		fmt.Println("could not unmarshal body")
	}

	fmt.Printf("got %s, %s\n", playerRequestResponse.GameServer, playerRequestResponse.Token)
	return playerRequestResponse
}

func main() {
	var wg sync.WaitGroup
	numPlayers := 6
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		player := Player{
			Name:  fmt.Sprintf("Player%d", i),
			BuyIn: 100 * (i + 1),
		}
		pl := findServer(player)
		go play(i, &wg, player, pl)
	}
	wg.Wait()
}
