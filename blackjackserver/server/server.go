package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	settings "blackjack/config"
	"blackjack/models"
)

func ReadData(conn net.Conn) string {
	// TODO: Improve connection closed vs read timed out error handling
	// Returns empty string if read failed, EOF if connection was closed
	conn.SetReadDeadline(time.Now().Add(settings.ReadTimeout))

	msg, err := wsutil.ReadClientText(conn)
	if err != nil {
		fmt.Println("Read failed")
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed by client")
		} else if errors.Is(err, wsutil.ClosedError{Code: 1001}) {
			fmt.Println("Connection closed by client, ws closed")
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			fmt.Println("Read timed out")
		} else {
			fmt.Printf("Some other error: %e\n", err)
		}
		return "EOF"
	}

	msg_str := string(msg)
	return msg_str
}

func SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerText(conn, []byte(msg))
	if err != nil {
		fmt.Println("Send failed, ", err)
	}
}

type Server struct {
	room               *Room
	newConnectionMutex sync.Mutex
	newConnectionCond  sync.Cond
	connQueue          []PlayerConn
}

func (server *Server) registerPlayer(conn *net.Conn) {
	server.newConnectionMutex.Lock()
	defer server.newConnectionMutex.Unlock()
	fmt.Println("Registering player")

	fmt.Println("Asking client for a session token")
	var sessionJwt Token
	msg := ReadData(*conn)
	fmt.Printf("Received %s", msg)
	err := json.Unmarshal([]byte(msg), &sessionJwt)
	if err != nil {
		fmt.Println("Failed to unmarshal")
		// TODO: handle failure
	}

	fmt.Println(sessionJwt)
	sessionId := parseJwt(sessionJwt.Token)

	fmt.Println("Getting player details")
	pd := getPlayerDetails(sessionId)

	newPlayer := models.Player{
		Name:    pd.Name,
		BuyIn:   pd.BuyIn,
		CurrBet: pd.CurrBet,
		// NOTE: since Hand is empty when JSON serialised it will be sent a null
		// so it's handled by the frontend. Maybe change the serialization by
		// making a custom serializer or instantiating the hand beforehand, pun intended
		Active: true,
	}

	playerConn := PlayerConn{
		sessionId: sessionId,
		player:    &newPlayer,
		Conn:      *conn,
	}

	server.connQueue = append(server.connQueue, playerConn)
	fmt.Printf("Player %s registered.\n", playerConn.player.Name)
	server.newConnectionCond.Signal()
	// server.newConnectionCond.Broadcast() TODO maybe use broadcast?
}

func MakeServer() *Server {
	server := Server{}
	server.room = MakeRoom()
	server.newConnectionCond = *sync.NewCond(&server.newConnectionMutex)

	return &server
}

func (server *Server) GetRoom() *Room {
	return server.room
}

// Process new players and return new game if necessary
func (server *Server) WaitForPlayers() *Room {
	server.newConnectionCond.L.Lock()
	defer server.newConnectionCond.L.Unlock()

	for {
		// Wait for a new connection
		server.newConnectionCond.Wait()

		fmt.Println("Player found")
		server.room.AddNewPlayerConn(server.connQueue[0])
		server.connQueue = server.connQueue[1:]

		// If it's the first player return room to start the game
		if len(server.room.playerConns) == 1 {
			return server.room
		}

		// Room is full, make a new one
		if len(server.room.playerConns) >= settings.RoomSize {
			server.room = MakeRoom()
		}
	}
}

func (server *Server) Serve() {
	port := 8080
	fmt.Println("Serving on port", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				fmt.Println("Upgrade error, ", err)
				return
			}

			fmt.Println("New player connected")
			go server.registerPlayer(&conn)
		}))
}
