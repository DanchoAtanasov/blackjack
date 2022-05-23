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

func main() {
	var wg sync.WaitGroup
	numPlayers := 12
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go play(i, &wg)
	}
	wg.Wait()
}
