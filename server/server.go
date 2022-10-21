package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	settings "blackjack/configs"
)

func ReadData(conn net.Conn) string {
	msg, err := wsutil.ReadClientText(conn)
	if err != nil {
		fmt.Println("Read failed, ", err)
	}

	msg_str := string(msg)
	// fmt.Println("Received: ", msg_str)
	return msg_str
}

func SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerText(conn, []byte(msg))
	if err != nil {
		fmt.Println("Send failed, ", err)
	}
}

type Room struct {
	connections   []net.Conn
	currentPlayer int
	Log           *logrus.Logger
	Id            string
}

func MakeRoom() *Room {
	room := Room{}
	room.Id = uuid.New().String()
	room.Log = MakeLog(room.Id)

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

func (room *Room) GetCurrPlayerConn() *net.Conn {
	return &room.connections[room.currentPlayer]
}

func (room *Room) ChangePlayer() {
	room.currentPlayer += 1
	room.currentPlayer %= len(room.connections) // TODO if player disconnects
}

func (room *Room) SendAll(msg string) {
	for _, conn := range room.connections {
		// fmt.Println("sending all ", msg)
		SendData(conn, msg)
	}
}

type Server struct {
	room               *Room
	newConnectionMutex sync.Mutex
	newConnectionCond  sync.Cond
	connQueue          []net.Conn
}

func (server *Server) registerPlayer(conn *net.Conn) {
	server.newConnectionMutex.Lock()
	defer server.newConnectionMutex.Unlock()
	fmt.Println("Registering player")

	server.connQueue = append(server.connQueue, *conn)
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
	room.connections = server.connQueue[:settings.RoomSize]
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
				// handle error
				fmt.Println("Upgrade error, ", err)
			}

			fmt.Println("New player connected")
			go server.registerPlayer(&conn)
		}))
}
