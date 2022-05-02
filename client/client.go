package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func sendData(conn net.Conn, msg string) string {
	err := wsutil.WriteClientMessage(conn, ws.OpText, []byte(msg))
	if err != nil {
		fmt.Printf("Send failed")
		return "Failed"
	}
	fmt.Println("Sent ", msg)
	return "OK"
}

func readData(conn net.Conn) string {
	msg_bytes, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		fmt.Println("Receive failed")
		return "Failed"
	}
	msg := string(msg_bytes)
	fmt.Println("Received ", msg)
	return msg
}

func play(i int, wg *sync.WaitGroup) {
	// TODO add env variable for host
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://app:8080/")
	defer wg.Done()
	if err != nil {
		fmt.Printf("%d can not connect: %v\n", i, err)
		return
	}

	fmt.Printf("%d connected\n", i)
	fmt.Println("Waiting for game to begin")

	dealerHand := readData(conn)
	fmt.Println(dealerHand)

	for {
		currentCountString := readData(conn)
		if currentCountString == "Failed" {
			break
		}
		if currentCountString == "Over" {
			fmt.Println("Game is over, ending")
			break
		}

		currentCount, err := strconv.Atoi(currentCountString)
		if err != nil {
			// TODO fix this, dealer hand is coming here, for now read another message
			currentCountString = readData(conn)
			currentCount, _ = strconv.Atoi(currentCountString)
		}

		fmt.Println("Current hand: ", currentCount)
		var action string
		if currentCount < 16 {
			action = "H"
		} else {
			action = "S"
		}

		res := sendData(conn, action)
		if res == "Failed" {
			break
		}
	}

	err = conn.Close()
	if err != nil {
		fmt.Printf("%d can not close: %v\n", i, err)
	} else {
		fmt.Printf("%d closed\n", i)
	}
}

func main() {
	var wg sync.WaitGroup
	numPlayers := 18
	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go play(i, &wg)
	}
	wg.Wait()
}
