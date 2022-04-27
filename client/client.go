package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func play(i int, wg *sync.WaitGroup) {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://127.0.0.1:8080/")
	defer wg.Done()
	if err != nil {
		fmt.Printf("%d can not connect: %v\n", i, err)
		return
	}

	fmt.Printf("%d connected\n", i)

	msg := []byte("I want to play")
	tries := 0
	for {
		fmt.Println("Sending message...")
		err = wsutil.WriteClientMessage(conn, ws.OpText, msg)
		if err != nil {
			fmt.Printf("%d can not send: %v\n", i, err)
			return
		}
		fmt.Printf("%d send: %s, type: %v\n", i, msg, ws.OpText)

		msg, op, err := wsutil.ReadServerData(conn)
		if err != nil {
			fmt.Printf("%d can not receive: %v\n", i, err)
			return
		}
		fmt.Printf("%d receive: %s, type: %v\n", i, msg, op)
		if bytes.Equal(msg, []byte("OK")) {
			fmt.Println("I'm in ", i)
			break
		}
		if tries >= 3 {
			fmt.Println("Tried, but failed, giving up :(")
			break
		}

		time.Sleep(time.Duration(3) * time.Second)
	}

	msg, op, err := wsutil.ReadServerData(conn)
	if err != nil {
		fmt.Printf("%d can not receive: %v\n", i, err)
		return
	}
	fmt.Printf("%d receive: %s, type: %v\n", i, msg, op)

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
