package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	settings "blackjack/config"
	"blackjack/models"
)

func ReadData(conn net.Conn) string {
	// Returns empty string if read failed, EOF if connection was closed
	msg, err := wsutil.ReadClientText(conn)
	if err != nil {
		fmt.Printf("Read failed, %s\n", err)
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed by client")
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

// Struct used to couple the player model and the connection to him
type PlayerConn struct {
	player models.Player
	Conn   net.Conn
}

type Room struct {
	Log         *logrus.Logger
	Id          string
	playerConns []PlayerConn
}

func MakeRoom() *Room {
	room := Room{}
	room.Id = uuid.New().String()
	room.Log = MakeLog(room.Id)
	room.playerConns = make([]PlayerConn, 0, settings.RoomSize)

	return &room
}

func MakeLog(id string) *logrus.Logger {
	filename := fmt.Sprintf("./logs/%s.log", id)
	log := logrus.New()
	log.Out = os.Stdout

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	return log
}

func (room *Room) GetCurrPlayerConn() *PlayerConn {
	return &room.playerConns[0]
}

func (room *Room) RemoveDisconnectedPlayer() {
	// TODO improve player diconnecting/rotating logic
	// maybe couple the player and room objects better
	room.playerConns = room.playerConns[1:]
}

func (room *Room) ChangePlayer() {
	// Change players by popping from queue and appending
	var currentConn PlayerConn
	currentConn, room.playerConns = room.playerConns[0], room.playerConns[1:]
	room.playerConns = append(room.playerConns, currentConn)
}

func (room *Room) SendAll(msg string) {
	for _, playerConnection := range room.playerConns {
		SendData(playerConnection.Conn, msg)
	}
}

func (room *Room) GetPlayers() []models.Player {
	var players []models.Player
	for i := 0; i < len(room.playerConns); i++ {
		players = append(players, room.playerConns[i].player)
	}
	return players
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
	}

	playerConn := PlayerConn{
		player: newPlayer,
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
