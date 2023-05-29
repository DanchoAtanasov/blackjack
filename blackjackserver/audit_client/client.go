package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"blackjack/messages"
	"blackjack/server"

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

func play(io *ConnIO, wg *sync.WaitGroup, pd *server.PlayerDetails) {
	defer wg.Done()

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

	fmt.Printf("%s Connected\n", pd.Name)
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
		_, player, err := messages.DecodePlayerHandMessage(msg)
		// TODO: send when it's this current player's turn
		if err == nil {
			if player.Name != pd.Name {
				continue
			}
			if player.Hands[0].IsBust() {
				continue
			}
			fmt.Println("Sending message")
			sendData(&conn, io.ReadData())
		}
	}
}

func main() {
	var wg sync.WaitGroup

	sessionLogFileName := "../audit/9048601e-afd0-47cc-9213-ffcc9b3293a4.log"
	sessionIO := MakeSessIO(sessionLogFileName)

	for {
		input, err := sessionIO.ReadData()
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				break
			} else {
				panic("Session read failed")
			}
		}

		server.SetPlayerDetails(input.SessionId, input.PlayerDetails)
		wg.Add(1)
		connectionLogFileName := fmt.Sprintf("../audit/%s.log", input.SessionId)
		connIO := MakeConnIO(connectionLogFileName)

		time.Sleep(1 * time.Second)

		go play(connIO, &wg, &input.PlayerDetails)

	}
	wg.Wait()
	fmt.Println("All players finished")
}
