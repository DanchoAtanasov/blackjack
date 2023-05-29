package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"

	settings "blackjack/config"
	"blackjack/models"
)

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
	msg := server.room.IO.ReadData(*conn)
	if msg == "EOF" {
		return
	}

	fmt.Printf("Received %s", msg)
	err := json.Unmarshal([]byte(msg), &sessionJwt)
	if err != nil {
		fmt.Println("Failed to unmarshal")
		// TODO: handle failure
	}

	fmt.Println(sessionJwt)
	sessionId := parseJwt(sessionJwt.Token)

	fmt.Printf("Making connection log: %s", sessionId)
	connectionLog := MakeAuditLog(sessionId)
	connectionLog.Info(msg)

	fmt.Println("Getting player details")
	pd := getPlayerDetails(sessionId)

	// Log to room audit log
	server.room.Audit.WithFields(
		logrus.Fields{
			"type":          "system",
			"action":        "newPlayer",
			"sessionId":     sessionId,
			"playerDetails": pd,
			"round":         server.room.CurrRound,
		},
	).Info()

	newPlayer := models.Player{
		Name:    pd.Name,
		BuyIn:   pd.BuyIn,
		CurrBet: pd.CurrBet,
		Hands:   []*models.Hand{{}},
		// NOTE: since Hand is empty when JSON serialised it will be sent a null
		// so it's handled by the frontend. Maybe change the serialization by
		// making a custom serializer or instantiating the hand beforehand, pun intended
		Active: true,
	}

	playerConn := PlayerConn{
		sessionId:     sessionId,
		player:        &newPlayer,
		Conn:          *conn,
		ConnectionLog: connectionLog,
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
