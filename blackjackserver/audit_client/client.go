package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"blackjack/messages"
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

func sendData(conn *net.Conn, msg string) (string, error) {
	err := wsutil.WriteClientText(*conn, []byte(msg))
	if err != nil {
		fmt.Printf("Send failed")
		return "", err
	}
	return "OK", err
}

func readData(conn *net.Conn) (string, error) {
	msg_bytes, err := wsutil.ReadServerText(*conn)
	if err != nil {
		fmt.Printf("Receive failed: %s\n", err)
		return "", err
	}
	msg := string(msg_bytes)
	return msg, err
}

func play(io *AuditIO) {
	// Start websocket connection to blackjack server
	conn, _, _, err := ws.DefaultDialer.Dial(
		context.Background(),
		fmt.Sprintf("ws://%s/", "localhost:8080"),
	)
	if err != nil {
		fmt.Printf("Can not connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected")
	fmt.Println("Getting session token from file")
	sessionTokenMsg := io.ReadData()
	sendData(&conn, sessionTokenMsg)

	for {
		fmt.Println("Waiting for ws message")
		msg, err := readData(&conn)
		if err != nil {
			fmt.Printf("read data failed: %s\n", err)
			return
		}
		fmt.Printf("Received: %s\n", msg)
		if msg == messages.PLAYING_THIS_HAND_MSG {
			fmt.Println("Sending in message")
			sendData(&conn, io.ReadData())
		} else if msg == messages.BUST_MSG {
			fmt.Println("BUST")
			continue
		}
		isDealerHand, player, err := messages.DecodePlayerHandMessage(msg)
		// TODO: send when it's this current player's turn
		if err == nil && !isDealerHand {
			if player.Hands[0].IsBust() {
				continue
			}
			fmt.Println("Sending message")
			sendData(&conn, io.ReadData())
		}
	}
}

func main() {
	filename := "../audit/95fa3923-00f4-4d06-90cd-6db76f649109.log"
	io := MakeAuditIO(filename)
	// go play(io)
	play(io)
}
