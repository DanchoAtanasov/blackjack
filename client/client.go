package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func sendData(conn net.Conn, msg string) string {
	fmt.Println("Sending message...")
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
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://127.0.0.1:8080/")
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
		fmt.Println("Sending S")
		res := sendData(conn, "S")
		if res == "Failed" {
			break
		}

		fmt.Println("Reading")
		res = readData(conn)
		if res == "Failed" {
			break
		}
		fmt.Println(res)
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
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go play(i, &wg)
	}
	wg.Wait()
}
