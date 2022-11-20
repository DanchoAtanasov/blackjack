package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	settings "blackjack/config"
	"blackjack/models"
)

func ReadData(conn net.Conn) string {
	// Returns empty string if read failed, EOF if connection was closed
	msg, err := wsutil.ReadClientText(conn)
	if err != nil {
		fmt.Printf("Read failed, %s, %e\n", err, err)
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed by client")
			return "EOF"
		} else if errors.Is(err, wsutil.ClosedError{Code: 1001}) {
			fmt.Println("Connection closed by client, ws closed")
			return "EOF"
		}
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

	fmt.Println("Getting player details")
	pd := getPlayerDetails(*conn)

	newPlayer := models.Player{
		Name:       pd.Name,
		BuyIn:      pd.BuyIn,
		CurrentBet: pd.CurrBet,
		// NOTE: since Hand is empty when JSON serialised it will be sent a null
		// so it's handled by the frontend. Maybe change the serialization by
		// making a custom serializer or instantiating the hand beforehand, pun intended
		Active: true,
	}

	playerConn := PlayerConn{
		player: &newPlayer,
		Conn:   *conn,
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

func (server *Server) WaitForPlayers() *Room {
	server.newConnectionCond.L.Lock()
	defer server.newConnectionCond.L.Unlock()

	for len(server.connQueue) < settings.RoomSize {
		server.newConnectionCond.Wait()
	}

	fmt.Println("Wait is over!, queue len is ", len(server.connQueue))
	room := MakeRoom()
	room.playerConns = server.connQueue[:settings.RoomSize]
	server.connQueue = server.connQueue[settings.RoomSize:]

	return room
}

func (server *Server) Serve() {
	port := 8080
	fmt.Println("Serving on port", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				fmt.Println("Upgrade error, ", err)
			}

			fmt.Println("New player connected")
			go server.registerPlayer(&conn)
		}))
}
