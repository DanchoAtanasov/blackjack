package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var hostName string = getEnv("BJ_HOST", "localhost")

func sendData(conn net.Conn, msg string) (string, error) {
	err := wsutil.WriteClientMessage(conn, ws.OpText, []byte(msg))
	if err != nil {
		fmt.Printf("Send failed")
		return "", err
	}
	// fmt.Println("Sent ", msg)
	return "OK", err
}

func readData(conn net.Conn) (string, error) {
	msg_bytes, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		fmt.Println("Receive failed")
		return "", err
	}
	msg := string(msg_bytes)
	// fmt.Println("Received ", msg)
	return msg, err
}

func play(i int, wg *sync.WaitGroup) {
	// TODO add env variable for host
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), fmt.Sprintf("ws://%s:8080/", hostName))
	defer wg.Done()
	if err != nil {
		fmt.Printf("[%d] can not connect: %v\n", i, err)
		return
	}
	// defer conn.Close()

	fmt.Printf("[%d] connected\n", i)
	fmt.Println("Waiting for game to begin")

	startMsg, err := readData(conn)
	if startMsg != "Start" {
		fmt.Printf("[%d] Wrong start msg received: %s\n", i, startMsg)
	}
	fmt.Printf("[%d] Game started\n", i)

	// Round loop
	for {
		dealerHand, err := readData(conn)
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
			currentCountString, err := readData(conn)
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

func main() {
	var wg sync.WaitGroup
	numPlayers := 6
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go play(i, &wg)
	}
	wg.Wait()
}
