package main

import (
	"context"
	"fmt"
	"io"
	"math"
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

func play(io *ConnIO, wg *sync.WaitGroup, roundCond *sync.Cond, currRound *int, playerConn *PlayerConn) {
	defer wg.Done()

	roundCond.L.Lock()

	for {
		fmt.Printf("Waiting on cond %s\n", playerConn.playerDetails.Name)
		roundCond.Wait()
		// fmt.Printf("%s unlocked\n", playerConn.playerDetails.Name)
		if *currRound >= playerConn.startRound {
			fmt.Printf("%s will start playing\n", playerConn.playerDetails.Name)
			break
		}
	}
	roundCond.L.Unlock()

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

	fmt.Printf("%s Connected\n", playerConn.playerDetails.Name)
	fmt.Println("Getting session token from file")
	sessionTokenMsg := io.ReadData()
	sendData(&conn, sessionTokenMsg)

	oldCurrRound := *currRound

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
			roundCond.L.Lock()
			fmt.Printf("Comparing %d and %d\n", oldCurrRound+1, *currRound)
			*currRound = int(math.Max(float64(oldCurrRound+1), float64(*currRound)))
			roundCond.Broadcast()
			roundCond.L.Unlock()
			sendData(&conn, io.ReadData())
			continue
		} else if msg == messages.BUST_MSG {
			fmt.Println("BUST")
			continue
		}
		_, player, err := messages.DecodePlayerHandMessage(msg)
		if err == nil {
			if player.Name != playerConn.playerDetails.Name {
				continue
			}
			if player.Hands[0].IsBust() {
				continue
			}
			fmt.Println("Sending message")
			sendMsg := io.ReadData()
			sendData(&conn, sendMsg)
			if sendMsg == messages.LEAVE_MSG {
				break
			}
		}
	}
}

type PlayerConn struct {
	sessionId     string
	playerDetails *server.PlayerDetails
	startRound    int
}

func main() {
	var wg sync.WaitGroup
	var currRound int = 0
	roundMutex := &sync.Mutex{}
	roundCond := sync.NewCond(roundMutex)

	sessionID := "8a905114-45f2-4d41-8d82-f084e58058cb"
	sessionLogFileName := fmt.Sprintf("../audit/%s.log", sessionID)
	sessionIO := MakeSessIO(sessionLogFileName)

	playersInSession := []PlayerConn{}

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
		playersInSession = append(playersInSession, PlayerConn{
			sessionId:     input.SessionId,
			playerDetails: &input.PlayerDetails,
			startRound:    input.Round,
		})

	}

	time.Sleep(2 * time.Second)

	for i := range playersInSession {
		playerSession := playersInSession[i]
		wg.Add(1)
		connectionLogFileName := fmt.Sprintf("../audit/%s.log", playerSession.sessionId)
		connIO := MakeConnIO(connectionLogFileName)

		go play(connIO, &wg, roundCond, &currRound, &playerSession)

	}
	time.Sleep(1 * time.Second)
	roundCond.Broadcast()
	wg.Wait()
	fmt.Println("All players finished")
}
